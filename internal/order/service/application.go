package service

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/xmhu2001/gorder-system/common/metrics"
	"github.com/xmhu2001/gorder-system/order/adapters"
	"github.com/xmhu2001/gorder-system/order/app"
	"github.com/xmhu2001/gorder-system/order/app/command"
	"github.com/xmhu2001/gorder-system/order/app/query"
)

func NewApplication(ctx context.Context) app.Application {
	orderRepo := adapters.NewMemoryOrderRepository()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{
			CreateOrder: command.NewCreateOrderHandler(orderRepo,logger, metricClient),
			UpdateOrder: command.NewUpdateOrderHandler(orderRepo, logger, metricClient),
		},
		Queries: app.Queries{
			GetCustomerOrder: query.NewGetCustomerOrderHandler(orderRepo, logger, metricClient),
		},
	}
}
