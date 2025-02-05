package inmem

import pb "github.com/juxue97/common/api"

type inmem struct{}

func NewInmem() *inmem {
	return &inmem{}
}

func (i *inmem) CreatePaymentLink(o *pb.Order) (string, error) {
	return "dummy-link", nil
}
