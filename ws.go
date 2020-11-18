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
		case c := <-manager.register:
			manager.clients[c] = true
		case c := <-manager.unregister:
			if _, ok := manager.clients[c]; ok {
				close(c.send)
				c.socket.Close()
				delete(manager.clients, c)
			}
		case line := <-manager.broadcast:
			for c := range manager.clients {
				if c.logName == line.LogName {
					c.send <- line
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
