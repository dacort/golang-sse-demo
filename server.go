package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dacort/golang-sse-demo/static"
)

type handler struct{}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/api/stream":
		h.handleStream(w, r)
	default:
		http.FileServer(http.FS(static.FS)).ServeHTTP(w, r)
	}
}

func (h *handler) handleStream(w http.ResponseWriter, r *http.Request) {
	log.Printf("stream connected %p", r)
	defer log.Printf("stream disconnected %p", r)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-store")

	// Create a lock and get a copy of the channel that sends updates
	notify.mu.Lock()
	notifyCh, value := notify.ch, notify.value
	notify.mu.Unlock()

	for {
		// Marshal data & write SSE message.
		if buf, err := json.Marshal(Event{Value: value}); err != nil {
			log.Printf("cannot marshal event: %s", err)
			return
		} else if _, err := fmt.Fprintf(w, "event: update\ndata: %s\n\n", buf); err != nil {
			log.Printf("cannot write update event: %s", err)
			return
		}
		w.(http.Flusher).Flush()

		// Wait for change to value.
		select {
		case <-r.Context().Done():
			return
		case <-notifyCh:
		}

		notify.mu.Lock()
		notifyCh, value = notify.ch, notify.value
		notify.mu.Unlock()
	}
}

type Event struct {
	Value int `json:"value"`
}
