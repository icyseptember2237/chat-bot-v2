package logger

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
	"reflect"
)

var AppName string

func initLoggerWithOptions(l Logger, options Options) (err error) {
	AppName = options.AppName
	// 如果配置里指定了日志等级，则解析并设置，否则默认等级是info。
	if options.Level != "" {
		level, err := ParseLevel(options.Level)
		if err != nil {
			return errors.Wrapf(err, "failed to parse level(%s)", options.Level)
		}
		l.SetLevel(level)
	}

	// 设置默认值
	AppName = "service"
	file, errfile := "./log/service.app.log", "./log/service.err.log"

	if options.AppName != "" {
		AppName = options.AppName
	}

	// 如果配置里指定了日志文件，则解析并设置，否则默认写到stderr。
	if options.File != "" {
		file = options.File
	}
	err = handleFileOutput(l, file) // 设置output、压测标志
	if err != nil {
		return errors.Wrapf(err, "failed to set logger.Output")
	}

	l.ResetHooks()

	l.AddHook(NewFileLineHook()) // 在日志中输出文件名和行号。
	if options.Format == "json" || options.Format == "" {
		l.SetFormatter(newJSONFormatter())
	} else {
		l.SetFormatter(newTextFormatter())
	}

	// 如果配置里指定了错误日志文件，则额外将等级为error(及以上)的日志复制一份写到该文件中。
	if options.ErrFile != "" {
		errfile = options.ErrFile
	}
	errWriter, err := os.OpenFile(errfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return errors.Wrapf(err, "failed to open err file(%s)", options.ErrFile)
	}
	l.AddHook(NewErrWriterHook(errWriter))

	return
}

func handleFileOutput(l Logger, fileName string) error {
	writer, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return errors.Wrapf(err, "failed to open file(%s)", fileName)
	}
	l.SetOutput(writer)
	return nil
}

func parseFieldsFromObj(o interface{}) logrus.Fields {
	logFields := logrus.Fields{}

	val := reflect.ValueOf(o)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return logFields
		}
		val = val.Elem()
	}
	for i := 0; i < val.NumField(); i++ {
		fValue := val.Field(i)
		fType := val.Type().Field(i)
		if !isZero(fValue) && fValue.IsValid() && fType.PkgPath == "" { // exported fields
			if fValue.Kind() == reflect.Struct ||
				(fValue.Kind() == reflect.Ptr &&
					fValue.Elem().Kind() == reflect.Struct) {
				fields := parseFieldsFromObj(fValue.Interface())
				if fType.Anonymous {
					for k, v := range fields {
						logFields[k] = v
					}
				} else {
					logFields[fType.Name] = fields
				}
			} else {
				logFields[fType.Name] = fValue.Interface()
			}
		}
	}
	return logFields
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return len(v.String()) == 0
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	case reflect.Slice:
		return v.Len() == 0
	case reflect.Map:
		return v.Len() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Struct:
		return false
	}
	return false
}
