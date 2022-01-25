package main

import (
	"log"

	"github.com/hsmtkk/solid-pancake/env"
	"github.com/hsmtkk/solid-pancake/msg"
	"github.com/nats-io/nats.go"
)

func main() {
	natsURL := env.MandatoryString("NATS_URL")
	natsSubject := env.MandatoryString("NATS_SUBJECT")

	natsConn, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("failed to connect NATS; %s", err)
	}
	defer natsConn.Close()

	for {
		m := msg.New()
		data, err := m.ToJSON()
		if err != nil {
			log.Print(err)
		}
		natsConn.Publish(natsSubject, data)
	}
}
