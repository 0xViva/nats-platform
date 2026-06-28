package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"log"
	"time"

	"nats-platform/distributor/api"
	"nats-platform/distributor/views"
)

func main() {
	r := gin.Default()

	api.Setup(r)
	views.Setup(r)

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	sub, err := js.PullSubscribe("foo", "worker")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("subscribed to foo and worker")

	// background consumer
	go func() {
		for {
			msgs, err := sub.Fetch(10)
			if err != nil {
				log.Println("fetch error:", err)
				continue
			}

			for _, msg := range msgs {
				meta, err := msg.Metadata()
				if err != nil {
					log.Println("metadata error:", err)
					continue
				}

				log.Printf(
					"Subject: %s | Message received: %s | published: %s | seq: %d",
					string(msg.Subject),
					string(msg.Data),
					meta.Timestamp.Format(time.RFC3339),
					meta.Sequence.Stream,
				)

				msg.Ack()
			}
		}
	}()

	// start HTTP server (this now runs)
	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
