package adapters

// crud implement

import (
	"context"
	"github.com/sirupsen/logrus"
	domain "github.com/xmhu2001/gorder-system/order/domain/order"
	"strconv"
	"sync"
	"time"
)

type MemoryOrderRepository struct {
	lock  *sync.RWMutex
	store []*domain.Order
}

func NewMemoryOrderRepository() *MemoryOrderRepository {
	s := make([]*domain.Order, 0)
	s = append(s, &domain.Order{
		ID:          "fake-id",
		CustomerID:  "fake-customer-id",
		Status:      "fake-status",
		PaymentLink: "fake-payment-link",
		Items:       nil,
	})
	return &MemoryOrderRepository{
		lock:  &sync.RWMutex{},
		store: s,
	}

}

func (m *MemoryOrderRepository) Create(_ context.Context, order *domain.Order) (*domain.Order, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	newOrder := &domain.Order{
		ID:          strconv.FormatInt(time.Now().Unix(), 10),
		CustomerID:  order.CustomerID,
		Status:      order.Status,
		PaymentLink: order.PaymentLink,
		Items:       order.Items,
	}
	m.store = append(m.store, newOrder)
	logrus.WithFields(logrus.Fields{
		"input order":        order,
		"store_after_create": m.store,
	}).Debug("memory_order_repo_Create")
	return newOrder, nil
}

func (m *MemoryOrderRepository) Get(_ context.Context, id, customerID string) (*domain.Order, error) {
	m.lock.RLock() // get 只需要加 读 锁
	defer m.lock.RUnlock()
	logrus.Debug("store: %v", m.store)
	for _, order := range m.store {
		if order.ID == id && customerID == order.CustomerID {
			logrus.Debug("memory_order_repo_get || found || id=%s || customerID=%s || res=%v", order.ID, order.CustomerID, *order)
			return order, nil
		}
	}
	return nil, domain.NotFoundError{OrderID: id}
}

func (m *MemoryOrderRepository) Update(ctx context.Context, o *domain.Order, updateFn func(context.Context, *domain.Order) (*domain.Order, error)) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	found := false
	for i, order := range m.store {
		if order.ID == o.ID && order.CustomerID == o.CustomerID {
			found = true
			updateOrder, err := updateFn(ctx, o)
			if err != nil {
				return err
			}
			m.store[i] = updateOrder
		}
	}
	if !found {
		return domain.NotFoundError{OrderID: o.ID}
	}
	return nil
}
