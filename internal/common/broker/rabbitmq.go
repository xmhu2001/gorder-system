package broker

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

// rabbitmq相关初始化
// 这里需要接入rabbitmq的sdk
func Connect(user, password, host, port string) (*amqp.Channel, func() error) {
	address := fmt.Sprintf("amqp://%s:%s@%s:%s", user, password, host, port)
	conn, err := amqp.Dial(address)
	if err != nil {
		logrus.Fatal(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		logrus.Fatal(err)
	}
	// 下面声明 2个 exchange（对应两个event）
	// direct: 点对点
	// fanout: 广播
	err = ch.ExchangeDeclare(EventOrderCreated, "direct", true, false, false, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}
	err = ch.ExchangeDeclare(EventOrderPaid, "fanout", true, false, false, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}
	return ch, conn.Close
	// 然后去 order/main.go 和 payment/main.go 给它初始化
}
