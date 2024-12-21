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

	if err := g.checkStock(ctx, query.Items); err != nil {
		return nil, err
	}
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

func (g checkIfItemsInStockHandler) checkStock(ctx context.Context, query []*entity.ItemWithQuantity) error {
	var ids []string
	for _, i := range query {
		ids = append(ids, i.ID)
	}
	records, err := g.stockRepo.GetStock(ctx, ids)
	if err != nil {
		return err
	}
	idQuantityMap := make(map[string]int32)
	for _, i := range records {
		idQuantityMap[i.ID] += i.Quantity
	}
	// check whether stock enough
	var (
		ok       = true
		failedOn []struct {
			ID   string
			Want int32
			Have int32
		}
	)
	for _, item := range query {
		if item.Quantity > idQuantityMap[item.ID] {
			ok = false
			failedOn = append(failedOn, struct {
				ID   string
				Want int32
				Have int32
			}{
				ID:   item.ID,
				Want: item.Quantity,
				Have: idQuantityMap[item.ID],
			})
		}
	}
	if ok {
		return nil
	}
	return domain.ExceedStockError{FailedOn: failedOn}
}
