package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xmhu2001/gorder-system/common/genproto/orderpb"
	"github.com/xmhu2001/gorder-system/order/app"
	"github.com/xmhu2001/gorder-system/order/app/command"
	"github.com/xmhu2001/gorder-system/order/app/query"
	"net/http"
)

type HTTPServer struct {
	app app.Application
}

func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {
	var req orderpb.CreateOrderRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	r, err := H.app.Commands.CreateOrder.Handle(c, command.CreateOrder{
		CustomerID: req.CustomerID,
		Items:      req.Items,
	})
	if err != nil {
		c.JSON(200, gin.H{"error": err})
		return
	}
	c.JSON(200, gin.H{"message": "success", "CustomerID": req.CustomerID, "OrderID": r.OrderID, "redirect_url": fmt.Sprintf("http://localhost:8080/success?customerID=%s&orderID=%s", req.CustomerID, r.OrderID)})
}

func (H HTTPServer) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string) {
	o, err := H.app.Queries.GetCustomerOrder.Handle(c, query.GetCustomerOrder{
		OrderID:    orderID,
		CustomerID: customerID,
	})
	if err != nil {
		c.JSON(200, gin.H{"error": err})
		return
	}
	c.JSON(200, gin.H{
		"msg": "success",
		"data": gin.H{
			"Order": o,
		},
	})
}
