package main

import (
	"context"
	"log"
	"time"

	"pubsub/subpub/internal/subpub"
)

func main() {
	ps := subpub.NewPubSub()

	if _, err := ps.Subscribe("news", func(msg interface{}) {
		log.Printf("[FastSubscriber] Got message: %v", msg)
	}); err != nil {
		log.Fatal("Subscribe failed:", err)
	}

	if _, err := ps.Subscribe("news", func(msg interface{}) {
		time.Sleep(3 * time.Second)
		log.Printf("[SlowSubscriber] Got message after delay: %v", msg)
	}); err != nil {
		log.Fatal("Subscribe failed:", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		i := 0
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				msg := "some random msg for test abobus " + time.Now().Format("15:04:05")
				log.Println("[Publisher] Sending:", msg)
				if err := ps.Publish("news", msg); err != nil {
					log.Println("Publish error:", err)
				}
				i++
				if i == 5 {
					log.Println("[Publisher] Done sending messages")
					cancel()
					return
				}
			}
		}
	}()

	<-ctx.Done()

	if err := ps.Close(ctx); err != nil {
		log.Println("Error during shutdown:", err)
	} else {
		log.Println("PubSub shut down cleanly.")
	}
}
