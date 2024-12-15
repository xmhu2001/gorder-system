package query

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/xmhu2001/gorder-system/common/decorator"
	domain "github.com/xmhu2001/gorder-system/stock/domain/stock"
	"github.com/xmhu2001/gorder-system/stock/entity"
	"github.com/xmhu2001/gorder-system/stock/infrastructure/integration"
)

type CheckIfItemsInStock struct {
	Items []*entity.ItemWithQuantity
}

type CheckIfItemsInStockHandler decorator.QueryHandler[CheckIfItemsInStock, []*entity.Item]

type checkIfItemsInStockHandler struct {
	stockRepo domain.Repository
	stripeAPI *integration.StripeAPI
}

func NewCheckIfItemsInStockHandler(
	stockRepo domain.Repository,
	stripeAPI *integration.StripeAPI,
	logger *logrus.Entry,
	metricClient decorator.MetricsClient,
) CheckIfItemsInStockHandler {
	if stockRepo == nil {
		panic("stockRepo is nil")
	}
	if stripeAPI == nil {
		panic("stripeAPI is nil")
	}
	return decorator.ApplyQueryDecorators[CheckIfItemsInStock, []*entity.Item](
		checkIfItemsInStockHandler{
			stockRepo: stockRepo,
			stripeAPI: stripeAPI,
		},
		logger,
		metricClient,
	)
}

func (g checkIfItemsInStockHandler) Handle(ctx context.Context, query CheckIfItemsInStock) ([]*entity.Item, error) {
	var res []*entity.Item
	for _, i := range query.Items {
		priceID, err := g.stripeAPI.GetPriceByProductID(ctx, i.ID)
		if err != nil {
			logrus.Warnf("GetPriceByProductID error, item ID=%s, err=%v", i.ID, err)
			continue
		}
		res = append(res, &entity.Item{
			ID:       i.ID,
			Quantity: i.Quantity,
			PriceID:  priceID,
		})
	}
	return res, nil
}
