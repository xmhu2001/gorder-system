package redis

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"github.com/xmhu2001/gorder-system/common/handler/factory"
	"time"
)

const (
	confName      = "redis"
	localSupplier = "local"
)

var (
	singleton = factory.NewSingleton(supplier)
)

func Init() {
	conf := viper.GetStringMap(confName)
	for supplyName := range conf {
		_ = Client(supplyName)
	}
}

func LocalClient() *redis.Client {
	return singleton.Get(localSupplier).(*redis.Client)
}

func Client(name string) *redis.Client {
	return singleton.Get(name).(*redis.Client)
}

// supplier: read redis config
func supplier(key string) any {
	confKey := confName + "." + key
	type Section struct {
		IP           string        `mapstructure:"ip"`
		Port         string        `mapstructure:"port"`
		PoolSize     int           `mapstructure:"pool_size"`
		MaxConn      int           `mapstructure:"max_conn"`
		ConnTimeout  time.Duration `mapstructure:"conn_timeout"`
		ReadTimeout  time.Duration `mapstructure:"read_timeout"`
		WriteTimeout time.Duration `mapstructure:"write_timeout"`
	}
	var c Section
	if err := viper.UnmarshalKey(confKey, &c); err != nil {
		panic(err)
	}
	return redis.NewClient(&redis.Options{
		Network:         "tcp",
		Addr:            c.IP + ":" + c.Port,
		PoolSize:        c.PoolSize,
		MaxActiveConns:  c.MaxConn,
		ConnMaxLifetime: c.ConnTimeout * time.Millisecond,
		ReadTimeout:     c.ReadTimeout * time.Millisecond,
		WriteTimeout:    c.WriteTimeout * time.Millisecond,
	})
}
