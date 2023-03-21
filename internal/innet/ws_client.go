package innet

import (
	"errors"
	"fmt"
	"insight/insight-im-sdk/pkg/common"
	"insight/insight-im-sdk/pkg/constant"
	"insight/insight-im-sdk/pkg/proto/msg"
	sdkstruct "insight/insight-im-sdk/sdk_struct"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
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
					//	log.Info("write msg retry")
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
		if c.WSConn.conn == nil {
			_, err := c.WSConn.ReConnect()
			if err != nil {
				time.Sleep(time.Duration(1) * time.Second)
			}
			continue
		}
		msgType, msg, err := c.WSConn.conn.ReadMessage()
		if err != nil {
			continue
		}
		if msgType == websocket.BinaryMessage {
			// 处理消息
			go c.doMsg(msg)
		} else {
			//其他消息
		}
	}
}

func (c *WSClient) doMsg(message []byte) {
	resp, err := c.decodeMsg(message)
	if err != nil {
		return
	}
	switch resp.ReqIdentifier {
	case constant.WSGetNewestSeq:
		c.doGetNewestSeq(*resp)
	case constant.WSPullMsgBySeqList:
		c.doPullMsgBySeqList(*resp)
	case constant.WSPushMsg:
		c.doPushMsg(*resp)
	case constant.WSSendMsg:
		c.doSendMsg(*resp)
	default:
		//其他情况
		fmt.Println("undefined identifier")
	}
}

func (c *WSClient) doSendMsg(resp WsResp) error {
	if err := c.NotifyResp(resp); err != nil {
		return err
	}
	return nil
}

func (c *WSClient) doPushMsg(resp WsResp) error {
	if resp.ErrCode != 0 {
		return errors.New(resp.ErrMsg)
	}
	//推送到消息同步(msg_sync)处理
	var data msg.MsgData
	err := proto.Unmarshal(resp.Data, &data)
	if err != nil {
		return err
	}
	return common.AddPushMsgTask(&sdkstruct.CmdPushMsg{Msg: &data, OperationID: 0}, c.pushMsgCh)
}

func (c *WSClient) doGetNewestSeq(resp WsResp) error {
	if err := c.NotifyResp(resp); err != nil {
		return err
	}
	return nil
}

func (c *WSClient) doPullMsgBySeqList(resp WsResp) error {
	if err := c.NotifyResp(resp); err != nil {
		return err
	}
	return nil
}
