package logger

import (
    "sync"
    "io"
    "fmt"
    "time"
    "os"
    "strings"
)

type LogType int

const (
    DEBUG =     LogType(0)
    INFO =      LogType(1)
    NOTICE =    LogType(2)
    WARN =      LogType(3)
    ERROR =     LogType(4)
    CRITICAL =  LogType(5)
    FATAL =     LogType(6)
)

var logTypeStrings = func() []string{
    // log 类型对应的 名称字符串，用于输出，所以统一了长度，故 DEBUG 为 "DEBUG..." 和 "CRITICAL"等长
    types :=[]string{"DEBUG", "INFO", "NOTICE", "WARN", "ERROR", "CRITICAL", "FATAL"}
    maxTypeLen := 0
    for _, t := range types{
        if len(t) > maxTypeLen {
            maxTypeLen = len(t)
        }
    }
    for index, t := range types{
        typeLen := len(t)
        if typeLen < maxTypeLen{
            types[index] += strings.Repeat(".", maxTypeLen - typeLen)
        }
    }
    return types
}()

var logTypesColors = []string{"0;35", "1;36", "1;37", "0;33", "1;31", "1;31", "1;31"}

/*
    创建 Logger
 */
func NewLogger()  *Logger{
    logger := &Logger{}
    logger.Init()
    return logger
}

func GetLogTypeString(t LogType)  string{
    return logTypeStrings[t]
}

/*
    Logger 日志

    l := NewLogger()
    l.Info("hello")
    l.Warn(1)

    输入格式可以通过 SetLoggerFormat 设置。默认输出格式定义在 见Logger的DefaultLogFormatFunc

    可以通过 SetLogLevel 设置输出等级。
 */
type Logger struct {
    mu              sync.Mutex
    out             io.Writer

    logFormatFunc   func(logType LogType, i interface{}) (string, []interface{}, bool)
    logLevel        LogType
}

// 声明接口实现者
var _ ILogger = &Logger{}

func (l *Logger) Init() {
    l.mu.Lock()
    defer l.mu.Unlock()
    l.logFormatFunc = l.DefaultLogFormatFunc
    l.out = os.Stdout
    l.logLevel = DEBUG
}

func (l* Logger) SetLogLevel(logType LogType) {
    l.mu.Lock()
    defer l.mu.Unlock()
    l.logLevel = logType
}

func (l* Logger) GetLogLevel() LogType{
    l.mu.Lock()
    defer l.mu.Unlock()
    return l.logLevel
}

// 设置格式化log输出函数
// 函数返回 format 和 对应格式 []interface{}
func (l *Logger) SetLoggerFormat(formatFunc func(logType LogType, i interface{}) (string, []interface{}, bool))  {
    l.mu.Lock()
    defer l.mu.Unlock()
    l.logFormatFunc = formatFunc
}

// 输出信息
func (l *Logger) Debug(i interface{})  {
    l.log(DEBUG, i)
}

func (l *Logger) Info(i interface{}){
    l.log(INFO, i)
}

func (l *Logger) Notice(i interface{}){
    l.log(NOTICE, i)
}

func (l *Logger) Warn(i interface{})  {
    l.log(WARN, i)
}

func (l *Logger) Error(i interface{})  {
    l.log(ERROR, i)
}

func (l *Logger) Critical(i interface{}) {
    l.log(CRITICAL, i)
}

func (l *Logger) Fatal(i interface{}) {
    l.log(FATAL, i)
}

func (l *Logger) DefaultLogFormatFunc(logType LogType, i interface{})  (string, []interface{}, bool){
    format := "\033["+ logTypesColors[logType] + "m%s [%s] %s \033[0m\n"
    layout :="2006-01-02 15:04:05.999"
    formatTime := time.Now().Format(layout)
    if len(formatTime) != len(layout) {
        // 可能出现结尾是0被省略 如 2006-01-02 15:04:05.9 2006-01-02 15:04:05.99，补上
        formatTime += ".000"[4 - (len(layout) - len(formatTime)):4]
    }

    values := make([]interface{}, 3)
    values[0] = logTypeStrings[logType]
    values[1] = formatTime
    values[2] = fmt.Sprint(i)

    return format, values, true
}

func (l *Logger) log(logType LogType, i interface{})  {
    l.mu.Lock()
    defer l.mu.Unlock()

    if l.logLevel > logType{
        return
    }

    format, data, isLog := l.logFormatFunc(logType, i)
    if !isLog{ return}

    _, err := fmt.Fprintf(l.out, format, data...)
    if err !=nil{
        panic(err)
    }
}
