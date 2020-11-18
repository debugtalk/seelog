package wstailog

import (
	"fmt"
	"golang.org/x/net/websocket"
	"html/template"
	"log"
	"net/http"
	"time"
)

// start http server
func startServer(port int) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("http startServer panic error: %v", err)
		}
	}()

	// socket
	http.Handle("/ws", websocket.Handler(createWSConnection))

	// page
	http.HandleFunc("/wstailog", func(writer http.ResponseWriter, request *http.Request) {
		renderWebPage(writer, webPageContent, slogs)
	})

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	log.Println(err)
}

// response page
func renderWebPage(writer http.ResponseWriter, page string, slogs interface{}) {
	t, err := template.New("").Parse(page)
	if err != nil {
		log.Printf("renderWebPage error: %v", err)
		return
	}
	t.Execute(writer, slogs)
}

// create wsClient
func createWSConnection(conn *websocket.Conn) {
	client := &wsClient{time.Now().String(), conn, make(chan logLine, 1), slogs[0].Name}
	manager.register <- client
	go client.read()
	client.write()
}
