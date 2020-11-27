package main

import (
	"fmt"
	"github.com/debugtalk/wstailog"
	"log"
	"os"
	"time"
)

const (
	DebugLog = "debug.log"
	ErrLog   = "err.log"
)

func main() {
	wstailog.Tail("错误日志", ErrLog)
	wstailog.Tail("调试日志", DebugLog)
	wstailog.Serve(27129)

	// 模拟服务输出日志
	go printLog("调试日志", DebugLog)
	go printLog("错误日志", ErrLog)
	select {}
}

func printLog(name, path string) {
	// 模拟日志输出
	err := os.Remove(path)
	if err != nil {
		log.Println(err)
	}

	f, err := os.Create(path)
	if err != nil {
		log.Println(err)
		return
	}

	for t := range time.Tick(time.Second * 1) {
		testLog := fmt.Sprintf("「%s」[%s]\n", name, t.String())
		_, err := f.WriteString(testLog)
		if err != nil {
			log.Println(err.Error())
		}
	}
}
