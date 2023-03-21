package constant

const (
	//Websocket Protocol
	WSGetNewestSeq     = 1001
	WSPullMsgBySeqList = 1002
	WSSendMsg          = 1003
	WSHeartbeat        = 1004
	WSSendSignalMsg    = 1005
	WSPushMsg          = 2001
	WSKickOnlineMsg    = 2002
	WsLogoutMsg        = 2003
	WSDataError        = 3001
)

// cmd
const (
	CmdNewMsgCome = "005"

	CmdMaxSeq  = "maxSeq"
	CmdPushMsg = "pushMsg"
)
