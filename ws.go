package wstailog

import (
	"encoding/json"
	"golang.org/x/net/websocket"
	"io"
	"log"
)

//  websocket客户端
type client struct {
	id     string
	socket *websocket.Conn
	send   chan msg
	see    string
}

// 客户端管理
type clientManager struct {
	clients    map[*client]bool
	broadcast  chan msg
	register   chan *client
	unregister chan *client
}

var manager = clientManager{
	broadcast:  make(chan msg),
	register:   make(chan *client),
	unregister: make(chan *client),
	clients:    make(map[*client]bool),
}

func (manager *clientManager) start() {
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
		case msg := <-manager.broadcast:
			for conn := range manager.clients {
				if conn.see == msg.LogName {
					conn.send <- msg
				}
			}
		}
	}
}

func (c *client) write() {

	for msg := range c.send {
		msgByte, _ := json.Marshal(msg) // 忽略错误
		_, err := c.socket.Write(msgByte)
		if err != nil {
			manager.unregister <- c
			log.Printf("write error: %v", err)
			break
		}
	}
}

func (c *client) read() {
	for {
		var reply string
		if err := websocket.Message.Receive(c.socket, &reply); err != nil {
			if err != io.EOF {
				log.Printf("read error: %v", err)
				manager.unregister <- c
			}
			break
		}
		type recv struct {
			LogName string `json:"logName"`
		}
		var rcv = &recv{}
		if err := json.Unmarshal([]byte(reply), &rcv); err != nil {
			manager.unregister <- c
			log.Printf("read error: %v", err)
			break
		}
		c.see = rcv.LogName
	}
}
