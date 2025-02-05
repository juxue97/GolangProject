package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	pb "github.com/juxue97/common/api"
	"github.com/juxue97/common/broker"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/webhook"
	"go.opentelemetry.io/otel"
)

type PaymentHTTPHandler struct {
	channel *amqp.Channel
}

func NewPaymentHTTPHandler(channel *amqp.Channel) *PaymentHTTPHandler {
	return &PaymentHTTPHandler{channel: channel}
}

func (h *PaymentHTTPHandler) registerRouters(router *http.ServeMux) {
	router.HandleFunc("/webhook", h.handlerCheckoutWebhook)
}

func (h *PaymentHTTPHandler) handlerCheckoutWebhook(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	event, err := webhook.ConstructEvent(payload, r.Header.Get("Stripe-Signature"),
		endpointStripeSecret)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
		w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}
	switch event.Type {
	case "payment_intent.succeeded":
		// log.Println("payment succeeded")

	case "payment_intent.payment_failed":
		// log.Println("payment failed")

	case "payment_intent.created":
		// log.Println("payment created")

	case "checkout.session.completed":
		var checkoutSession stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &checkoutSession)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if checkoutSession.PaymentStatus == "paid" {
			log.Printf("Payment for checkout session %s succeeded.", checkoutSession.ID)
			// publish broker message here next
			orderID := checkoutSession.Metadata["orderID"]
			customerID := checkoutSession.Metadata["customerID"]

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			o := &pb.Order{
				Status:      "paid",
				PaymentLink: "",
				ID:          orderID,
				CustomerID:  customerID,
			}

			marshalledOrder, err := json.Marshal(o)
			if err != nil {
				log.Fatal(err)
			}

			tr := otel.Tracer("amqp")
			amqpContext, messageSpan := tr.Start(ctx, fmt.Sprintf(
				"AMQP - publish - %s", broker.OrderPaidEvent,
			))
			defer messageSpan.End()

			headers := broker.InjectAMQPHeaders(amqpContext)

			h.channel.PublishWithContext(amqpContext, broker.OrderPaidEvent, "", false, false, amqp.Publishing{
				ContentType:  "application/json",
				Headers:      headers,
				Body:         marshalledOrder,
				DeliveryMode: amqp.Persistent,
			})
			// log.Println("event published: order paid")
		}

	case "mandate.updated":
		// log.Println("mandate updated")

	case "charge.succeeded":
		// log.Println("charge succeeded")

	case "charge.updated":
		// log.Println("charge updated")

	default:
		fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)

	}
}
