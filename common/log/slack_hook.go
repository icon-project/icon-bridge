package log

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
)

/*
SlackHook defines a mechanism for forwarding logrus entries to a slack channel, whose endpoint is given.
logrus entries are saved to a fixed length buffer, which is read from and forwarded with HTTP Client
For Buffer overflow, last entry is updated with the latest one

Example config:
    "log_forwarder": {
        "vendor":"slack",
        "address":"https://hooks.slack.com/services/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",
        "level":"info"
    }

NOTE: invoking logrus.logging inside some of the functions defined here (like read or write) can cause deadlock
But logging and forwarding errors encountered in these functions is still possible
For a possible way: check forward() where message is concatenated to a string and is then sent directly
instead of relying on logrus.log()
*/

const DefaultHTTPTimeout = 10   // 10 seconds; wait until 10 seconds before HTTP request times out
const STACK_LEN = 1000          // size of buffer that holds logrus entries
const RATE_LIMIT = 1000         // in milliseconds; rate limit of slack is 1 request per second
const MAX_PAYLOAD_SIZE = 100000 //100K bytes; HTTP POST request have >= 2MB, but depends upon server

type SlackHook struct {
	url        string         //webhook url for the slack channel
	levels     []logrus.Level //log levels for which message is forwarded
	httpClient *http.Client   // http client that forwards message
	stack      *slackStack    // thread-safe buffer to hold logs
}

type slackStack struct {
	buf      []*logrus.Entry // STACK_LEN sized array
	mu       sync.Mutex      // mutex to lock all other fields of this structure
	cursor   uint64          // points to the array index where next write should be; previous indexes are filled
	lastRead time.Time       // last time data was read from buffer; used to rate limit reads to RATE_LIMIT milliseconds
}

func NewSlackClient(url string, lvs []logrus.Level) (*SlackHook, error) {
	sh := &SlackHook{
		url:        url,
		levels:     lvs,
		httpClient: &http.Client{Timeout: time.Second * time.Duration(DefaultHTTPTimeout)},
		stack:      &slackStack{buf: make([]*logrus.Entry, STACK_LEN), mu: sync.Mutex{}, cursor: 0, lastRead: time.Now()},
	}
	sh.forward()
	return sh, nil
}

func (sh *SlackHook) Levels() []logrus.Level {
	return sh.levels // sh.levels has been initialized from minimum log_level present in forwarder config
}

// WARNING: Calling logging inside writeToStack will spawn another write function which will wait for the current mutex to release
// while the second function waits for the first one. This causes deadlock.
// So do not call logging inside this function
func (st *slackStack) writeToStack(e *logrus.Entry) {
	st.mu.Lock()
	defer st.mu.Unlock()
	if st.cursor >= STACK_LEN { // If stack full, overwrite the last array element
		st.buf[STACK_LEN-1] = e
	} else {
		st.buf[st.cursor] = e
		st.cursor++
	}
}

// WARNING: Same warning as for writeToStack() because of the same mutex lock being used
func (st *slackStack) readFromStack() []*logrus.Entry {
	st.mu.Lock()
	defer st.mu.Unlock()
	if st.cursor == 0 || time.Now().Sub(st.lastRead).Milliseconds() < RATE_LIMIT { //either nothing to read or too early to read
		return nil
	}
	//copy slice st.buf, do not return it as slice is returned by reference
	// changes outside mutex lock is not thread safe
	newArr := make([]*logrus.Entry, st.cursor)
	for i := 0; i < int(st.cursor); i++ {
		newArr[i] = st.buf[i]
	}
	// everything has been read; reinitialize cursor;
	// this way don't need to delete array entries
	st.cursor = 0
	st.lastRead = time.Now()
	return newArr
}

// Fire is called for every logging invocation
// So, minimize processing overhead inside this function
func (sh *SlackHook) Fire(e *logrus.Entry) (err error) {
	sh.stack.writeToStack(e)
	return
}

