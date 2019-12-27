logger日志组件

## 使用案例
```
package main

import (
	// "fmt"
	"github.com/pxlh007/logger"
)

func main() {
	var l *logger.Logger = logger.NewLogger()
	l.Info("200 | ok! | 1.052µs | POST /PING ")
	l.Info("200 | ok! | 1.053µs | GET  /PONG ")

	var s = []string{
		"200-g",
		"请求成功",
		"0.55ms",
		"GET-y",
		"/hello",
	}
	l.Info(s)

	var serror = []string{
		"404-r",
		"NOT FOUND!",
		"2.556ms",
		"POST-b",
		"/",
	}
	l.Error(serror)

	l.Debug("输出debug信息！")

	// l.Error("This is an error")
	// l.Warn("This is warning!")
	// l.Debug("This is debugging!")
	// l.Fatal("This is Fatal!")
	// fmt.Println(l)

	// 文件测试
	var lf *logger.RotateFileLogger = logger.NewRotateFileLogger("./")
	lf.Info("记录进文件测试...")
	lf.Error("记录错误信息...")

	var sf = []string{
		"200",
		"请求成功",
		"0.55ms",
		"GET",
		"/hello",
	}
	lf.Info(sf)

	var sferror = []string{
		"404",
		"NOT FOUND!",
		"2.556ms",
		"POST",
		"/",
	}
	lf.Error(sferror)

}
```
