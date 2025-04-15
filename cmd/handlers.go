package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

//go:embed templates/index.html.tmpl
var tmplFS embed.FS

var (
	tmpl     = template.Must(template.ParseFS(tmplFS, "templates/index.html.tmpl"))
	upgrader = websocket.Upgrader{}
)

func registerRoutes() {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/ws", handleWS)
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade connection error:", err)
		return
	}
	clients.Add(ws)
	defer clients.Delete(ws)

	// replay history for new clients
	logBuffer.Do(func(line []byte) {
		if err := ws.WriteMessage(websocket.TextMessage, append(line, '\n')); err != nil {
			log.Println("new connection error:", err)
			return
		}
	})

	for {
		if _, _, err := ws.ReadMessage(); err != nil {
			break
		}
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, struct{ AppName, Title string }{appName, title}); err != nil {
		log.Fatal("template error:", err)
	}
}
