package innet

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"insight/insight-im-sdk/callback"
	sdkstruct "insight/insight-im-sdk/sdk_struct"
	"insight/pkg/utils"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	MaxLimitMsgLen     = 51200 //最大消息长度
	WriteTimeoutSecond = 30
)

// 客户端连接websocket服务器
type WSConn struct {
	mtx      sync.Mutex
	conn     *websocket.Conn
	token    string
	userId   string
	listener callback.OnConnListener //连接监听
}

func NewConn(listener callback.OnConnListener, token string) *WSConn {
	c := WSConn{
		listener: listener,
		token:    token,
	}
	conn, _ := c.ReConnect()
	c.conn = conn
	return &c
}

func (w *WSConn) ReConnect() (*websocket.Conn, error) {
	w.mtx.Lock()
	defer w.mtx.Unlock()
	if w.conn != nil {
		w.conn.Close()
		w.conn = nil
	}
	w.listener.OnConnecting()
	operationID := utils.OperationIDGenerator()
	url := fmt.Sprintf("%s?sendID=%s&token=%s&platformID=%d&operationID=%s", sdkstruct.SvrConf.WsAddr, w.userId, w.token, sdkstruct.SvrConf.Platform, operationID)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		w.listener.OnConnectFailed(1001, err.Error())
		return nil, err
	}
	w.listener.OnConnectSuccess()
	w.conn = conn
	return conn, nil
}

func (w *WSConn) CloseConn() error {
	w.mtx.Lock()
	defer w.mtx.Unlock()
	if w.conn != nil {
		err := w.conn.Close()
		if err != nil {
			return err
		}
		w.conn = nil
	}
	return nil
}

func (w *WSConn) writeBinaryMsg(msg WsReq) (*websocket.Conn, error) {
	w.mtx.Lock()
	defer w.mtx.Unlock()
	if w.conn == nil {
		return nil, errors.New("conn is nil")
	}
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(msg)
	if err != nil {
		return nil, errors.New("Encode error")
	}
	err = w.SetWriteTimeout(WriteTimeoutSecond)
	if err != nil {
		return nil, errors.New("SetWriteTimeout")
	}
	if len(buff.Bytes()) > MaxLimitMsgLen {
		return nil, errors.New("msg too long")
	}
	return w.conn, w.conn.WriteMessage(websocket.BinaryMessage, buff.Bytes())
}

func (w *WSConn) SendPingMsg() error {
	w.mtx.Lock()
	defer w.mtx.Unlock()
	if w.conn == nil {
		return errors.New("conn is nil")
	}
	ping := "try ping"
	err := w.SetWriteTimeout(WriteTimeoutSecond)
	if err != nil {
		return err
	}
	err = w.conn.WriteMessage(websocket.PingMessage, []byte(ping))
	if err != nil {
		return err
	}
	return nil
}

func (w *WSConn) SetWriteTimeout(timeout int) error {
	return w.conn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
}

func (w *WSConn) IsWriteTimeout(err error) bool {
	if strings.Contains(err.Error(), "timeout") {
		return true
	}
	return false
}

func (w *WSConn) decodeMsg(msg []byte) (*WsResp, error) {
	buff := bytes.NewBuffer(msg)
	dec := gob.NewDecoder(buff)
	var data WsResp
	err := dec.Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
