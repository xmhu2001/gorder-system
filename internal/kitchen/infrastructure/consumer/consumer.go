package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/xmhu2001/gorder-system/common/broker"
	"github.com/xmhu2001/gorder-system/common/genproto/orderpb"
	"go.opentelemetry.io/otel"
	"time"
)

type Order struct {
	ID          string
	CustomerID  string
	Status      string
	PaymentLink string
	Items       []*orderpb.Item
}

type OrderService interface {
	UpdateOrder(ctx context.Context, req *orderpb.Order) error
}

type Consumer struct {
	orderGRPC OrderService
}

func NewConsumer(orderGRPC OrderService) *Consumer {
	return &Consumer{
		orderGRPC: orderGRPC,
	}
}

func (c *Consumer) Listen(ch *amqp.Channel) {
	q, err := ch.QueueDeclare("", true, false, true, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}

	err = ch.QueueBind(q.Name, "", broker.EventOrderPaid, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		logrus.Warnf("fail to consume: queue=%s, err=%v", q.Name, err)
	}
	var forever chan struct{}
	go func() {
		for msg := range msgs {
			c.handleMessage(ch, msg, q)
		}
	}()
	<-forever
}

func (c *Consumer) handleMessage(ch *amqp.Channel, msg amqp.Delivery, q amqp.Queue) {
	var err error
	logrus.Infof("Kitchen received a message: %s from %s", string(msg.Body), q.Name)
	ctx := broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers)
	tr := otel.Tracer("rabbitmq")
	_, span := tr.Start(ctx, fmt.Sprintf("rabbitmq.%s.consume", q.Name))

	defer func() {
		span.End()
		if err != nil {
			_ = msg.Nack(false, false)
		} else {
			_ = msg.Ack(false)
		}
	}()

	o := &Order{}
	if err := json.Unmarshal(msg.Body, o); err != nil {
		logrus.Infof("fail to unmarshal msg to order, err=%v", err)
		return
	}
	if o.Status != "paid" {
		err = errors.New("order status not paid, cannot cook")
		return
	}
	cook(o)
	span.AddEvent(fmt.Sprintf("order_cook: %v", o))
	if err = c.orderGRPC.UpdateOrder(ctx, &orderpb.Order{
		ID:          o.ID,
		CustomerID:  o.CustomerID,
		Status:      "ready",
		Items:       o.Items,
		PaymentLink: o.PaymentLink,
	}); err != nil {
		if err = broker.HandleRetry(ctx, ch, &msg); err != nil {
			logrus.Warnf("kitchen: error handling retry, err=%v", err)
		}
		return
	}
	span.AddEvent("kitchen.order.finished.updated")
}

func cook(o *Order) {
	logrus.Printf("cooking order %s", o.ID)
	time.Sleep(5 * time.Second)
	logrus.Printf("order %s cooking done", o.ID)
}
