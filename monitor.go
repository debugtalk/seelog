package wstailog

import (
	"context"
	"errors"
	"fmt"
	"github.com/hpcloud/tail"
	"log"
	"os"
	"time"
)

func (manager *wsClientManager) monitorLogFile(sl slog) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("monitorLogFile panic error: %v", err)
		}
	}()

	fileInfo, err := os.Stat(sl.LogPath)
	if err != nil {
		log.Printf("wait log file to be created: %s", sl.LogPath)
		fileInfo, err = blockUntilFileExists(sl.LogPath)
		if err != nil {
			log.Fatalf(fmt.Sprintf("log file is not created, error: %v", err))
			return
		}
	}

	log.Printf("start to monitor log file: %s", sl.LogPath)

	t, err := tail.TailFile(
		sl.LogPath,
		tail.Config{
			Follow:   true,
			ReOpen:   true,
			Location: &tail.SeekInfo{Offset: fileInfo.Size(), Whence: 0},
			Logger:   tail.DiscardingLogger,
		},
	)

	for line := range t.Lines {
		manager.broadcast <- logLine{sl.LogName, line.Text}
	}
}

// monitor all log files
func (manager *wsClientManager) monitorAllLogs(slogs []slog) {
	for _, sl := range slogs {
		go manager.monitorLogFile(sl)
	}
}

// wait log file to be created
func blockUntilFileExists(fileName string) (os.FileInfo, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Minute*5)
	for {
		if f, err := os.Stat(fileName); err == nil {
			return f, nil
		}

		select {
		case <-time.After(time.Millisecond * 200):
			continue
		case <-ctx.Done():
			return nil, errors.New(fmt.Sprintf("TimeoutError for waiting log file: %s", fileName))
		}
	}
}
