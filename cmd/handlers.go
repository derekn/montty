package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"

	"github.com/coder/websocket"
)

//go:embed templates/index.html.tmpl
var tmplFS embed.FS

//go:embed static/*.css static/*.js static/*.svg static/*.png
var staticFS embed.FS

var tmpl = template.Must(template.ParseFS(tmplFS, "templates/index.html.tmpl"))

func registerRoutes() {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/ws", handleWS)
	http.Handle("/static/", http.FileServer(http.FS(staticFS)))
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println("upgrade connection error:", err)
		return
	}
	clients.Add(ws)
	defer clients.Delete(ws)

	// replay history for new clients
	for _, line := range logBuffer.Lines() {
		if err := ws.Write(r.Context(), websocket.MessageText, fmtOutput(line)); err != nil {
			log.Println("replay error:", err)
			break
		}
	}

	for {
		if _, _, err := ws.Read(r.Context()); err != nil {
			break
		}
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	tmplData := struct {
		AppName string
		Title   string
		CSSUrl  string
	}{
		appName, args.Title, args.CSSUrl,
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, tmplData); err != nil {
		log.Println("template error:", err)
	}
}
