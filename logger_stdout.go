package logger

import (
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

type LogType int

const (
	DEBUG    = LogType(0)
	INFO     = LogType(1)
	NOTICE   = LogType(2)
	WARN     = LogType(3)
	ERROR    = LogType(4)
	CRITICAL = LogType(5)
	FATAL    = LogType(6)
)

// 定义数据段颜色
var dataColor map[string]string = map[string]string{
	"r": "41;97", // 深红
	"g": "42;97", // 绿色
	"y": "43;97", // 土黄
	"b": "44;97", // 深蓝
}

var logTypeStrings = func() []string {
	// log 类型对应的 名称字符串，用于输出，所以统一了长度，故 DEBUG 为 "DEBUG..." 和 "CRITICAL"等长
	types := []string{"DEBUG", "INFO", "NOTICE", "WARN", "ERROR", "CRITICAL", "FATAL"}
	maxTypeLen := 0
	for _, t := range types {
		if len(t) > maxTypeLen {
			maxTypeLen = len(t)
		}
	}
	for index, t := range types {
		typeLen := len(t)
		if typeLen < maxTypeLen {
			types[index] += strings.Repeat(" ", maxTypeLen-typeLen)
		}
	}
	return types
}()

// 41;97 底色深红, 加亮白色;
// 42;97 底色绿字, 加亮白色;
// 43;97 底色黄色, 加亮白色;
// 45;97 底色紫色, 加亮白色;
var logTypesColors = []string{"45;97", "42;97", "43;97", "43;97", "41;97", "41;97", "41;97"}

/*
 * 创建 Logger
 */
func NewLogger() *Logger {
	logger := &Logger{}
	logger.Init()
	return logger
}

func GetLogTypeString(t LogType) string {
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
	mu            sync.Mutex
	out           io.Writer
	logFormatFunc FormatFunc
	logLevel      LogType
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

func (l *Logger) SetLogLevel(logType LogType) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logLevel = logType
}

func (l *Logger) GetLogLevel() LogType {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.logLevel
}

// 设置格式化log输出函数
// 函数返回 format 和 对应格式 []interface{}
func (l *Logger) SetLoggerFormat(formatFunc FormatFunc) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logFormatFunc = formatFunc
}

// 输出信息
func (l *Logger) Debug(i interface{}) {
	l.log(DEBUG, i)
}

func (l *Logger) Info(i interface{}) {
	l.log(INFO, i)
}

func (l *Logger) Notice(i interface{}) {
	l.log(NOTICE, i)
}

func (l *Logger) Warn(i interface{}) {
	l.log(WARN, i)
}

func (l *Logger) Error(i interface{}) {
	l.log(ERROR, i)
}

func (l *Logger) Critical(i interface{}) {
	l.log(CRITICAL, i)
}

func (l *Logger) Fatal(i interface{}) {
	l.log(FATAL, i)
}

func (l *Logger) DefaultLogFormatFunc(logType LogType, i interface{}) (string, []interface{}, bool) {
	// 异常捕获
	defer func() {
		e := recover()
		if e != nil {
			panic(debug.Stack())
		}
	}()

	// 计算日期format
	layout := "2006/01/02 - 15:04:05.9999"
	formatTime := time.Now().Format(layout)
	if len(formatTime) != len(layout) {
		// 可能出现结尾是0被省略如：2006/01/02 - 15:04:05.9 补足成 2006/01/02 - 15:04:05.9000
		formatTime += ".000"[4-(len(layout)-len(formatTime)) : 4]
	}

	// 计算数据format
	format := ""
	values := []interface{}{}
	if iSli, ok := i.([]string); ok {
		// 切片
		l := len(iSli)
		format = "[\033[" + logTypesColors[logType] + "m%s\033[0m] %s | "
		values = make([]interface{}, l+2)
		values[0] = logTypeStrings[logType]
		values[1] = formatTime
		for j := 0; j < l; j++ {
			ls := len(iSli[j])
			tj := ""
			if ls >= 2 {
				// 截取最后两个字符
				// 颜色标志：-g 绿色; -r 红色; -b 蓝色
				color := iSli[j][ls-2:]
				if color[0] == '-' && (color[1] == 'g' || color[1] == 'b' || color[1] == 'r' || color[1] == 'y') {
					tj = iSli[j][0 : ls-2] // 去除颜色标志
					format += "\033[" + dataColor[string(color[1])] + "m%s\033[0m | "
				} else {
					tj = iSli[j]
					format += "%s | "
				}
			} else {
				tj = iSli[j]
				format += "%s | "
			}
			// 计算输出值
			values[j+2] = tj
		}
		format += "\n"
	} else if iStr, ok := i.(string); ok {
		// 文本
		format = "[\033[" + logTypesColors[logType] + "m%s\033[0m] %s | %s | \n"
		// 计算输出值
		values = make([]interface{}, 3)
		values[0] = logTypeStrings[logType]
		values[1] = formatTime
		values[2] = iStr
	}

	// 返回格式/值
	return format, values, true
}

func (l *Logger) log(logType LogType, i interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.logLevel > logType {
		return
	}

	format, data, isLog := l.logFormatFunc(logType, i)
	if !isLog {
		return
	}

	_, err := fmt.Fprintf(l.out, format, data...)
	if err != nil {
		panic(err)
	}
}
