package inmem

import pb "github.com/juxue97/common/api"

type inmem struct{}

func NewInmem() *inmem {
	return &inmem{}
}

func (i *inmem) CreatePaymentLink(o *pb.Order) (string, error) {
	return "dummy-link", nil
}

func (s *inmem) CreateProduct(p *pb.Product) (string, string, error) {
	return "dummy-product-id", "dummy-price-id", nil
}
