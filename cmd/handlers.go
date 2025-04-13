package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

//go:embed index.html.tmpl
var tmplFS embed.FS

var (
	tmpl     = template.Must(template.ParseFS(tmplFS, "index.html.tmpl"))
	upgrader = websocket.Upgrader{}
)

func registerRoutes() {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/ws", handleWS)
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error:", err)
		return
	}
	clients.Add(ws)
	defer clients.Delete(ws)


	for {
		if _, _, err := ws.ReadMessage(); err != nil {
			break
		}
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, struct{ Title string }{title}); err != nil {
		log.Fatal("error:", err)
	}
}
