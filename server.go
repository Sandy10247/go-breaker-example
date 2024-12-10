// server.go
package main

import (
	"log"
	"net/http"
	"os"
)

// ExampleServer is a test server to check the "CircuitBreaker" pattern
type ExampleServer struct {
	addr      string
	logger    *log.Logger
	isEnabled bool
}

// NewExampleServer creates the instance of our server
func NewExampleServer(addr string) *ExampleServer {
	return &ExampleServer{
		addr:      addr,
		logger:    log.New(os.Stdout, "Server\t", log.LstdFlags),
		isEnabled: true,
	}
}

// ListenAndServe starts listening on the address provided
// on creating the instance.
func (s *ExampleServer) ListenAndServe() error {
	// The main endpoint we will request to
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if s.isEnabled {
			s.logger.Println("responded with OK ✅")
			w.WriteHeader(http.StatusOK)
		} else {
			s.logger.Println("responded with Error 🔴")
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	// Toggle endpoint to switch on and off responses from the main one
	http.HandleFunc("/toggle", func(w http.ResponseWriter, r *http.Request) {
		s.isEnabled = !s.isEnabled
		var state = "🟢"
		if !s.isEnabled {
			state = "🔴"
		}
		s.logger.Printf("Toggled to %v\n", state)
		w.WriteHeader(http.StatusOK)
	})

	return http.ListenAndServe(s.addr, nil)
}
