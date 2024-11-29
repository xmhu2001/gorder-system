package order

import (
	"context"
	"fmt"
)

type Repository interface {
	Create(context.Context, *Order) (*Order, error) // 需要订单的主键，因此多返回了它
	Get(ctx context.Context, id, customerID string) (*Order, error)
	Update(
		ctx context.Context,
		o *Order,
		updateFn func(context.Context, *Order) (*Order, error),
	) error
}

type NotFoundError struct {
	OrderID string
}

// 实现内置 error 接口
func (e NotFoundError) Error() string {
	return fmt.Sprintf("order not found: %s", e.OrderID)
}
