package main

import (
	"context"

	pb "github.com/juxue97/common/api"
	"github.com/juxue97/payment/gateway"
	"github.com/juxue97/payment/processor"
)

type paymentService struct {
	stripeProcessor processor.PaymentProcessor
	gateway         gateway.OrdersGateway
}

func NewPaymentService(stripeProcessor processor.PaymentProcessor, gateway gateway.OrdersGateway) *paymentService {
	return &paymentService{
		stripeProcessor: stripeProcessor,
		gateway:         gateway,
	}
}

func (s *paymentService) CreatePayment(ctx context.Context, o *pb.Order) (string, error) {
	// connect to payment processor, return link
	link, err := s.stripeProcessor.CreatePaymentLink(o)
	if err != nil {
		return "", err
	}

	// update order with the link
	if err := s.gateway.UpdateOrderAfterPaymentLink(ctx, o.ID, link); err != nil {
		return "", err
	}

	return link, nil
}
