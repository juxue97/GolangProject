package gateway

import (
	"context"

	pb "github.com/juxue97/common/api"
)

type StocksGateway interface {
	CheckIfItemsInStock(context.Context, string, []*pb.ItemsWithQuantity) (bool, []*pb.Item, error)
}
