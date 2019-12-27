package logger

// 定义格式函数
type FormatFunc func(LogType, interface{}) (string, []interface{}, bool)

// 定义日志接口
type ILogger interface {
	SetLogLevel(LogType)        // 设置日志级别
	GetLogLevel() LogType       // 获取日志级别
	SetLoggerFormat(FormatFunc) // 设置日志格式
	Debug(interface{})          // 打印debug日志
	Info(interface{})           // 打印info日志
	Notice(interface{})         // 打印notice日志
	Warn(interface{})           // 打印warn日志
	Error(interface{})          // 打印error日志
	Fatal(interface{})          // 打印fatal日志
}
