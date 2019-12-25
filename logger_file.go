package logs

import (
    "os"
    "time"
    "fmt"
    "runtime/debug"
)

/*
 *
 */
func NewRotateFileLogger(dir string)  *RotateFileLogger{
    logger := &RotateFileLogger{}
    logger.Init(dir)
    return logger
}



// 一段时间自动创建新的log文件
// 如果在间隔时间无日志，此间隔不会创建log文件
type RotateFileLogger struct {
    Logger

    file                    *os.File                // 正在操作文件
    dirPath                 string                  // logs 文件所在文件夹
    fileNameFormatFunc      func(t time.Time) string    // 获取文件名称格式

    newFileGapTime          time.Duration       // 创建新log的间隔时间
    lastFileTime            time.Time           // 上次创建文件，文件对应时间(依据间隔时间，不是真实创建时间)
}

func (l * RotateFileLogger) Init(dir string) {
    l.Logger.Init()
    l.mu.Lock()
    defer l.mu.Unlock()

    l.fileNameFormatFunc = l.DefaultFileNameFormat
    l.logFormatFunc = l.DefaultLogFormatFunc
    l.newFileGapTime = 0
    l.lastFileTime = time.Now()
    l.dirPath = dir

    file, err := l.createLogFile(l.fileNameFormatFunc(l.lastFileTime))
    if err != nil{
        panic(err)
        return
    }

    l.file = file
    l.out = file
}

func (l * RotateFileLogger) SetNewFileGapTime(gapTime time.Duration)  {
    l.mu.Lock()
    defer l.mu.Unlock()
    l.newFileGapTime = gapTime
}
// 默认 创建格式化的文件名. fileTime 不一定是创建文件的时间
func (l * RotateFileLogger) DefaultFileNameFormat(fileTime time.Time) string{
    // tips: linux 不支持 2006-01-02 15:04:05.999 ":"名称
    // tips: os x  2006/01-02 15-04-05-999 带有"/"报no such file or directory
    layout :="2006-01-02 15-04-05-999"
    formatTime := fileTime.Format(layout)
    if len(formatTime) != len(layout) {
        // 对于如果是 出现2006-01-02 15:04:05.99  适配处理 成2006-01-02 15:04:05.990
        formatTime += ".000"[4 - (len(layout) - len(formatTime)):4]
    }
    return formatTime + ".log"
}

func (l * RotateFileLogger) createLogFile(filename string)  (*os.File, error){

    if len(l.dirPath) != 0{
        filename = l.dirPath + "/" + filename
        err:=os.MkdirAll(l.dirPath, 0777)
        if err != nil{
            panic(err)
        }
    }

    file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
    if err != nil{
        return nil, err
    }
    return file, nil
}

func (l *RotateFileLogger) DefaultLogFormatFunc(logType LogType, i interface{})  (string, []interface{}, bool){
    format := "%s [%s] %s\n"
    now := time.Now()
    gapTime := now.Sub(l.lastFileTime)
    if gapTime > l.newFileGapTime && l.newFileGapTime > 0{
        l.file.Close()

        rate := int(int64(gapTime) / int64(l.newFileGapTime))
        l.lastFileTime = l.lastFileTime.Add(l.newFileGapTime * time.Duration(rate))
        file, err := l.createLogFile(l.fileNameFormatFunc(l.lastFileTime))
        if err != nil{
            file.Close()
            panic(err)
        }

        l.file = file
        l.out = file
    }

    layout :="2006-01-02 15:04:05.999"
    formatTime := now.Format(layout)
    defer func() {
        e := recover()
        if e!=nil{
            panic(debug.Stack())
        }
    }()
    if len(formatTime) != len(layout) {
        // 对于如果是 出现2006-01-02 15:04:05.99  适配处理 成2006-01-02 15:04:05.990
        formatTime += ".000"[4 - (len(layout) - len(formatTime)):4]
    }

    values := make([]interface{}, 3)
    values[0] = GetLogTypeString(logType)
    values[1] = formatTime
    values[2] = fmt.Sprint(i)

    return format, values, true
}
