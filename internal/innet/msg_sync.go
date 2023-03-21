package innet

import (
	"insight/insight-im-sdk/pkg/common"
	"insight/insight-im-sdk/pkg/constant"
	"insight/insight-im-sdk/pkg/db"
	"insight/insight-im-sdk/pkg/log"
	"insight/insight-im-sdk/pkg/proto/msg"
	sdkstruct "insight/insight-im-sdk/sdk_struct"
)

// 消息同步
type MsgSync struct {
	*db.DataBase
	*WSClient
	PushMsgAndMaxSeqCh chan common.Cmd2Value
	userId             string
	seqMaxSync         uint32
	seqMaxNeedSync     uint32 //需要同步的
}

// ----------- 实现Task interface -------------
func (s *MsgSync) Consume(cmd common.Cmd2Value) {
	switch cmd.Cmd {
	case constant.CmdPushMsg:
	case constant.CmdMaxSeq:
	}
}

func (s *MsgSync) Product() chan common.Cmd2Value {
	return s.PushMsgAndMaxSeqCh
}

// ----------- task interface end -----------

func (s *MsgSync) doPushMsg(cmd common.Cmd2Value) {
	if s.seqMaxNeedSync == 0 {
		return
	}
	CmdPushMsg := cmd.Value.(sdkstruct.CmdPushMsg)
	m := CmdPushMsg.Msg
	operationID := CmdPushMsg.Platform
	if m.Seq == 0 {
		//触发新消息
		s.NewMsgCome([]*msg.MsgData{m}, operationID)
		return
	}
	if m.Seq == s.seqMaxSync+1 {
		//触发新消息
		s.NewMsgCome([]*msg.MsgData{m}, operationID)
		s.seqMaxSync = m.Seq
	}
	if m.Seq > s.seqMaxNeedSync {
		s.seqMaxNeedSync = m.Seq
	}
	s.syncMsg()
}

func (s *MsgSync) doMaxSeq(cmd common.Cmd2Value) {

}

func (s *MsgSync) syncMsg() {
	if s.seqMaxNeedSync <= s.seqMaxSync {
		return
	}

	s.seqMaxSync = s.seqMaxNeedSync
}

// 新消息来了
func (s *MsgSync) NewMsgCome(msgs []*msg.MsgData, operationID string) {
	err := common.AddNewMsgCome(&sdkstruct.CmdNewMsgCome{msgs, operationID}, s.PushMsgAndMaxSeqCh)
	if err != nil {
		log.Warn(operationID, "Add NewMsgCome failed", err.Error(), s.userId)
	}
}
