package seelog

import (
	"log"
)

var slogs []slog

// register monitor log file
func See(name, path string) {

	if name == "" || path == "" {
		log.Fatal("log name and path should not be empty")
		return
	}

	for _, sl := range slogs {
		if sl.Name == name {
			log.Fatalf("log name has been registered: %s", name)
			return
		}
	}
	slogs = append(slogs, slog{name, path})
}

// start monitor
func Serve(port int) {
	maxPort := 1<<16 - 1
	if port <= 0 || port > maxPort {
		log.Fatalf("port should be ranged in (0, %d]", maxPort)
		return
	}

	if len(slogs) < 1 {
		log.Fatalf("no log file has been registered")
		return
	}

	// start websocket manager
	go manager.start()

	// start to monitor all registered log files
	go monitorAllLogs(slogs)

	// start log server
	go server(port)
}
