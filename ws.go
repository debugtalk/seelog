package wstailog

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net/http"
	"time"
)

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

func (manager *wsClientManager) read(c *wsClient) {
	for {
		var reply string
		if err := websocket.Message.Receive(c.socket, &reply); err != nil {
			if err != io.EOF {
				log.Printf("receive message error: %v", err)
				manager.unregister <- c
			}
			break
		}
		var line = &logLine{}
		if err := json.Unmarshal([]byte(reply), &line); err != nil {
			manager.unregister <- c
			log.Printf("parse received message error: %v", err)
			break
		}
		c.logName = line.LogName
	}
}

func (manager *wsClientManager) write(c *wsClient) {
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

// create wsClient
func (manager *wsClientManager) createWSConnection(conn *websocket.Conn) {
	client := &wsClient{time.Now().String(), conn, make(chan logLine, 1), slogs[0].LogName}
	manager.register <- client
	go manager.read(client)
	manager.write(client)
}

// start http server
func (manager *wsClientManager) listenAndServe(port int) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("http startServer panic error: %v", err)
		}
	}()

	// socket
	http.Handle("/ws", websocket.Handler(manager.createWSConnection))

	// page
	http.HandleFunc("/wstailog", func(writer http.ResponseWriter, request *http.Request) {
		renderWebPage(writer, slogs)
	})

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	log.Println(err)
}
