package adapters

import (
	"context"
	"github.com/xmhu2001/gorder-system/common/genproto/orderpb"
)

type OrderGRPC struct {
	client orderpb.OrderServiceClient
}

func NewOrderGRPC(client orderpb.OrderServiceClient) *OrderGRPC {
	return &OrderGRPC{client: client}
}

func (g *OrderGRPC) UpdateOrder(ctx context.Context, req *orderpb.Order) error {
	_, err := g.client.UpdateOrder(ctx, req)
	return err
}
