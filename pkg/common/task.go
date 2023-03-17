package common

import (
	"errors"
	"insight/insight-im-sdk/pkg/constant"
	sdkstruct "insight/insight-im-sdk/sdk_struct"
	"time"
)

type Cmd2Value struct {
	Cmd   string
	Value interface{}
}

type Task interface {
	Consume(cmd Cmd2Value)   //消费
	Product() chan Cmd2Value //生产
}

// 处理任务
func ProcessTask(task Task) {
	for {
		select {
		case cmd := <-task.Product():
			task.Consume(cmd)
		}
	}
}

// 以下是自定义好的任务
func AddPushMsgTask(c *sdkstruct.CmdPushMsg, ch chan Cmd2Value) error {
	c2v := Cmd2Value{Cmd: constant.CmdPushMsg, Value: c}
	return sendTask(ch, c2v, 1)
}

func sendTask(ch chan Cmd2Value, value Cmd2Value, timeout int64) error {
	var flag = 0
	select {
	case ch <- value:
		flag = 1
	case <-time.After(time.Second * time.Duration(timeout)):
		flag = 2
	}
	if flag == 1 {
		return nil
	} else {
		return errors.New("send cmd timeout")
	}
}
