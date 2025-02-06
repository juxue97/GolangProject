package stripe

import (
	"log"

	"github.com/juxue97/common"
	pb "github.com/juxue97/common/api"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/price"
	"github.com/stripe/stripe-go/v81/product"
)

var gatewayHTTPAddr = common.GetString("HTTP_ADDR", "http://localhost:8080")

type Stripe struct{}

func NewProcessor() *Stripe {
	return &Stripe{}
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
