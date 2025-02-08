package gateway

import (
	"context"

	pb "github.com/juxue97/common/api"
)

type StocksGateway interface {
	GetOrder(context.Context, *pb.Order) (*pb.Order, error)
}
