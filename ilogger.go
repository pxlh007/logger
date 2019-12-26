package logger

type FormatFunc func(logType LogType, i interface{}) (string, []interface{}, bool)

// 定义日志接口
type ILogger interface {
	func Init()
	func SetLogLevel(ltyp LogType)
	func GetLogLevel()
	func SetLoggerFormat(func FormatFunc)
	func Debug(i interface{})
	func Info(i interface{})
	func Notice(i interface{})
	func Warn(i interface{})
	func Error(i interface{})
	func Fatal(i interface{})
}
