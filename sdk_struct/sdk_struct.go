package sdkstruct

var SvrConf IMConfig

type IMConfig struct {
	Platform      int32  `json:"platform"`
	ApiAddr       string `json:"api_addr"`
	WsAddr        string `json:"ws_addr"`
	DataDir       string `json:"data_dir"`
	LogLevel      uint32 `json:"log_level"`
	ObjectStorage string `json:"object_storage"` //"cos"(default)  "oss"
}
