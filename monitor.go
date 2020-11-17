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

func monitorLogFile(sl slog) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("monitorLogFile panic error: %v", err)
		}
	}()

	fileInfo, err := os.Stat(sl.Path)
	if err != nil {
		log.Printf("wait log file to be created: %s", sl.Path)
		fileInfo, err = blockUntilFileExists(sl.Path)
		if err != nil {
			log.Fatalf(fmt.Sprintf("log file is not created, error: %v", err))
			return
		}
	}

	log.Printf("start to monitor log file: %s", sl.Path)

	t, err := tail.TailFile(
		sl.Path,
		tail.Config{
			Follow:   true,
			ReOpen:   true,
			Location: &tail.SeekInfo{Offset: fileInfo.Size(), Whence: 0},
			Logger:   tail.DiscardingLogger,
		},
	)

	for line := range t.Lines {
		manager.broadcast <- logLine{sl.Name, line.Text}
	}
}

// monitor all log files
func monitorAllLogs(slogs []slog) {
	for _, sl := range slogs {
		go monitorLogFile(sl)
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
