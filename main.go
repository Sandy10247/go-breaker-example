// main.go
package main

import (
	"log"
	"os"
	"time"
)

func main() {
	logger := log.New(os.Stdout, "Main\t", log.LstdFlags)
	server := NewExampleServer(":8080")

	go func() {
		_ = server.ListenAndServe()
	}()

	var client NotificationClient

	client = NewSmsClient("http://127.0.0.1:8080")

	client = NewClientCircuitBreakerProxy(client)

	for {
		err := client.Send()
		time.Sleep(1 * time.Second)
		if err != nil {
			logger.Println("caught an error", err)
		}
	}
}
