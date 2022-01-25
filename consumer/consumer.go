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

	sub, err := natsConn.Subscribe(natsSubject, handler)
	if err != nil {
		log.Fatalf("failed to subscribe subject; %s; %s", natsSubject, err)
	}
	defer sub.Unsubscribe()

	select {}
}

func handler(natsMsg *nats.Msg) {
	m, err := msg.FromJSON(natsMsg.Data)
	if err != nil {
		log.Print(err)
	}
	log.Print(m.ID)
}
