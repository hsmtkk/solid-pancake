package main

import (
	"context"
	"log"

	"github.com/hsmtkk/solid-pancake/env"
	"github.com/hsmtkk/solid-pancake/msg"
	"github.com/hsmtkk/solid-pancake/traceprovider"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func main() {
	natsURL := env.MandatoryString("NATS_URL")
	natsSubject := env.MandatoryString("NATS_SUBJECT")

	tp, err := traceprovider.New("consumer")
	if err != nil {
		log.Fatal(err)
	}
	defer tp.Shutdown(context.Background())
	otel.SetTracerProvider(tp)

	natsConn, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("failed to connect NATS; %s", err)
	}
	defer natsConn.Close()

	hdl := newHandler(context.Background())
	sub, err := natsConn.Subscribe(natsSubject, hdl.Handle)
	if err != nil {
		log.Fatalf("failed to subscribe subject; %s; %s", natsSubject, err)
	}
	defer sub.Unsubscribe()

	select {}
}

type handler struct {
	ctx context.Context
}

func newHandler(ctx context.Context) *handler {
	return &handler{ctx}
}

func (hdl *handler) Handle(natsMsg *nats.Msg) {
	tr := otel.Tracer("handle")
	_, span := tr.Start(hdl.ctx, "handle")
	defer span.End()

	m, err := msg.FromJSON(natsMsg.Data)
	if err != nil {
		log.Print(err)
	}
	log.Print(m.ID)
	span.SetAttributes(attribute.Key("id").String(m.ID))
}
