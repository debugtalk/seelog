package seelog

import (
	"errors"
	"fmt"
	"golang.org/x/net/websocket"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

// start http server
func server(port int, password string) {

	defer func() {
		if err := recover(); err != nil {
			printError(errors.New("server panic"))
		}
	}()

	// socket
	http.Handle("/ws", websocket.Handler(genConn))

	// page
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if !(strings.Replace(request.RequestURI, "/", "", -1) == password) {
			showPage(writer, Page403, nil)
			return
		}
		showPage(writer, PageIndex, slogs)
	})
	log.Println(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

// response page
func showPage(writer http.ResponseWriter, page string, data interface{}) {
	t, err := template.New("").Parse(page)
	if err != nil {
		printError(err)
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
