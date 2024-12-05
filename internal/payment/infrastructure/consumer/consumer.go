package consumer

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/xmhu2001/gorder-system/common/broker"
	"github.com/xmhu2001/gorder-system/common/genproto/orderpb"
	"github.com/xmhu2001/gorder-system/payment/app"
	"github.com/xmhu2001/gorder-system/payment/app/command"
)

// mq中的order event由payment消费
// mq,db属于基础设施，因此在 payment 下建立 infrastructure 文件夹
type Consumer struct {
	app app.Application
}

func NewConsumer(app app.Application) *Consumer {
	return &Consumer{
		app: app,
	}
}

func (c *Consumer) Listen(ch *amqp.Channel) {
	q, err := ch.QueueDeclare(broker.EventOrderCreated, true, false, false, false, nil)
	if err != nil {
		logrus.Fatal(err) // 初始化里出错fatal
	}
	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		logrus.Warnf("fail to consume: queue=%s, err=%v", q.Name, err)
	}
	var forever chan struct{}
	go func() {
		for msg := range msgs {
			c.HandleMessage(msg, q, ch)
		}
	}()
	// 永远读不到，因此会永远阻塞，以读取消息
	<-forever
}

func (c *Consumer) HandleMessage(msg amqp.Delivery, q amqp.Queue, ch *amqp.Channel) {
	logrus.Infof("Payment received a message: %s from %s", string(msg.Body), q.Name)

	o := &orderpb.Order{}
	if err := json.Unmarshal(msg.Body, o); err != nil {
		logrus.Infof("fail to unmarshal msg to order, err=%v", err)
		_ = msg.Nack(false, false)
		return
	}
	logrus.Infof("create payment link for customer: %v success", o)
	if _, err := c.app.Commands.CreatePayment.Handle(context.TODO(), command.CreatePayment{
		Order: o,
	}); err != nil {
		// TODO: retry
		logrus.Infof("fail to create order's paymentlink, err=%v", err)
		_ = msg.Nack(false, false)
		return
	}

	_ = msg.Ack(false)
	logrus.Info("Payment consume success")
}
