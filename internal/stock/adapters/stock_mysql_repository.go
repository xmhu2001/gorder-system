package adapters

import (
	"context"
	"github.com/xmhu2001/gorder-system/stock/entity"
	"github.com/xmhu2001/gorder-system/stock/infrastructure/persistent"
)

type MySQLStockRepository struct {
	db *persistent.MySQL
}

func NewMySQLStockRepository(db *persistent.MySQL) *MySQLStockRepository {
	return &MySQLStockRepository{db: db}
}

func (m MySQLStockRepository) GetItems(ctx context.Context, ids []string) ([]*entity.Item, error) {
	//TODO implement me
	panic("implement me")
}

func (m MySQLStockRepository) GetStock(ctx context.Context, ids []string) ([]*entity.ItemWithQuantity, error) {
	data, err := m.db.BatchGetStockByID(ctx, ids)
	if err != nil {
		return nil, err
	}
	var result []*entity.ItemWithQuantity
	for _, item := range data {
		result = append(result, &entity.ItemWithQuantity{
			ID:       item.ProductID,
			Quantity: item.Quantity,
		})
	}
	return result, nil
}
