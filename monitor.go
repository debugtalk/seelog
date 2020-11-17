package seelog

import (
	"context"
	"errors"
	"fmt"
	"github.com/hpcloud/tail"
	"os"
	"time"
)

type msg struct {
	LogName string `json:"logName"`
	Data    string `json:"data"`
}

// 监控日志文件
func monitor() {

	for _, sl := range slogs {
		go func(sl slog) {
			defer func() {
				if err := recover(); err != nil {
					printError(errors.New("monitor() panic"))
				}
			}()

			// 等待文件
			fileInfo, err := os.Stat(sl.Path)
			if err != nil {
				printInfo(fmt.Sprintf("等待文件 %s 生成", sl.Path))
				ctx, _ := context.WithTimeout(context.Background(), time.Minute*5)
				fileInfo, err = BlockUntilExists(sl.Path, ctx)
				if err != nil {
					printError(err)
					return
				}
			}

			printInfo(fmt.Sprintf("开始监控文件 %s", sl.Path))

			t, err := tail.TailFile(sl.Path,
				tail.Config{
					Follow:   true,
					ReOpen:   true,
					Location: &tail.SeekInfo{Offset: fileInfo.Size(), Whence: 0},
					Logger:   tail.DiscardingLogger,
				})

			for line := range t.Lines {
				manager.broadcast <- msg{sl.Name, line.Text}
			}
		}(sl)
	}

}

func BlockUntilExists(fileName string, ctx context.Context) (os.FileInfo, error) {

	for {
		f, err := os.Stat(fileName)
		if err == nil {
			return f, nil
		}

		select {
		case <-time.After(time.Millisecond * 200):
			continue
		case <-ctx.Done():
			return nil, errors.New(fmt.Sprintf("等待 %s 超时", fileName))
		}
	}
}
