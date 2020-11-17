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
func server(port int) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("http server panic error: %v", err)
		}
	}()

	// socket
	http.Handle("/ws", websocket.Handler(genConn))

	// page
	http.HandleFunc("/wstailog", func(writer http.ResponseWriter, request *http.Request) {
		renderWebPage(writer, webPageContent, slogs)
	})

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	log.Println(err)
}

// response page
func renderWebPage(writer http.ResponseWriter, page string, data interface{}) {
	t, err := template.New("").Parse(page)
	if err != nil {
		log.Printf("renderWebPage error: %v", err)
		return
	}
	t.Execute(writer, data)
}

// create client
func genConn(ws *websocket.Conn) {
	client := &client{time.Now().String(), ws, make(chan msg, 1), slogs[0].Name}
	manager.register <- client
	go client.read()
	client.write()
}
