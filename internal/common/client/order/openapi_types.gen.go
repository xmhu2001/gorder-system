// Package order provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
package order

// CreateOrderRequest defines model for CreateOrderRequest.
type CreateOrderRequest struct {
	CustomerID string             `json:"CustomerID"`
	Items      []ItemWithQuantity `json:"Items"`
}

// Error defines model for Error.
type Error struct {
	Message *string `json:"message,omitempty"`
}

// Item defines model for Item.
type Item struct {
	Id       *string `json:"id,omitempty"`
	Name     *string `json:"name,omitempty"`
	PriceID  *string `json:"priceID,omitempty"`
	Quantity *int32  `json:"quantity,omitempty"`
}

// ItemWithQuantity defines model for ItemWithQuantity.
type ItemWithQuantity struct {
	Id       *string `json:"id,omitempty"`
	Quantity *int32  `json:"quantity,omitempty"`
}

// Order defines model for Order.
type Order struct {
	CustomerID  *string `json:"customerID,omitempty"`
	ID          *string `json:"ID,omitempty"`
	Items       *[]Item `json:"items,omitempty"`
	PaymentLink *string `json:"paymentLink,omitempty"`
	Status      *string `json:"status,omitempty"`
}

// PostCustomerCustomerIDOrdersJSONRequestBody defines body for PostCustomerCustomerIDOrders for application/json ContentType.
type PostCustomerCustomerIDOrdersJSONRequestBody = CreateOrderRequest

