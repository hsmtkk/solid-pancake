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

	tp, err := traceprovider.New("producer")
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

	ctx := context.Background()
	for {
		publish(ctx, natsConn, natsSubject)
	}
}

func publish(ctx context.Context, conn *nats.Conn, subj string) {
	tr := otel.Tracer("publish")
	_, span := tr.Start(ctx, "publish")
	defer span.End()

	m := msg.New()
	span.SetAttributes(attribute.Key("id").String(m.ID))
	data, err := m.ToJSON()
	if err != nil {
		log.Print(err)
	}
	conn.Publish(subj, data)
}
