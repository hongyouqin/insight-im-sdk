package log

import (
	"runtime"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

type fileHook struct{}

func newFileHook() *fileHook {
	return &fileHook{}
}

func (f *fileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (f *fileHook) Fire(entry *logrus.Entry) error {
	var s string
	_, b, c, _ := runtime.Caller(8)
	i := strings.LastIndex(b, "/")
	if i != -1 {
		l := len(b)
		s = b[i+1:l] + ":" + strconv.FormatInt(int64(c), 10)
	}
	entry.Data["FilePath"] = s
	return nil
}
