package seelog

import (
	"errors"
	"fmt"
)

type slog struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

var slogs = []slog{}

// 启动seelog
func See(name, path string) {

	if name == "" || path == "" {
		printError(errors.New("log名称或者路径不可为空"))
		return
	}

	for _, sl := range slogs {
		if sl.Name == name {
			printError(errors.New(fmt.Sprintf("log名称 %s 已存在,不可重复", name)))
			return
		}
	}
	slogs = append(slogs, slog{name, path})
}

// 开始监控
func Serve(port int, password string) {

	if port < 0 || port > 65535 {
		printError(errors.New("端口号不符合规范，port(0,65535)"))
		return
	}

	if len(slogs) < 1 {
		printError(errors.New("至少监听一个日志文件,请使用 seelog.See(name,path string)"))
		return
	}
	// 开启socket管理器
	go manager.start()

	// 监控文件
	go monitor()

	// 开启httpServer
	go server(port, password)
}
