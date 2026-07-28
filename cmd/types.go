package main

import (
	"container/ring"
	"context"
	"sync"

	"github.com/coder/websocket"
)

type Clients struct {
	clients map[*websocket.Conn]struct{}
	mu      sync.Mutex
}

func NewClients() *Clients {
	return &Clients{
		clients: make(map[*websocket.Conn]struct{}),
	}
}

func (c *Clients) Add(ws *websocket.Conn) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.clients[ws] = struct{}{}
}

func (c *Clients) Delete(ws *websocket.Conn) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.clients, ws)
	ws.CloseNow()
}

func (c *Clients) Broadcast(msg []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for ws := range c.clients {
		ws.Write(context.Background(), websocket.MessageText, msg)
	}
}

type LogBuffer struct {
	ring *ring.Ring
}

func NewLogBuffer(size int) *LogBuffer {
	return &LogBuffer{
		ring: ring.New(size),
	}
}

func (r *LogBuffer) AddLine(line []byte) {
	r.ring.Value = append([]byte(nil), line...)
	r.ring = r.ring.Next()
}

func (r *LogBuffer) Do(f func([]byte)) {
	r.ring.Do(func(s any) {
		if s, ok := s.([]byte); ok {
			f(s)
		}
	})
}
