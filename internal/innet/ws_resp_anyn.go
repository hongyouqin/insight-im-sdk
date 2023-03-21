package innet

import (
	"errors"
	"insight/insight-im-sdk/pkg/utils"
	"sync"
	"time"
)

type WsResp struct {
	ReqIdentifier int    `json:"reqIdentifier"`
	ErrCode       int    `json:"errCode"`
	ErrMsg        string `json:"errMsg"`
	MsgIncr       string `json:"msgIncr"`
	OperationID   string `json:"operationID"`
	Data          []byte `json:"data"`
}

type WsReq struct {
	ReqIdentifier int32  `json:"reqIdentifier"`
	Token         string `json:"token"`
	SendID        string `json:"sendID"`
	OperationID   string `json:"operationID"`
	MsgIncr       string `json:"msgIncr"`
	Data          []byte `json:"data"`
}

// 收到的消息，进行异步通知
type WsRespAsyn struct {
	notify map[string]chan WsResp //key 每一条消息产生一个唯一的id值，value是返回的消息提
	mtx    sync.RWMutex
}

func NewWsRespAsyn() *WsRespAsyn {
	return &WsRespAsyn{
		notify: make(map[string]chan WsResp, 1000),
	}
}

func (r *WsRespAsyn) AddNotify(userId string) (string, chan WsResp, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	msgIncr := utils.GenMsgIncr(userId)
	ch := make(chan WsResp, 1)
	if _, ok := r.notify[msgIncr]; ok {
		r.notify[msgIncr] = ch
		return msgIncr, ch, nil
	}
	close(ch)
	return msgIncr, nil, errors.New("msgIncr is exist")
}

func (r *WsRespAsyn) GetNotify(msgIncr string) chan WsResp {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	if ch, ok := r.notify[msgIncr]; ok {
		return ch
	}
	return nil
}

func (r *WsRespAsyn) DelNotify(msgIncr string) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	if ch, ok := r.notify[msgIncr]; ok {
		close(ch)
		delete(r.notify, msgIncr)
	}
}

// 获取socket返回的消息，通知出去
func (r *WsRespAsyn) NotifyResp(resp WsResp) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	ch := r.GetNotify(resp.MsgIncr)
	if ch == nil {
		return nil
	}
	return notifyCh(ch, resp, 1)
}

func notifyCh(ch chan WsResp, value WsResp, timeout int64) error {
	var flag = false
	select {
	case ch <- value:
		flag = true
	case <-time.After(time.Second * time.Duration(timeout)):
		flag = false
	}
	if flag {
		return nil
	} else {
		return errors.New("send cmd timeout")
	}
}
