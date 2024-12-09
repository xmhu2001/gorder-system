package query

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/xmhu2001/gorder-system/common/decorator"
	domain "github.com/xmhu2001/gorder-system/order/domain/order"
)

type GetCustomerOrder struct {
	CustomerID string
	OrderID    string
}

type GetCustomerOrderHandler decorator.QueryHandler[GetCustomerOrder, *domain.Order]

type getCustomerOrderHandler struct {
	orderRepo domain.Repository
}

func (g getCustomerOrderHandler) Handle(ctx context.Context, query GetCustomerOrder) (*domain.Order, error) {
	o, err := g.orderRepo.Get(ctx, query.OrderID, query.CustomerID)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func NewGetCustomerOrderHandler(
	orderRepo domain.Repository,
	logger *logrus.Entry,
	metricClient decorator.MetricsClient,
) GetCustomerOrderHandler {
	if orderRepo == nil {
		panic("orderRepo is nil")
	}
	return decorator.ApplyQueryDecorators[GetCustomerOrder, *domain.Order](
		getCustomerOrderHandler{orderRepo: orderRepo},
		logger,
		metricClient,
	)
}
