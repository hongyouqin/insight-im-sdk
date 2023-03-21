package sdkstruct

import "insight/insight-im-sdk/pkg/proto/msg"

var SvrConf IMConfig

type IMConfig struct {
	Platform      int32  `json:"platform"`
	ApiAddr       string `json:"api_addr"`
	WsAddr        string `json:"ws_addr"`
	DataDir       string `json:"data_dir"`
	LogLevel      uint32 `json:"log_level"`
	ObjectStorage string `json:"object_storage"` //"cos"(default)  "oss"
}

// 推送命令消息
type CmdPushMsg struct {
	Msg         *msg.MsgData
	OperationID string
}

// 新消息命令
type CmdNewMsgCome struct {
	MsgList     []*msg.MsgData
	OperationID string
}
