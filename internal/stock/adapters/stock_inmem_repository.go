package adapters

import (
	"context"
	"github.com/xmhu2001/gorder-system/common/genproto/orderpb"
	domain "github.com/xmhu2001/gorder-system/stock/domain/stock"
	"sync"
)

type MemoryStockRepository struct {
	lock  *sync.RWMutex
	store map[string]*orderpb.Item
}

func NewMemoryStockRepository() *MemoryStockRepository {
	return &MemoryStockRepository{
		lock:  &sync.RWMutex{},
		store: stub,
	}
}

var stub = map[string]*orderpb.Item{
	"item_1": {
		ID:       "foo_item",
		Name:     "stub item",
		Quantity: 1000,
		PriceID:  "stub_item_price_id",
	},
	"item_2": {
		ID:       "item_2",
		Name:     "stub item 2",
		Quantity: 1000,
		PriceID:  "stub_item_2_price_id",
	},
	"item_3": {
		ID:       "item_3",
		Name:     "stub item 3",
		Quantity: 1000,
		PriceID:  "stub_item_price_id",
	},
}

func (m MemoryStockRepository) GetItems(_ context.Context, ids []string) ([]*orderpb.Item, error) {
	// get 只需要加 读 锁
	m.lock.RLock()
	defer m.lock.RUnlock()
	var (
		res     []*orderpb.Item
		missing []string
	)
	for _, id := range ids {
		if item, exist := m.store[id]; exist {
			res = append(res, item)
		} else {
			missing = append(missing, id)
		}
	}
	if len(res) == len(ids) {
		return res, nil
	}
	return res, domain.NotFoundError{Missing: missing}

}
