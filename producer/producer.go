package main

import (
	"context"
	"log"
	"math/rand"
	"time"

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

	publish2(ctx)

	data, err := m.ToJSON()
	if err != nil {
		log.Print(err)
	}
	conn.Publish(subj, data)
}

func publish2(ctx context.Context) {
	tr := otel.Tracer("publish2")
	_, span := tr.Start(ctx, "publish2")
	defer span.End()

	// some slow work
	duration := time.Duration(rand.Intn(1000))
	time.Sleep(duration * time.Millisecond)
}
