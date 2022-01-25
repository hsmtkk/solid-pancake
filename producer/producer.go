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
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

func main() {
	natsURL := env.MandatoryString("NATS_URL")
	natsSubject := env.MandatoryString("NATS_SUBJECT")

	tp, err := traceprovider.New("producer")
	if err != nil {
		log.Fatal(err)
	}
	defer tp.Shutdown(context.Background())
	otel.SetTracerProvider(tp)
	tracer = tp.Tracer("producer")

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
	_, span := tracer.Start(ctx, "publish")
	defer span.End()

	m := msg.New()
	log.Print(m.ID)
	span.SetAttributes(attribute.Key("id").String(m.ID))

	publish2(ctx)

	data, err := m.ToJSON()
	if err != nil {
		log.Print(err)
	}
	conn.Publish(subj, data)
}

func publish2(ctx context.Context) {
	_, span := tracer.Start(ctx, "publish-child")
	defer span.End()

	// some slow work
	duration := time.Duration(rand.Intn(1000))
	time.Sleep(duration * time.Millisecond)
}
