package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"

	"github.com/juxue97/common"
	pb "github.com/juxue97/common/api"
	"github.com/juxue97/common/broker"
	"google.golang.org/grpc"
)

type gRPCHandler struct {
	pb.UnimplementedOrderServiceServer
	service *loggingMiddleware
	channel *amqp.Channel
}

func NewGRPCHandler(gRPCServer *grpc.Server, service *loggingMiddleware, channel *amqp.Channel) {
	handler := &gRPCHandler{
		service: service,
		channel: channel,
	}
	pb.RegisterOrderServiceServer(gRPCServer, handler)
}

func (h *gRPCHandler) CreateOrder(ctx context.Context, payload *pb.CreateOrderRequest) (*pb.Order, error) {
	q, err := h.channel.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	tr := otel.Tracer("amqp")
	amqpContext, span := tr.Start(ctx, fmt.Sprintf(
		"AMQP - publish - %s", q.Name,
	))

	defer span.End()

	items, err := h.service.validateOrder(amqpContext, payload)
	if err != nil {
		return nil, err
	}

	o, err := h.service.createOrder(amqpContext, payload, items)
	if err != nil {
		return nil, err
	}
	if len(o.Items) == 0 {
		return nil, common.ErrNoDoc
	}

	marshalledOrder, err := json.Marshal(o)
	if err != nil {
		log.Fatal(err)
	}

	headers := broker.InjectAMQPHeaders(amqpContext)

	h.channel.PublishWithContext(amqpContext, "", q.Name, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Headers:      headers,
		Body:         marshalledOrder,
		DeliveryMode: amqp.Persistent,
	})

	return o, nil
}

func (h *gRPCHandler) GetOrder(ctx context.Context, payload *pb.GetOrderRequest) (*pb.Order, error) {
	return h.service.getOrder(ctx, payload)
}

func (h *gRPCHandler) UpdateOrder(ctx context.Context, payload *pb.Order) (*pb.Order, error) {
	return h.service.updateOrder(ctx, payload)
}

func (h *gRPCHandler) GetOrderForStockUpdate(ctx context.Context, payload *pb.GetOrderRequest) (*pb.Order, error) {
	return h.service.getOrderForStock(ctx, payload)
}
