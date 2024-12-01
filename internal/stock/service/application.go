package service

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/xmhu2001/gorder-system/common/metrics"
	"github.com/xmhu2001/gorder-system/stock/adapters"
	"github.com/xmhu2001/gorder-system/stock/app"
	"github.com/xmhu2001/gorder-system/stock/app/query"
)

func NewApplication(_ context.Context) app.Application {
	stockRepo := adapters.NewMemoryStockRepository()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			GetItems:            query.NewGetItemsHandler(stockRepo, logger, metricClient),
			CheckIfItemsInStock: query.NewCheckIfItemsInStockHandler(stockRepo, logger, metricClient),
		},
	}
}
