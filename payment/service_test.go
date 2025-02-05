package main

import (
	"context"
	"testing"

	"github.com/juxue97/common/api"
	inmemRegistry "github.com/juxue97/common/discovery/inmem"
	"github.com/juxue97/payment/gateway"
	"github.com/juxue97/payment/processor/inmem"
)

func TestStripeService(t *testing.T) {
	processor := inmem.NewInmem()
	registry := inmemRegistry.NewRegistry()
	gateway := gateway.NewGateway(registry)
	service := NewPaymentService(processor, gateway)
	t.Run("should create payment link", func(t *testing.T) {
		link, err := service.CreatePayment(context.Background(), &api.Order{})
		if err != nil {
			t.Errorf("CreatePaymentLink failed: %v", err)
		}

		if link == "" {
			t.Error("CreatePaymentLink failed: link is empty")
		}

		if link != "dummy-link" {
			t.Errorf("CreatePaymentLink failed: Supposed to be %v, but got %v", "dummy-link", link)
		}
	})
}
