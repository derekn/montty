package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

//go:embed index.html.tmpl
var tmplFS embed.FS
var tmpl = template.Must(template.ParseFS(tmplFS, "index.html.tmpl"))

func registerRoutes() {
	http.HandleFunc("/", handleIndex)
	http.Handle("/ws", websocket.Handler(handleWS))
}

func handleWS(ws *websocket.Conn) {
	mu.Lock()
	clients[ws] = struct{}{}
	mu.Unlock()

	defer func() {
		mu.Lock()
		delete(clients, ws)
		mu.Unlock()
		_ = ws.Close()
	}()

	buf := make([]byte, 1)
	for {
		if _, err := ws.Read(buf); err != nil {
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
