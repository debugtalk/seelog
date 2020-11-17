package wstailog

import (
	"encoding/json"
	"golang.org/x/net/websocket"
	"io"
	"log"
)

type wsClientManager struct {
	clients    map[*wsClient]bool
	broadcast  chan logLine
	register   chan *wsClient
	unregister chan *wsClient
}

func (manager *wsClientManager) start() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("start manager panic error: %v", err)
		}
	}()

	for {
		select {
		case conn := <-manager.register:
			manager.clients[conn] = true
		case conn := <-manager.unregister:
			if _, ok := manager.clients[conn]; ok {
				close(conn.send)
				conn.socket.Close()
				delete(manager.clients, conn)
			}
		case line := <-manager.broadcast:
			for conn := range manager.clients {
				if conn.logName == line.LogName {
					conn.send <- line
				}
			}
		}
	}
}

type wsClient struct {
	id      string
	socket  *websocket.Conn
	send    chan logLine
	logName string
}

func (c *wsClient) write() {
	for msg := range c.send {
		msgByte, err := json.Marshal(msg)
		if err != nil {
			continue
		}
		_, err = c.socket.Write(msgByte)
		if err != nil {
			manager.unregister <- c
			log.Printf("write error: %v", err)
			break
		}
	}
}

func (c *wsClient) read() {
	for {
		var reply string
		if err := websocket.Message.Receive(c.socket, &reply); err != nil {
			if err != io.EOF {
				log.Printf("read error: %v", err)
				manager.unregister <- c
			}
			break
		}
		var line = &logLine{}
		if err := json.Unmarshal([]byte(reply), &line); err != nil {
			manager.unregister <- c
			log.Printf("read error: %v", err)
			break
		}
		c.logName = line.LogName
	}
}
