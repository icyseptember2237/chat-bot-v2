package logger

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type Entry = logrus.Entry
type Fields = logrus.Fields
type Formatter = logrus.Formatter
type Hook = logrus.Hook
type Logger = LogRusLogger
type Level = logrus.Level
type LevelHooks = logrus.LevelHooks

const PanicLevel = logrus.PanicLevel
const FatalLevel = logrus.FatalLevel
const ErrorLevel = logrus.ErrorLevel
const WarnLevel = logrus.WarnLevel
const InfoLevel = logrus.InfoLevel
const DebugLevel = logrus.DebugLevel
const TraceLevel = logrus.TraceLevel

var AllLevels = logrus.AllLevels
var StdLogger = StdLoggerNew()

func newJSONFormatter() logrus.Formatter {
	formatter := new(JSONFormatter)
	formatter.TimestampFormat = "2006-01-02T15:04:05.000Z07:00"
	return formatter
}

func newTextFormatter() logrus.Formatter {
	formatter := new(logrus.TextFormatter)
	formatter.TimestampFormat = "2006-01-02T15:04:05.000Z07:00"
	return formatter
}

// StdLoggerNew 生成带有标准格式的logger
func StdLoggerNew() Logger {
	formatter := newJSONFormatter()

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

func StandardLogger() Logger {
	return StdLogger
}

func SetOutput(out, shadowOut io.Writer) {
	StdLogger.SetOutput(out)
}

func GetOutput() (out io.Writer) {
	return StdLogger.GetOutput()
}

func SetFormatter(formatter Formatter) {
	StdLogger.SetFormatter(formatter)
}

func SetReportCaller(include bool) {
	StdLogger.SetReportCaller(include)
}

func SetLevel(level Level) {
	StdLogger.SetLevel(level)
}

func GetLevel() Level {
	return StdLogger.GetLevel()
}

func SetLevelWithShadow(level, shadowLevel Level) {
	StdLogger.SetLevel(level)
}

func AddHook(hook Hook) {
	StdLogger.AddHook(hook)
}

func ParseLevel(level string) (Level, error) {
	return logrus.ParseLevel(level)
}

func ParseLevelOrInfo(level string) Level {
	l, err := logrus.ParseLevel(level)
	if err != nil {
		return InfoLevel
	}
	return l
}

func NewLogrusEntry(l *logrus.Logger) *Entry {
	return logrus.NewEntry(l)
}

func WithError(err error) LoggerInterface {
	return StdLogger.WithError(err)
}

func WithField(key string, value interface{}) LoggerInterface {
	return StdLogger.WithField(key, value)
}

func WithFields(fields Fields) LoggerInterface {
	return StdLogger.WithFields(fields)
}

func WithTime(t time.Time) LoggerInterface {
	return StdLogger.WithTime(t)
}

func WithObject(obj interface{}) LoggerInterface {
	return StdLogger.WithObject(obj)
}

func Trace(ctx context.Context, args ...interface{}) {
	StdLogger.Trace(ctx, args...)
}

func Debug(ctx context.Context, args ...interface{}) {
	StdLogger.Debug(ctx, args...)
}

func Print(ctx context.Context, args ...interface{}) {
	StdLogger.Print(ctx, args...)
}

func Info(ctx context.Context, args ...interface{}) {
	StdLogger.Info(ctx, args...)
}

func Warn(ctx context.Context, args ...interface{}) {
	StdLogger.Warn(ctx, args...)
}

func Warning(ctx context.Context, args ...interface{}) {
	StdLogger.Warning(ctx, args...)
}

func Error(ctx context.Context, args ...interface{}) {
	StdLogger.Error(ctx, args...)
}

func Panic(ctx context.Context, args ...interface{}) {
	StdLogger.Panic(ctx, args...)
}

func Fatal(ctx context.Context, args ...interface{}) {
	StdLogger.Fatal(ctx, args...)
}

func Tracef(ctx context.Context, format string, args ...interface{}) {
	StdLogger.Tracef(ctx, format, args...)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	StdLogger.Debugf(ctx, format, args...)
}

func Printf(ctx context.Context, format string, args ...interface{}) {
	StdLogger.Printf(ctx, format, args...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	StdLogger.Infof(ctx, format, args...)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	StdLogger.Warnf(ctx, format, args...)
}

func Warningf(ctx context.Context, format string, args ...interface{}) {
	StdLogger.Warningf(ctx, format, args...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	StdLogger.Errorf(ctx, format, args...)
}

func Panicf(ctx context.Context, format string, args ...interface{}) {
	StdLogger.Panicf(ctx, format, args...)
}

func Fatalf(ctx context.Context, format string, args ...interface{}) {
	StdLogger.Fatalf(ctx, format, args...)
}

func Traceln(ctx context.Context, args ...interface{}) {
	StdLogger.Traceln(ctx, args...)
}

func Debugln(ctx context.Context, args ...interface{}) {
	StdLogger.Debugln(ctx, args...)
}

func Println(ctx context.Context, args ...interface{}) {
	StdLogger.Println(ctx, args...)
}

func Infoln(ctx context.Context, args ...interface{}) {
	StdLogger.Infoln(ctx, args...)
}

func Warnln(ctx context.Context, args ...interface{}) {
	StdLogger.Warnln(ctx, args...)
}

func Warningln(ctx context.Context, args ...interface{}) {
	StdLogger.Warningln(ctx, args...)
}

func Errorln(ctx context.Context, args ...interface{}) {
	StdLogger.Errorln(ctx, args...)
}

func Panicln(ctx context.Context, args ...interface{}) {
	StdLogger.Panicln(ctx, args...)
}

func Fatalln(ctx context.Context, args ...interface{}) {
	StdLogger.Fatalln(ctx, args...)
}
