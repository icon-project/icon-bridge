package icon

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"sync"

	"github.com/icon-project/icon-bridge/common/log"

	"github.com/gorilla/websocket"
	"github.com/icon-project/icon-bridge/common/jsonrpc"
)

type client struct {
	*jsonrpc.Client
	conns map[string]*websocket.Conn
	log   log.Logger
	mtx   sync.Mutex
}

func (c *client) MonitorEvent(ctx context.Context, p *EventRequest, cb func(conn *websocket.Conn, v *EventNotification) error, errCb func(*websocket.Conn, error)) error {
	resp := &EventNotification{}
	return c.Monitor(ctx, "/event", p, resp, func(conn *websocket.Conn, v interface{}) error {
		switch t := v.(type) {
		case *EventNotification:
			if err := cb(conn, t); err != nil {
				c.log.Debugf("MonitorEvent callback return err:%+v", err)
			}
		case WSEvent:
			c.log.Debugf("MonitorBlock WSEvent %s %+v", conn.LocalAddr().String(), t)
			switch t {
			case WSEventInit:
				c.log.Debug("WSEventInit")
			default:
				errCb(conn, fmt.Errorf("not supported type %T", t))
			}
		case error:
			errCb(conn, t)
		default:
			errCb(conn, fmt.Errorf("not supported type %T", t))
		}
		return nil
	})
}

func (c *client) Monitor(ctx context.Context, reqUrl string, reqPtr, respPtr interface{}, cb wsReadCallback) error {
	if cb == nil {
		return fmt.Errorf("callback function cannot be nil")
	}
	conn, err := c.wsConnect(reqUrl, nil)
	if err != nil {
		return ErrConnectFail
	}
	defer func() {
		c.log.Debugf("Monitor finish %s", conn.LocalAddr().String())
		c.wsClose(conn)
	}()
	if err = c.wsRequest(conn, reqPtr); err != nil {
		return err
	}
	if err := cb(conn, WSEventInit); err != nil {
		return err
	}
	return c.wsReadJSONLoop(ctx, conn, respPtr, cb)
}

func (c *client) CloseMonitor(conn *websocket.Conn) {
	c.log.Debugf("CloseMonitor %s", conn.LocalAddr().String())
	c.wsClose(conn)
}

func (c *client) CloseAllMonitor() {
	for _, conn := range c.conns {
		c.log.Debugf("CloseAllMonitor %s", conn.LocalAddr().String())
		c.wsClose(conn)
	}
}

func (c *client) _addWsConn(conn *websocket.Conn) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	la := conn.LocalAddr().String()
	c.conns[la] = conn
}

func (c *client) _removeWsConn(conn *websocket.Conn) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	la := conn.LocalAddr().String()
	_, ok := c.conns[la]
	if ok {
		delete(c.conns, la)
	}
}

func (c *client) wsConnect(reqUrl string, reqHeader http.Header) (*websocket.Conn, error) {
	wsEndpoint := strings.Replace(c.Endpoint, "http", "ws", 1)
	conn, httpResp, err := websocket.DefaultDialer.Dial(wsEndpoint+reqUrl, reqHeader)
	if err != nil {
		wsErr := wsConnectError{error: err}
		wsErr.httpResp = httpResp
		return nil, wsErr
	}
	c._addWsConn(conn)
	return conn, nil
}

func (c *client) wsRequest(conn *websocket.Conn, reqPtr interface{}) error {
	if reqPtr == nil {
		log.Panicf("reqPtr cannot be nil")
	}
	var err error
	wsResp := &WSResponse{}
	if err = conn.WriteJSON(reqPtr); err != nil {
		return wsRequestError{fmt.Errorf("fail to WriteJSON err:%+v", err), nil}
	}

	if err = conn.ReadJSON(wsResp); err != nil {
		return wsRequestError{fmt.Errorf("fail to ReadJSON err:%+v", err), nil}
	}

	if wsResp.Code != 0 {
		return wsRequestError{
			fmt.Errorf("invalid WSResponse code:%d, message:%s", wsResp.Code, wsResp.Message),
			wsResp}
	}
	return nil
}

func (c *client) wsClose(conn *websocket.Conn) {
	c._removeWsConn(conn)
	if err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
		c.log.Debugf("fail to WriteMessage CloseNormalClosure err:%+v", err)
	}
	if err := conn.Close(); err != nil {
		c.log.Debugf("fail to Close err:%+v", err)
	}
}

func (c *client) wsRead(conn *websocket.Conn, respPtr interface{}) error {
	mt, r, err := conn.NextReader()
	if err != nil {
		return err
	}
	if mt == websocket.CloseMessage {
		return io.EOF
	}
	return json.NewDecoder(r).Decode(respPtr)
}

func (c *client) wsReadJSONLoop(ctx context.Context, conn *websocket.Conn, respPtr interface{}, cb wsReadCallback) error {
	elem := reflect.ValueOf(respPtr).Elem()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			v := reflect.New(elem.Type())
			ptr := v.Interface()
			if _, ok := c.conns[conn.LocalAddr().String()]; !ok {
				c.log.Debugf("wsReadJSONLoop c.conns[%s] is nil", conn.LocalAddr().String())
				return errors.New("wsReadJSONLoop c.conns is nil")
			}
			if err := c.wsRead(conn, ptr); err != nil {
				c.log.Debugf("wsReadJSONLoop c.conns[%s] ReadJSON err:%+v", conn.LocalAddr().String(), err)
				if cErr, ok := err.(*websocket.CloseError); !ok || cErr.Code != websocket.CloseNormalClosure {
					cb(conn, err)
				}
				return err
			}
			if err := cb(conn, ptr); err != nil {
				return err
			}
		}

	}
}

func newClient(uri string, l log.Logger) *client {
	//TODO options {MaxRetrySendTx, MaxRetryGetResult, MaxIdleConnsPerHost, Debug, Dump}
	tr := &http.Transport{MaxIdleConnsPerHost: 1000}
	c := &client{
		Client: jsonrpc.NewJsonRpcClient(&http.Client{Transport: tr}, uri),
		conns:  make(map[string]*websocket.Conn),
		log:    l,
	}
	opts := IconOptions{}
	opts.SetBool(IconOptionsDebug, true)
	c.CustomHeader[HeaderKeyIconOptions] = opts.ToHeaderValue()
	return c
}