func (sh *SlackHook) forward() {
	send := func(reqStr string) (err error) {
		body, err := json.Marshal(map[string]interface{}{"text": string(reqStr)}) //entire message has to be inside "text" key; This key is used by slack API
		if err != nil {
			return errors.Wrap(err, "SlackHook; sendFunc; Json Marshal string; Err:")
		}
		httpReq, err := http.NewRequest("POST", sh.url, bytes.NewReader(body))
		if err != nil {
			return errors.Wrap(err, "SlackHook; sendFunc; http NewRequst; Err:")
		}
		httpReq.Header.Set("Content-Type", "application/json")
		resp, err := sh.httpClient.Do(httpReq)
		if err != nil {
			return errors.Wrap(err, "SlackHook; sendFunc; httpClient Do; Err: ")
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			if t, err := ioutil.ReadAll(resp.Body); err == nil && t != nil {
				err = errors.New("SlackHook; sendFunc; HTTP Request returned Err: " + string(t) + " Status Code: " + strconv.FormatInt(int64(resp.StatusCode), 10))
			} else {
				err = errors.New("SlackHook; sendFunc; HTTP Request returned. " + " Status Code: " + strconv.FormatInt(int64(resp.StatusCode), 10))
			}
		}
		return
	}

	tick := time.NewTicker(time.Millisecond * RATE_LIMIT / 10) // throttle (one-tenth) readFromStack() which mutex locks read/write operation
	go func() {
		reqStr := ""
		bold := "*"
		mono := "```"
		header := ""
		var err error
		defer tick.Stop()
		for {
			select {
			case <-tick.C:
				if res := sh.stack.readFromStack(); res != nil && len(res) > 0 {
					for _, e := range res {
						req := map[string]interface{}{}
						if len(e.Data) > 0 { // If log was invoked WithFields, get the key value pairs
							req = e.Data
						}
						// Add the message passed to logging
						srv := ""
						if vi, ok := req[FieldKeyService]; ok {
							if vs, ok := vi.(string); ok {
								srv = vs
							}
						}
						header = "[" + srv + "]" + "[" + e.Level.String() + "][ICON-BRIDGE][" + e.Time.UTC().Format("2006-01-02T15:04:05.000Z") + "]\n"
						req["Message"] = e.Message
						req["Level"] = e.Level
						req["Time"] = e.Time.UTC().Format("2006-01-02T15:04:05.000Z") // UTC Event log

						if reqBytes, err := json.Marshal(req); err == nil && reqBytes != nil {
							if e.Level == logrus.WarnLevel || e.Level == logrus.ErrorLevel || e.Level == logrus.FatalLevel || e.Level == logrus.PanicLevel {
								reqStr += header + bold + mono + string(reqBytes) + mono + bold
							} else {
								reqStr += header + mono + string(reqBytes) + mono + "\n"
							}
						} else {
							// If couldn't process message; save the error, so the error can be reported instead
							reqStr += header + bold + mono + "SlackHook; forwardFunc; JSON Marshal log entry; Err: " + err.Error() + mono + bold + "\n"
						}
						if len(reqStr) > MAX_PAYLOAD_SIZE {
							if err = send(reqStr); err != nil {
								reqStr = header + bold + mono +
									"forwardFunc; Message of length " + strconv.FormatInt(int64(len(reqStr)), 10) + " dropped because of error " + err.Error() +
									mono + bold + "\n"
							} else {
								reqStr = ""
							}
						}
					}
					if len(reqStr) > 0 { // send concatenated string if hasn't been sent
						if err = send(reqStr); err != nil { // forwarding this message is delayed as the reqStr can get sent only the next time readFromStack returns entries
							header = "[Common][error][ICON-BRIDGE][" + time.Now().UTC().Format("2006-01-02T15:04:05.000Z") + "]\n"
							reqStr = header + bold + mono + "forwardFunc; Message of length " + strconv.FormatInt(int64(len(reqStr)), 10) + " dropped because of error " + err.Error() + mono + bold + "\n"
						} else {
							reqStr = ""
						}
					}
				}
			}
		}
	}()

	return
}
