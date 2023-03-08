package innet

import (
	"errors"
	"insight/insight-im-sdk/pkg/common"
	"time"

	"github.com/emicklei/proto"
	"github.com/gorilla/websocket"
)

type WSClient struct {
	*WSConn
	*WsRespAsyn
	pushMsgCh chan common.Cmd2Value // recv push msg --> channel
}

func NewWsClient(respAsync *WsRespAsyn, conn *WSConn, pushMsg chan common.Cmd2Value) *WSClient {
	client := WSClient{
		WSConn:     conn,
		WsRespAsyn: respAsync,
		pushMsgCh:  pushMsg,
	}
	return &client
}

func (c *WSClient) SendReqWaitResp(msg proto.Message, identifier int32, senderId string, operationId string) (*WsResp, error) {
	msgIncr, ch, err := c.AddNotify(senderId)
	if err != nil {
		return nil, err
	}
	defer c.DelNotify(msgIncr)

	var req WsReq
	req.ReqIdentifier = identifier
	req.OperationID = operationId
	req.MsgIncr = msgIncr
	req.Data, err = proto.Marshal(msg)
	if err != nil {
		return nil, err
	}
	writeMsg := func() (*websocket.Conn, bool) {
		const retryTimes = 3 //重试次数
		for i := 0; i < retryTimes; i++ {
			conn, err := c.writeBinaryMsg(req)
			if err != nil {
				if !c.IsWriteTimeout(err) {
					time.Sleep(time.Duration(1) * time.Second)
					continue
				} else {
					return nil, false
				}
			}
			return conn, true
		}
		return nil, false
	}
	_, success := writeMsg()
	if success {
		return c.WaitResp(ch)
	}
	return nil, errors.New("SendReqWaitResp failed")
}

func (c *WSClient) WaitResp(ch chan WsResp) (*WsResp, error) {
	select {
	case resp := <-ch:
		if resp.ErrCode != 0 {
			return nil, errors.New("recv code err")
		} else {
			return &resp, nil
		}
	case <-time.After(time.Second * 1):
		return nil, errors.New("recv resp msg timeout")
	}
}

// 接收来自服务端的消息
func (c *WSClient) recvMessage() {
	for {

	}
}
