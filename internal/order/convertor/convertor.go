package convertor

import (
	client "github.com/xmhu2001/gorder-system/common/client/order"
	"github.com/xmhu2001/gorder-system/common/genproto/orderpb"
	domain "github.com/xmhu2001/gorder-system/order/domain/order"
	"github.com/xmhu2001/gorder-system/order/entity"
)

// 单例模式

type OrderConvertor struct {
}

type ItemConvertor struct{}

type ItemWithQuantityConvertor struct{}

// []*Foo 不能被 []*interface{} 接收

func (c *OrderConvertor) EntityToProto(o *domain.Order) *orderpb.Order {
	c.check(o)
	return &orderpb.Order{
		ID:          o.ID,
		CustomerID:  o.CustomerID,
		Status:      o.Status,
		Items:       NewItemConvertor().EntitiesToProtos(o.Items),
		PaymentLink: o.PaymentLink,
	}
}

func (c *OrderConvertor) ProtoToEntity(o *orderpb.Order) *domain.Order {
	c.check(o)
	return &domain.Order{
		ID:          o.ID,
		CustomerID:  o.CustomerID,
		Status:      o.Status,
		Items:       NewItemConvertor().ProtosToEntities(o.Items),
		PaymentLink: o.PaymentLink,
	}
}

func (c *OrderConvertor) ClientToEntity(o *client.Order) *domain.Order {
	c.check(o)
	return &domain.Order{
		ID:          o.ID,
		CustomerID:  o.CustomerID,
		Status:      o.Status,
		Items:       NewItemConvertor().ClientsToEntities(o.Items),
		PaymentLink: o.PaymentLink,
	}
}

func (c *OrderConvertor) EntityToClient(o *domain.Order) *client.Order {
	c.check(o)
	return &client.Order{
		ID:          o.ID,
		CustomerID:  o.CustomerID,
		Status:      o.Status,
		Items:       NewItemConvertor().EntitiesToClients(o.Items),
		PaymentLink: o.PaymentLink,
	}
}

func (c *OrderConvertor) check(o interface{}) {
	if o == nil {
		panic("cannot convert nil order")
	}
}

func (c *ItemConvertor) EntitiesToProtos(items []*entity.Item) (res []*orderpb.Item) {
	for _, item := range items {
		res = append(res, c.EntityToProto(item))
	}
	return
}

func (c *ItemConvertor) ProtosToEntities(items []*orderpb.Item) (res []*entity.Item) {
	for _, item := range items {
		res = append(res, c.ProtoToEntity(item))
	}
	return
}

func (c *ItemConvertor) ClientsToEntities(items []client.Item) (res []*entity.Item) {
	for _, item := range items {
		res = append(res, c.ClientToEntity(item))
	}
	return
}

func (c *ItemConvertor) EntitiesToClients(items []*entity.Item) (res []client.Item) {
	for _, item := range items {
		res = append(res, c.EntityToClient(item))
	}
	return
}

func (c *ItemConvertor) EntityToProto(item *entity.Item) *orderpb.Item {
	return &orderpb.Item{
		ID:       item.ID,
		Name:     item.Name,
		Quantity: item.Quantity,
		PriceID:  item.PriceID,
	}
}

func (c *ItemConvertor) ProtoToEntity(item *orderpb.Item) *entity.Item {
	return &entity.Item{
		ID:       item.ID,
		Name:     item.Name,
		Quantity: item.Quantity,
		PriceID:  item.PriceID,
	}
}

func (c *ItemConvertor) ClientToEntity(item client.Item) *entity.Item {
	return &entity.Item{
		ID:       item.ID,
		Name:     item.Name,
		Quantity: item.Quantity,
		PriceID:  item.PriceID,
	}
}

func (c *ItemConvertor) EntityToClient(item *entity.Item) client.Item {
	return client.Item{
		ID:       item.ID,
		Name:     item.Name,
		PriceID:  item.PriceID,
		Quantity: item.Quantity,
	}
}

func (c *ItemWithQuantityConvertor) EntitiesToProtos(items []*entity.ItemWithQuantity) (res []*orderpb.ItemWithQuantity) {
	for _, item := range items {
		res = append(res, c.EntityToProto(item))
	}
	return
}

func (c *ItemWithQuantityConvertor) ProtosToEntities(items []*orderpb.ItemWithQuantity) (res []*entity.ItemWithQuantity) {
	for _, item := range items {
		res = append(res, c.ProtoToEntity(item))
	}
	return
}

func (c *ItemWithQuantityConvertor) EntityToProto(item *entity.ItemWithQuantity) *orderpb.ItemWithQuantity {
	return &orderpb.ItemWithQuantity{
		ID:       item.ID,
		Quantity: item.Quantity,
	}
}

func (c *ItemWithQuantityConvertor) ProtoToEntity(item *orderpb.ItemWithQuantity) *entity.ItemWithQuantity {
	return &entity.ItemWithQuantity{
		ID:       item.ID,
		Quantity: item.Quantity,
	}
}

func (c *ItemWithQuantityConvertor) ClientsToEntities(items []client.ItemWithQuantity) (res []*entity.ItemWithQuantity) {
	for _, item := range items {
		res = append(res, c.ClientToEntity(item))
	}
	return
}

func (c *ItemWithQuantityConvertor) ClientToEntity(item client.ItemWithQuantity) *entity.ItemWithQuantity {
	return &entity.ItemWithQuantity{
		ID:       item.ID,
		Quantity: item.Quantity,
	}
}
