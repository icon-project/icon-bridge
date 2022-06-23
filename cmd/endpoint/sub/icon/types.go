package icon

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

var (
	ErrConnectFail            = fmt.Errorf("fail to connect")
	ErrSendFailByExpired      = fmt.Errorf("reject by expired")
	ErrSendFailByFuture       = fmt.Errorf("reject by future")
	ErrSendFailByOverflow     = fmt.Errorf("reject by overflow")
	ErrGetResultFailByPending = fmt.Errorf("fail to getresult by pending")
)

//T_BIN_DATA, T_HASH
type HexBytes string

func (hs HexBytes) Value() ([]byte, error) {
	if hs == "" {
		return nil, nil
	}
	return hex.DecodeString(string(hs[2:]))
}
func NewHexBytes(b []byte) HexBytes {
	return HexBytes("0x" + hex.EncodeToString(b))
}

//T_INT
type HexInt string

func (i HexInt) Value() (int64, error) {
	s := string(i)
	if strings.HasPrefix(s, "0x") {
		s = s[2:]
	}
	return strconv.ParseInt(s, 16, 64)
}

func (i HexInt) Int() (int, error) {
	s := string(i)
	if strings.HasPrefix(s, "0x") {
		s = s[2:]
	}
	v, err := strconv.ParseInt(s, 16, 32)
	return int(v), err
}

func NewHexInt(v int64) HexInt {
	return HexInt("0x" + strconv.FormatInt(v, 16))
}

type Address string

func (a Address) Value() ([]byte, error) {
	var b [21]byte
	switch a[:2] {
	case "cx":
		b[0] = 1
	case "hx":
	default:
		return nil, fmt.Errorf("invalid prefix %s", a[:2])
	}
	n, err := hex.Decode(b[1:], []byte(a[2:]))
	if err != nil {
		return nil, err
	}
	if n != 20 {
		return nil, fmt.Errorf("invalid length %d", n)
	}
	return b[:], nil
}

func NewAddress(b []byte) Address {
	if len(b) != 21 {
		return ""
	}
	switch b[0] {
	case 1:
		return Address("cx" + hex.EncodeToString(b[1:]))
	case 0:
		return Address("hx" + hex.EncodeToString(b[1:]))
	default:
		return ""
	}
}

type wsConnectError struct {
	error
	httpResp *http.Response
}

type wsReadCallback func(*websocket.Conn, interface{}) error

type wsRequestError struct {
	error
	wsResp *WSResponse
}

type BlockRequest struct {
	Height       HexInt         `json:"height"`
	EventFilters []*EventFilter `json:"eventFilters,omitempty"`
}
type EventNotification struct {
	Hash   HexBytes   `json:"hash"`
	Height HexInt     `json:"height"`
	Index  HexInt     `json:"index"`
	Events []HexInt   `json:"events,omitempty"`
	Logs   []EventLog `json:"logs,omitempty"`
}

type EventLog struct {
	Addr    string        `json:"scoreAddress"`
	Indexed []interface{} `json:"indexed"`
	Data    []interface{} `json:"data"`
}

type EventRequest struct {
	EventFilter
	Height HexInt `json:"height"`
	Logs   HexInt `json:"logs,omitempty"`
}

type EventFilter struct {
	Addr      Address   `json:"addr,omitempty"`
	Signature string    `json:"event"`
	Indexed   []*string `json:"indexed,omitempty"`
	Data      []*string `json:"data,omitempty"`
}

type WSEvent string

const (
	WSEventInit WSEvent = "WSEventInit"
)

type WSResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}

const (
	HeaderKeyIconOptions = "Icon-Options"
	IconOptionsDebug     = "debug"
	IconOptionsTimeout   = "timeout"
)

type IconOptions map[string]string

func (opts IconOptions) Set(key, value string) {
	opts[key] = value
}

func (opts IconOptions) Get(key string) string {
	if opts == nil {
		return ""
	}
	v := opts[key]
	if len(v) == 0 {
		return ""
	}
	return v
}

func (opts IconOptions) Del(key string) {
	delete(opts, key)
}

func (opts IconOptions) SetBool(key string, value bool) {
	opts.Set(key, strconv.FormatBool(value))
}

func (opts IconOptions) GetBool(key string) (bool, error) {
	return strconv.ParseBool(opts.Get(key))
}

func (opts IconOptions) SetInt(key string, v int64) {
	opts.Set(key, strconv.FormatInt(v, 10))
}

func (opts IconOptions) GetInt(key string) (int64, error) {
	return strconv.ParseInt(opts.Get(key), 10, 64)
}

func (opts IconOptions) ToHeaderValue() string {
	if opts == nil {
		return ""
	}
	strs := make([]string, len(opts))
	i := 0
	for k, v := range opts {
		strs[i] = fmt.Sprintf("%s=%s", k, v)
		i++
	}
	return strings.Join(strs, ",")
}

func NewIconOptionsByHeader(h http.Header) IconOptions {
	s := h.Get(HeaderKeyIconOptions)
	if s != "" {
		kvs := strings.Split(s, ",")
		m := make(map[string]string)
		for _, kv := range kvs {
			if kv != "" {
				idx := strings.Index(kv, "=")
				if idx > 0 {
					m[kv[:idx]] = kv[(idx + 1):]
				} else {
					m[kv] = ""
				}
			}
		}
		return m
	}
	return nil
}
