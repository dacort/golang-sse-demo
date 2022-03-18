package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

var (
	addr = flag.String("addr", "localhost:3000", "bind address")
)

var notify struct {
	mu sync.Mutex
	ch chan struct{}

	value int
}

func main() {
	flag.Parse()
	log.Printf("golang SSE demo, listening on http://%s", *addr)

	// Initialize the notify channel
	notify.ch = make(chan struct{})

	// We also start a background goroutine to send the SSE events
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for range ticker.C {
			notify.mu.Lock()
			notify.value = rand.Intn(100)
			notify.mu.Unlock()

			// notify.ch <- struct{}{}
			close(notify.ch)
			notify.ch = make(chan struct{})
		}
	}()

	err := http.ListenAndServe(*addr, &handler{})
	if err != nil {
		log.Fatal(err)
	}
}
