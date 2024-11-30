package app

import (
	"github.com/xmhu2001/gorder-system/order/app/command"
	"github.com/xmhu2001/gorder-system/order/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateOrder command.CreateOrderHandler
	UpdateOrder command.UpdateOrderHandler
}

type Queries struct {
	GetCustomerOrder query.GetCustomerOrderHandler // 注入 handler
}
