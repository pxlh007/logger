package logger

type FormatFunc func( LogType, interface{}) (string, []interface{}, bool)

// 定义日志接口
type ILogger interface {
	Init()
	SetLogLevel(LogType)
	GetLogLevel()
	SetLoggerFormat(FormatFunc)
	Debug(interface{})
	Info(interface{})
	Notice(interface{})
	Warn(interface{})
	Error(interface{})
	Fatal(interface{})
}
