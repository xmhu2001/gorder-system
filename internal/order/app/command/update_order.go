package command

import (
	"github.com/sirupsen/logrus"
	"github.com/xmhu2001/gorder-system/common/decorator"
	domain "github.com/xmhu2001/gorder-system/order/domain/order"
	"context"
)

type UpdateOrder struct {
	Order *domain.Order
	UpdateFn func(context.Context, *domain.Order) (*domain.Order, error)
}

type UpdateOrderHandler decorator.CommandHandler[UpdateOrder, interface{}]

type updateOrderHandler struct {
	orderRepo domain.Repository
}

func NewUpdateOrderHandler(
	orderRepo domain.Repository,
	logger *logrus.Entry,
	metricClient decorator.MetricsClient,
) UpdateOrderHandler {
	if orderRepo == nil {
		panic("orderRepo is nil")
	}
	return decorator.ApplyCommandDecorators[UpdateOrder, interface{}](
		updateOrderHandler{orderRepo: orderRepo},
		logger,
		metricClient,
	)
}


func (c updateOrderHandler) Handle(ctx context.Context, cmd UpdateOrder) (interface{}, error) {
	if cmd.UpdateFn == nil {
		logrus.Warnf("updateOrderHandler got nil UpdateFn, order=%#v",cmd.Order)
		cmd.UpdateFn = func(_ context.Context, order *domain.Order) (*domain.Order, error) {
			return order, nil
		}
	}
	err := c.orderRepo.Update(ctx, cmd.Order, cmd.UpdateFn)
	if err != nil {
		return nil, err
	}
	return nil, nil
}