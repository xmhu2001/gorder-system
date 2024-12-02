package discovery

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/xmhu2001/gorder-system/common/discovery/consul"
	"math/rand"
	"time"
)

func RegisterToConsul(serviceName string) (func() error, error) {
	ctx := context.Background()
	registry, err := consul.New(viper.GetString("consul.addr"))
	if err != nil {
		return func() error { return nil }, err
	}
	instanceID := GenerateInstanceID(serviceName)
	grpcAddr := viper.Sub(serviceName).GetString("grpc-addr")
	if err := registry.Register(ctx, instanceID, serviceName, grpcAddr); err != nil {
		return func() error { return nil }, err
	}
	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				logrus.Panicf("no heartbeat from %s to registry, err=%v", serviceName, err)
			}
			time.Sleep(1 * time.Second)
		}
	}()
	logrus.WithFields(logrus.Fields{
		"serviceName": serviceName,
		"addr":        grpcAddr,
		"instanceID":  instanceID,
	}).Info("register to consul")
	return func() error {
		return registry.Deregister(ctx, instanceID, serviceName, grpcAddr)
	}, nil
}

func GetServiceAddr(ctx context.Context, serviceName string) (string, error) {
	registry, err := consul.New(viper.GetString("consul.addr"))
	if err != nil {
		return "", err
	}
	addrs, err := registry.Discover(ctx, serviceName)
	if err != nil {
		return "", err
	}
	if len(addrs) == 0 {
		return "", fmt.Errorf("no %s service addr from consul", serviceName)
	}
	i := rand.Intn(len(addrs))
	logrus.Infof("Discovered %d instance of %s, addrs = %v", len(addrs), serviceName, addrs)
	return addrs[i], nil
}
