package logger

// 单元测试和基准测试文件

import (
	"testing"

	"github.com/pxlh007/logger"
)

// go test -test.bench=".*" -run=none  -test.benchmem  -benchtime=3s
func BenchmarkLogger(b *testing.B) {
	// 复位计时器
	b.ResetTimer()

	// 初始化变量
	//var l *logger.Logger = logger.NewLogger()
	var lf *logger.RotateFileLogger = logger.NewRotateFileLogger("./")
	//	var s = []string{
	//		"200",
	//		"请求成功",
	//		"0.55ms",
	//		"GET",
	//		"/hello",
	//	}

	ss := "1"

	// 循环执行测试代码
	for i := 0; i < b.N; i++ {
		// 这里书写测试代码
		// l.Info(ss)
		lf.Info(ss)
	}

}
