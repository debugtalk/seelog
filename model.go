package wstailog

import "golang.org/x/net/websocket"

type slog struct {
	LogName string `json:"logName"`
	LogPath string `json:"logPath"`
}

type logLine struct {
	LogName string `json:"logName"`
	Text    string `json:"text"`
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
