package discovery

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type Registry interface {
	Register(ctx context.Context, instanceID, serviceName, hostPort string) error
	Deregister(ctx context.Context, instanceID, serviceName, hostPort string) error // 注销
	Discover(ctx context.Context, serviceName string) ([]string, error)             // 返回 ip 地址数组
	HealthCheck(ctx context.Context, instanceID, serviceName string) error
}

func GenerateInstanceID(serviceName string) string {
	x := rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	return fmt.Sprintf("%s-%d", serviceName, x)
}
