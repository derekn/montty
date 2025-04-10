package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

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
	if _, err := fmt.Fprint(w, `
<!DOCTYPE html>
<html>
<head><title>Stdin Mirror</title></head>
<body>
	<pre id="output"></pre>
	<script>
		const out = document.getElementById("output");
		const ws = new WebSocket("ws://" + location.host + "/ws");
		ws.onmessage = e => {
			out.textContent += e.data;
			window.scrollTo(0, document.body.scrollHeight);
		};
	</script>
</body>
</html>`); err != nil {
		log.Println("error:", err)
	}
}
