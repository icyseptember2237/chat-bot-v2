package logger

import (
	"context"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

type LogRusLogger interface {
	SetOutput(out io.Writer)
	GetOutput() (out io.Writer)
	SetFormatter(formatter Formatter)
	SetReportCaller(include bool)
	GetLevel() Level
	SetLevel(level Level)
	AddHook(hook Hook)
	ResetHooks()

	LoggerInterface
}

type CtxLogger struct {
	*logrus.Entry                // entry
	l             *logrus.Logger // logger
}

func NewCtxLogger() LogRusLogger {
	formatter := new(JSONFormatter)
	formatter.TimestampFormat = "2006-01-02T15:04:05.000Z07:00"

	l := logrus.Logger{
		Out:          os.Stderr,
		Formatter:    formatter,
		Hooks:        make(LevelHooks),
		Level:        InfoLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}

	return &CtxLogger{logrus.NewEntry(&l), &l}
}

func (cl *CtxLogger) SetOutput(out io.Writer) {
	cl.l.SetOutput(out)
}

func (cl *CtxLogger) GetOutput() (out io.Writer) {
	return cl.l.Out
}

func (cl *CtxLogger) SetFormatter(formatter Formatter) {
	cl.l.SetFormatter(formatter)
}

func (cl *CtxLogger) SetReportCaller(include bool) {
	cl.l.SetReportCaller(include)
}

func (cl *CtxLogger) GetLevel() Level {
	return cl.l.Level
}

func (cl *CtxLogger) SetLevel(level Level) {
	cl.l.SetLevel(level)
}

func (cl *CtxLogger) AddHook(hook Hook) {
	cl.l.AddHook(hook)
}

func (cl *CtxLogger) ResetHooks() {
	cl.l.ReplaceHooks(make(LevelHooks))
}

func (cl *CtxLogger) WithField(key string, value interface{}) LoggerInterface {
	// 借用logrus.Logger本身Entry的管理机制来创建Entry,下同
	return &CtxLogger{cl.l.WithField(key, value), cl.l}
}

func (cl *CtxLogger) WithFields(fields Fields) LoggerInterface {
	return &CtxLogger{cl.l.WithFields(fields), cl.l}
}

func (cl *CtxLogger) WithError(err error) LoggerInterface {
	return &CtxLogger{cl.l.WithError(err), cl.l}
}

func (cl *CtxLogger) WithTime(t time.Time) LoggerInterface {
	return &CtxLogger{cl.l.WithTime(t), cl.l}
}

func (cl *CtxLogger) WithObject(obj interface{}) LoggerInterface {
	fields := parseFieldsFromObj(obj)
	return &CtxLogger{cl.l.WithFields(fields), cl.l}
}

func logLevelTrans(originLevel Level) Level {
	if originLevel > logrus.InfoLevel {
		return logrus.InfoLevel
	}

	return originLevel
}

func (cl *CtxLogger) Tracef(ctx context.Context, format string, args ...interface{}) {
	cl.Entry.WithContext(ctx).Logf(logLevelTrans(TraceLevel), format, args...)
}

func (cl *CtxLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	cl.Entry.WithContext(ctx).Logf(logLevelTrans(DebugLevel), format, args...)
}

func (cl *CtxLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	cl.Entry.WithContext(ctx).Logf(logLevelTrans(InfoLevel), format, args...)
}

func (cl *CtxLogger) Printf(ctx context.Context, format string, args ...interface{}) {
	cl.Entry.WithContext(ctx).Printf(format, args...)
}

func (cl *CtxLogger) Warnf(ctx context.Context, format string, args ...interface{}) {
	cl.Entry.WithContext(ctx).Logf(logLevelTrans(WarnLevel), format, args...)
}

func (cl *CtxLogger) Warningf(ctx context.Context, format string, args ...interface{}) {
	cl.Entry.WithContext(ctx).Warnf(format, args...)
}

func (cl *CtxLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	cl.Entry.WithContext(ctx).Logf(ErrorLevel, format, args...)
}

func (cl *CtxLogger) Fatalf(ctx context.Context, format string, args ...interface{}) {
	cl.Entry.WithContext(ctx).Fatalf(format, args...)
}

func (cl *CtxLogger) Panicf(ctx context.Context, format string, args ...interface{}) {
	cl.Entry.WithContext(ctx).Logf(PanicLevel, format, args...)
}

func (cl *CtxLogger) Logf(ctx context.Context, level Level, format string, args ...interface{}) {
	cl.Entry.WithContext(ctx).Logf(logLevelTrans(level), format, args...)
}

func (cl *CtxLogger) Log(ctx context.Context, level Level, args ...interface{}) {
	cl.Entry.WithContext(ctx).Log(logLevelTrans(level), args...)
}

func (cl *CtxLogger) Trace(ctx context.Context, args ...interface{}) {
	cl.Entry.WithContext(ctx).Log(logLevelTrans(TraceLevel), args...)
}

func (cl *CtxLogger) Debug(ctx context.Context, args ...interface{}) {
	cl.Entry.WithContext(ctx).Log(logLevelTrans(DebugLevel), args...)
}

func (cl *CtxLogger) Info(ctx context.Context, args ...interface{}) {
	cl.Entry.WithContext(ctx).Log(InfoLevel, args...)
}

func (cl *CtxLogger) Print(ctx context.Context, args ...interface{}) {
	cl.Entry.WithContext(ctx).Print(args...)
}

func (cl *CtxLogger) Warn(ctx context.Context, args ...interface{}) {
	cl.Entry.WithContext(ctx).Log(WarnLevel, args...)
}

func (cl *CtxLogger) Warning(ctx context.Context, args ...interface{}) {
	cl.Entry.WithContext(ctx).Warn(args...)
}

func (cl *CtxLogger) Error(ctx context.Context, args ...interface{}) {
	cl.Entry.WithContext(ctx).Log(ErrorLevel, args...)
}

func (cl *CtxLogger) Fatal(ctx context.Context, args ...interface{}) {
	cl.Entry.WithContext(ctx).Fatal(args...)
}

func (cl *CtxLogger) Panic(ctx context.Context, args ...interface{}) {
	cl.Entry.WithContext(ctx).Panic(args...)
}

func (cl *CtxLogger) Logln(ctx context.Context, level Level, args ...interface{}) {
	cl.Entry.WithContext(ctx).Logln(logLevelTrans(level), args...)
}

func (cl *CtxLogger) Traceln(ctx context.Context, args ...interface{}) {
	cl.Entry.WithContext(ctx).Logln(logLevelTrans(TraceLevel), args...)
}

func (cl *CtxLogger) Debugln(ctx context.Context, args ...interface{}) {
	cl.WithContext(ctx).Logln(logLevelTrans(DebugLevel), args...)
}

func (cl *CtxLogger) Infoln(ctx context.Context, args ...interface{}) {
	cl.Entry.WithContext(ctx).Logln(InfoLevel, args...)
}

func (cl *CtxLogger) Println(ctx context.Context, args ...interface{}) {
	cl.Entry.WithContext(ctx).Println(args...)
}

func (cl *CtxLogger) Warnln(ctx context.Context, args ...interface{}) {
	cl.Entry.WithContext(ctx).Logln(WarnLevel, args...)
}

func (cl *CtxLogger) Warningln(ctx context.Context, args ...interface{}) {
	cl.Entry.WithContext(ctx).Logln(WarnLevel, args...)
}

func (cl *CtxLogger) Errorln(ctx context.Context, args ...interface{}) {
	cl.Entry.WithContext(ctx).Logln(ErrorLevel, args...)
}

func (cl *CtxLogger) Fatalln(ctx context.Context, args ...interface{}) {
	cl.Entry.WithContext(ctx).Fatalln(args...)
}

func (cl *CtxLogger) Panicln(ctx context.Context, args ...interface{}) {
	cl.Entry.WithContext(ctx).Logln(PanicLevel, args...)
}
