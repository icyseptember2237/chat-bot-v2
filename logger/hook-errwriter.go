package logger

import (
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"io"
)

type LfsHook struct {
	*lfshook.LfsHook
	writer io.Writer
}

// 将等级为error(及以上)的日志复制一份写到errWriter。
func NewErrWriterHook(errWriter io.Writer) *LfsHook {
	lfsh := NewXdLfsHook(
		lfshook.WriterMap{
			ErrorLevel: errWriter,
			FatalLevel: errWriter,
			PanicLevel: errWriter,
		}, newJSONFormatter())
	lfsh.writer = errWriter
	return lfsh
}

func NewXdLfsHook(output interface{}, formatter logrus.Formatter) *LfsHook {
	return &LfsHook{
		LfsHook: lfshook.NewHook(output, formatter),
	}
}

func (hook *LfsHook) Fire(entry *logrus.Entry) error {
	return hook.LfsHook.Fire(entry)
}
