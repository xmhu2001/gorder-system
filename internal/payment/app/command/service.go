package command

import (
	"context"

	"github.com/xmhu2001/gorder-system/common/genproto/orderpb"
)

type OrderService interface {
	UpdateOrder(ctx context.Context, order *orderpb.Order) error
}
