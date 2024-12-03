package adapters

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/xmhu2001/gorder-system/common/genproto/orderpb"
)

// 实现OrderService接口
type OrderGRPC struct {
	client orderpb.OrderServiceClient // embed orderpb的一个client
}

func NewOrderGRPC(client orderpb.OrderServiceClient) *OrderGRPC {
	return &OrderGRPC{client: client}
}

// server 端：order/ports
func (o OrderGRPC) UpdateOrder(ctx context.Context, order *orderpb.Order) error {
	_, err := o.client.UpdateOrder(ctx, order)
	logrus.Infof("payment_adapter||update_order, err=%v", err)
	return err
}
