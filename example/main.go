package main

import (
	"fmt"
	"github.com/debugtalk/seelog"
	"log"
	"os"
	"time"
)

const (
	DebugLog = "debug.log"
	ErrLog   = "err.log"
)

func main() {

	// 测试
	seelog.See("错误日志", ErrLog)
	seelog.See("调试日志", DebugLog)
	seelog.Serve(27129, "password")

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
