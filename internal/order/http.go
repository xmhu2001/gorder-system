package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/xmhu2001/gorder-system/common"
	client "github.com/xmhu2001/gorder-system/common/client/order"
	"github.com/xmhu2001/gorder-system/order/app"
	"github.com/xmhu2001/gorder-system/order/app/command"
	"github.com/xmhu2001/gorder-system/order/app/dto"
	"github.com/xmhu2001/gorder-system/order/app/query"
	"github.com/xmhu2001/gorder-system/order/convertor"
)

type HTTPServer struct {
	app app.Application
	common.BaseResponse
}

func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {
	var (
		err  error
		req  client.CreateOrderRequest
		resp dto.CreateOrderResponse
	)
	defer func() {
		H.Response(c, err, &resp)
	}()
	if err = c.ShouldBind(&req); err != nil {
		return
	}
	r, err := H.app.Commands.CreateOrder.Handle(c.Request.Context(), command.CreateOrder{
		CustomerID: req.CustomerID,
		Items:      convertor.NewItemWithQuantityConvertor().ClientsToEntities(req.Items),
	})
	if err != nil {
		return
	}
	//c.JSON(200, gin.H{"message": "success", "trace_id": tracing.TraceID(c.Request.Context()), "customer_id": req.CustomerID, "order_id": r.OrderID, "redirect_url": fmt.Sprintf("http://localhost:8080/success?customerID=%s&orderID=%s", req.CustomerID, r.OrderID)})
	resp = dto.CreateOrderResponse{
		OrderID:     r.OrderID,
		CustomerID:  req.CustomerID,
		RedirectURL: fmt.Sprintf("http://localhost:8080/success?customerID=%s&orderID=%s", req.CustomerID, r.OrderID),
	}
}

func (H HTTPServer) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string) {
	var (
		err  error
		resp struct {
			Order *client.Order
		}
	)
	defer func() {
		H.Response(c, err, &resp)
	}()

	o, err := H.app.Queries.GetCustomerOrder.Handle(c.Request.Context(), query.GetCustomerOrder{
		OrderID:    orderID,
		CustomerID: customerID,
	})
	if err != nil {
		return
	}
	resp.Order = convertor.NewOrderConvertor().EntityToClient(o)
}
