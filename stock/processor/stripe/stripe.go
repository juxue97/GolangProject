package stripe

import (
	"fmt"
	"log"

	pb "github.com/juxue97/common/api"
	"github.com/juxue97/stock/processor"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/price"
	"github.com/stripe/stripe-go/v81/product"
)

// var gatewayHTTPAddr = common.GetString("HTTP_ADDR", "http://localhost:8080")

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
		Metadata:    p.Metadata,
	}
	prod, err := product.New(productParams)
	if err != nil {
		return "", "", err
	}

	// Create a price for the product (one-time payment)
	priceParams := &stripe.PriceParams{
		Product:    stripe.String(prod.ID),
		Currency:   stripe.String(p.Currency),            // Malaysian Ringgit
		UnitAmount: stripe.Int64(int64((p.Price * 100))), // RM 999 (Stripe uses cents)
	}

	price, err := price.New(priceParams)
	if err != nil {
		log.Fatalf("Error creating price: %v", err)
	}

	return prod.ID, price.ID, nil
}

func (s *Stripe) UpdateProduct(prodID, priceID string, i processor.Item) (string, error) {
	productParams := &stripe.ProductParams{
		Active: stripe.Bool(i.Active),
	}

	if i.Name != "" {
		productParams.Name = stripe.String(i.Name)
	}

	if i.Description != "" {
		productParams.Description = stripe.String(i.Description)
	}

	if i.Metadata != nil {
		productParams.Metadata = i.Metadata
	}

	newProduct, err := product.Update(prodID, productParams)
	if err != nil {
		return "", err
	}

	// if no price, no need proceed further
	if i.Price <= 0 {
		return "", nil
	}

	newPriceParams := &stripe.PriceParams{
		Product:    stripe.String(newProduct.ID),
		Currency:   stripe.String(i.Currency),
		UnitAmount: stripe.Int64(int64(i.Price * 100)),
	}

	// Deactivate the old price
	_, err = price.Update(priceID, &stripe.PriceParams{
		Active: stripe.Bool(false), // Mark old price as inactive
	})
	if err != nil {
		return "", fmt.Errorf("failed to deactivate old price: %w", err)
	}

	// Create new price
	newPrice, err := price.New(newPriceParams)
	if err != nil {
		return "", fmt.Errorf("failed to create new price: %w", err)
	}

	return newPrice.ID, nil
}
