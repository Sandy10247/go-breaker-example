// circuit_breaker.go
package main

import (
	"log"
	"os"
	"time"

	"github.com/sony/gobreaker/v2"
)

type ClientCircuitBreakerProxy struct {
	client NotificationClient
	logger *log.Logger
	gb     *gobreaker.CircuitBreaker[any] // downloaded lib structure
}

// shouldBeSwitchedToOpen checks if the circuit breaker should
// switch to the Open state
func shouldBeSwitchedToOpen(counts gobreaker.Counts) bool {
	failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
	return counts.Requests >= 3 && failureRatio >= 0.6
}

func NewClientCircuitBreakerProxy(client NotificationClient) *ClientCircuitBreakerProxy {
	logger := log.New(os.Stdout, "CB\t", log.LstdFlags)

	// We need circuit breaker configuration
	cfg := gobreaker.Settings{
		// When to flush counters int the Closed state
		Interval: 5 * time.Second,
		// Time to switch from Open to Half-open
		Timeout: 7 * time.Second,
		// Function with check when to switch from Closed to Open
		ReadyToTrip: shouldBeSwitchedToOpen,
		// set Max Request in Half Open state to 5
		MaxRequests: 5,
		// On State change Handler Fn
		OnStateChange: func(_ string, from gobreaker.State, to gobreaker.State) {
			// Handler for every state change. We'll use for debugging purpose
			logger.Println("state changed from", from.String(), "to", to.String())
		},
	}

	return &ClientCircuitBreakerProxy{
		client: client,
		logger: logger,
		gb:     gobreaker.NewCircuitBreaker[any](cfg),
	}
}

func (c *ClientCircuitBreakerProxy) Send() error {
	// We call the Execute method and wrap our client's call
	_, err := c.gb.Execute(func() (any, error) {
		err := c.client.Send()
		return nil, err
	})
	return err
}
