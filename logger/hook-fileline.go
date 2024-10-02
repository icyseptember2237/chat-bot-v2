package logger

import (
	"fmt"
	"path"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

// 在日志中输出文件名和行号，实际消耗要比想象中小得多。
const FileKey = "file"
const LineKey = "line"
const FuncKey = "func"

type FileLineHook struct {
}

func NewFileLineHook() *FileLineHook {
	return &FileLineHook{}
}

func searchFileLine() (string, int, string) {
	for skip := 6; skip < 18; skip++ {
		pc, file, line, _ := runtime.Caller(skip)
		if pc == 0 {
			break
		}

		f := runtime.FuncForPC(pc)
		if f == nil {
			break
		}

		fn := f.Name()
		if !strings.Contains(fn, "/logger") {
			return file, line, fn
		}
	}

	return "", 0, ""
}

func (h *FileLineHook) Fire(entry *logrus.Entry) error {
	file, line, fn := searchFileLine()
	if file != "" && line != 0 {
		if _, ok := entry.Data[FileKey]; !ok {
			entry.Data[FileKey] = fmt.Sprintf("%v:%v", path.Base(file), line)
		}

		if _, ok := entry.Data[FuncKey]; !ok {
			idx := strings.LastIndex(fn, "/")
			entry.Data[FuncKey] = fn[idx+1:]
		}
	}

	return nil
}

func (h *FileLineHook) Levels() []Level {
	return AllLevels
}
