package main

import (
	"fmt"
	client "github.com/xmhu2001/gorder-system/common/client/order"
	"github.com/xmhu2001/gorder-system/order/convertor"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xmhu2001/gorder-system/common/tracing"
	"github.com/xmhu2001/gorder-system/order/app"
	"github.com/xmhu2001/gorder-system/order/app/command"
	"github.com/xmhu2001/gorder-system/order/app/query"
)

type HTTPServer struct {
	app app.Application
}

func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {
	ctx, span := tracing.Start(c, "PostCustomerCustomerIDOrders")
	defer span.End()

	var req client.CreateOrderRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	r, err := H.app.Commands.CreateOrder.Handle(ctx, command.CreateOrder{
		CustomerID: req.CustomerID,
		Items:      convertor.NewItemWithQuantityConvertor().ClientsToEntities(req.Items),
	})
	if err != nil {
		c.JSON(200, gin.H{"error": err})
		return
	}
	c.JSON(200, gin.H{"message": "success", "trace_id": tracing.TraceID(ctx), "customer_id": req.CustomerID, "order_id": r.OrderID, "redirect_url": fmt.Sprintf("http://localhost:8080/success?customerID=%s&orderID=%s", req.CustomerID, r.OrderID)})
}

func (H HTTPServer) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string) {
	ctx, span := tracing.Start(c, "GetCustomerCustomerIDOrdersOrderID")
	defer span.End()

	o, err := H.app.Queries.GetCustomerOrder.Handle(ctx, query.GetCustomerOrder{
		OrderID:    orderID,
		CustomerID: customerID,
	})
	if err != nil {
		c.JSON(200, gin.H{"error": err})
		return
	}
	c.JSON(200, gin.H{
		"message":  "success",
		"trace_id": tracing.TraceID(ctx),
		"data": gin.H{
			"Order": o,
		},
	})
}
