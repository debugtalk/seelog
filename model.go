package wstailog

import "golang.org/x/net/websocket"

type slog struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type logLine struct {
	LogName string `json:"logName"`
	Data    string `json:"data"`
}

type wsClient struct {
	id      string
	socket  *websocket.Conn
	send    chan logLine
	logName string
}

type wsClientManager struct {
	clients    map[*wsClient]bool
	broadcast  chan logLine
	register   chan *wsClient
	unregister chan *wsClient
}
