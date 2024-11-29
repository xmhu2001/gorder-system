package adapters

import (
	"github.com/xmhu2001/gorder-system/common/genproto/orderpb"
	domain "github.com/xmhu2001/gorder-system/stock/domain/stock"
	"sync"
)

type MemoryStockRepository struct {
	lock  *sync.RWMutex
	store map[string]*orderpb.Item
}

func NewMemoryStockRepository(lock *sync.RWMutex, store map[string]*orderpb.Item) *MemoryStockRepository {
	return &MemoryStockRepository{
		lock:  &sync.RWMutex{},
		store: make(map[string]*orderpb.Item),
	}
}

var stub = map[string]*orderpb.Item{
	"item_id": {
		ID:       "foo_item",
		Name:     "stub item",
		Quantity: 1000,
		PriceID:  "stub_item_price_id",
	},
}

func (m MemoryStockRepository) GetItems(_ context.Context, ids []string) ([]*orderpb.Item, error) {
	m.lock.RLock() // get 只需要加 读 锁
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
	return nil, domain.NotFoundError{Missing: missing}

}
