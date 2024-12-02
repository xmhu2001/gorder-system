package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/xmhu2001/gorder-system/common/config"
	"github.com/xmhu2001/gorder-system/common/logging"
	"github.com/xmhu2001/gorder-system/common/server"
)

func init() {
	logging.Init()
	if err := config.NewViperConfig(); err != nil {
		logrus.Fatal(err)
	}
}

func main() {
	serverType := viper.GetString("payment.server-to-run")
	paymentHandler := NewPaymentHandler()
	switch serverType {
	case "http":
		server.RunHTTPServer(viper.GetString("payment.service-name"), paymentHandler.RegisterRoutes)
	case "grpc":
		logrus.Panic("unsupported type grpc") // 还会一层一层执行调用栈内地一些defer函数，Fatal则是直接退出
	default:
		logrus.Panic("unsupported type")
	}
}
