package callback

// 关于conn连接callback
type OnConnListener interface {
	OnConnecting()
	OnConnectSuccess()
	OnConnectFailed(errCode int32, errMsg string)
	OnUserTokenExpired()
}
