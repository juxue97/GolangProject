package processor

import (
	pb "github.com/juxue97/common/api"
)

type StockProcessor interface {
	CreateProduct(*pb.Product) (string, string, error)
}
