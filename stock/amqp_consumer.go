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
	q, err := ch.QueueDeclare("", true, false, true, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = ch.QueueBind(q.Name, "", broker.OrderPaidEvent, false, nil)
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

			log.Printf("Received a message: %s", d.Body)

			order := string(d.Body)

			// access the order service database using grpc gateway
			// using orderID, look for the product id and quantity purchased
			itemsWithQuantity, err := c.service.GetOrderService(ctx, o)
			if err != nil {
				log.Printf("failed to get order: %v", err)
				d.Nack(false, false)
				continue
			}
			for _, item := range itemsWithQuantity {
				q := int(item.Quantity)
				_, err := c.service.DeductStock(ctx, item.ID, q)
				if err != nil {
					log.Printf("failed to update stock: %v", err)
					if err := broker.HandleRetry(ch, &d); err != nil {
						log.Printf("failed to handle retry: %v", err)
					}
					continue
				}
			}

			messageSpan.AddEvent("stock.updated")
			messageSpan.End()
			log.Printf("Stock quantity updated for order: %s", order)

			d.Ack(false)
		}
	}()

	<-forever
}
