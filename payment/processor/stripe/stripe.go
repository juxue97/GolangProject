package stripe

import (
	"fmt"
	"log"

	"github.com/juxue97/common"
	pb "github.com/juxue97/common/api"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
	"github.com/stripe/stripe-go/v81/price"
	"github.com/stripe/stripe-go/v81/product"
)

var gatewayHTTPAddr = common.GetString("HTTP_ADDR", "http://localhost:8080")

type Stripe struct{}

func NewProcessor() *Stripe {
	return &Stripe{}
}

func (s *Stripe) CreatePaymentLink(o *pb.Order) (string, error) {
	log.Printf("Creating payment link for order %v", o)

	items := []*stripe.CheckoutSessionLineItemParams{}
	for _, item := range o.Items {
		items = append(items, &stripe.CheckoutSessionLineItemParams{
			Price:    stripe.String(item.PriceID),
			Quantity: stripe.Int64(int64(item.Quantity)),
		})
	}

	gatewaySuccessURL := fmt.Sprintf("%s/success.html?customerID=%s&orderID=%s", gatewayHTTPAddr, o.CustomerID, o.ID)
	gatewayCancelURL := fmt.Sprintf("%s/cancel.html", gatewayHTTPAddr)

	params := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String(gatewaySuccessURL),
		CancelURL:  stripe.String(gatewayCancelURL),
		LineItems:  items,
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		Metadata: map[string]string{
			"orderID":    o.ID,
			"customerID": o.CustomerID,
		},
	}
	result, err := session.New(params)
	if err != nil {
		return "", err
	}

	return result.URL, nil
}

func (s *Stripe) CreateProduct(p *pb.Product) (string, string, error) {
	// Create a new product
	productParams := &stripe.ProductParams{
		Name:        stripe.String(p.Name),
		Description: stripe.String(p.Description),
		TaxCode:     stripe.String("txcd_00000000"), // no tax
		Active:      stripe.Bool(true),
	}
	prod, err := product.New(productParams)
	if err != nil {
		return "", "", err
	}

	// Create a price for the product (one-time payment)
	priceParams := &stripe.PriceParams{
		Product:    stripe.String(prod.ID),
		Currency:   stripe.String(p.Currency),          // Malaysian Ringgit
		UnitAmount: stripe.Int64(int64(p.Price * 100)), // RM 999 (Stripe uses cents)
	}

	price, err := price.New(priceParams)
	if err != nil {
		log.Fatalf("Error creating price: %v", err)
	}

	return prod.ID, price.ID, nil
}
