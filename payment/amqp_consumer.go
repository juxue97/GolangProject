package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	pb "github.com/juxue97/common/api"
	"github.com/juxue97/common/broker"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
)

type consumer struct {
	service *loggingMiddleware
}

func NewConsumer(service *loggingMiddleware) *consumer {
	return &consumer{service: service}
}

func (c *consumer) Listen(ch *amqp.Channel) {
	q, err := ch.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	messages, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	var forever chan struct{}

	go func() {
		for d := range messages {
			o := &pb.Order{}
			if err := json.Unmarshal(d.Body, o); err != nil {
				d.Nack(false, false)
				log.Printf("failed to unmarshal order: %v", err)
				continue
			}

			ctx := broker.ExtractAMQPHeaders(context.Background(), d.Headers)
			tr := otel.Tracer("amqp")
			_, messageSpan := tr.Start(ctx, fmt.Sprintf(
				"AMQP - consume - %s", q.Name,
			))
			paymentLink, err := c.service.CreatePayment(ctx, o)
			if err != nil {
				log.Printf("failed to create payment: %v", err)
				if err := broker.HandleRetry(ch, &d); err != nil {
					log.Printf("failed to handle retry: %v", err)
				}

				continue
			}
			messageSpan.AddEvent(fmt.Sprintf("payment.created: %s", paymentLink))
			messageSpan.End()

			d.Ack(false)
		}
	}()

	<-forever
}
