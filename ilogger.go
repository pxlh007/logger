package logger

type FormatFunc func( LogType, interface{}) (string, []interface{}, bool)

// 定义日志接口
type ILogger interface {
	func Init()
	func SetLogLevel(LogType)
	func GetLogLevel()
	func SetLoggerFormat(FormatFunc)
	func Debug(interface{})
	func Info(interface{})
	func Notice(interface{})
	func Warn(interface{})
	func Error(interface{})
	func Fatal(interface{})
}
