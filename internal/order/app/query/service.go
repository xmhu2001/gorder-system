package query

import (
	"context"
	"github.com/xmhu2001/gorder-system/common/genproto/orderpb"
)

// 与 stock service 进行GRPC通信

type StockService interface {
	GetItems(ctx context.Context, itemIDs []string) ([]*orderpb.Item, error)
	CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) error
}
