package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/xmhu2001/gorder-system/common/broker"
	"github.com/xmhu2001/gorder-system/common/genproto/orderpb"
	"github.com/xmhu2001/gorder-system/payment/app"
	"github.com/xmhu2001/gorder-system/payment/app/command"
	"go.opentelemetry.io/otel"
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
			c.handleMessage(ch, msg, q)
		}
	}()
	// 永远读不到，因此会永远阻塞，以读取消息
	<-forever
}

func (c *Consumer) handleMessage(ch *amqp.Channel, msg amqp.Delivery, q amqp.Queue) {
	logrus.Infof("Payment received a message: %s from %s", string(msg.Body), q.Name)
	ctx := broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers)
	tr := otel.Tracer("rabbitmq")
	_, span := tr.Start(ctx, fmt.Sprintf("rabbitmq.%s.consume", q.Name))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			_ = msg.Nack(false, false)
		} else {
			_ = msg.Ack(false)
		}
	}()

	o := &orderpb.Order{}
	if err := json.Unmarshal(msg.Body, o); err != nil {
		logrus.Infof("fail to unmarshal msg to order, err=%v", err)
		return
	}
	logrus.Infof("create payment link for customer: %v success", o)
	if _, err := c.app.Commands.CreatePayment.Handle(ctx, command.CreatePayment{Order: o}); err != nil {
		logrus.Infof("fail to create order's paymentlink, err=%v", err)
		if err = broker.HandleRetry(ctx, ch, &msg); err != nil {
			logrus.Warnf("retry_error, error handling retry, messageID=%s, err=%v", msg.MessageId, err)
		}
		return
	}

	span.AddEvent("payment.created")
	logrus.Info("Payment consume success")
}
