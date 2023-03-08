package main

import (
	"context"
	"os"
	"runtime"

	"github.com/natefinch/lumberjack"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	fx.New(
		fx.Provide(),
		fx.Invoke(Server),
	).Run()
}

func Server(lc fx.Lifecycle, log *zap.Logger) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	lc.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				go func() {
					//启动服务
					//startRpc(log, chat)
				}()
				return nil
			},
			OnStop: func(context.Context) error {
				log.Info("server exiting")
				return nil
			},
		})
}

func newLogger() (*zap.Logger, error) {
	//return zap.NewProduction()
	//获取编码器,NewJSONEncoder()输出json格式，NewConsoleEncoder()输出普通文本格式
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder //指定时间格式
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	//文件writeSyncerí
	fileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./logs/sdk.log", //日志文件存放目录
		MaxSize:    10,               //文件大小限制,单位MB
		MaxBackups: 20,               //最大保留日志文件数量
		MaxAge:     30,               //日志文件保留天数
		Compress:   false,            //是否压缩处理
	})
	fileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(fileWriteSyncer, zapcore.AddSync(os.Stdout)), zapcore.DebugLevel) //第三个及之后的参数为写入文件的日志级别,ErrorLevel模式只记录error级别的日志

	logger := zap.New(fileCore, zap.AddCaller()) //AddCaller()为显示文件名和行号
	return logger, nil
}
