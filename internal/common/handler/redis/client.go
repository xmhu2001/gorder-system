package redis

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"time"
)

func SetNX(ctx context.Context, client *redis.Client, key, value string, ttl time.Duration) (err error) {
	now := time.Now()
	defer func() {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"start": now,
			"key":   key,
			"value": value,
			"err":   err,
			"cost":  time.Since(now).Nanoseconds(),
		}).Info("redis_setnx")
	}()
	if client == nil {
		return errors.New("redis_setnx: client is nil")
	}
	_, err = client.SetNX(ctx, key, value, ttl).Result()
	return err
}

func Del(ctx context.Context, client *redis.Client, key string) (err error) {
	now := time.Now()
	defer func() {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"start": now,
			"key":   key,
			"err":   err,
			"cost":  time.Since(now).Nanoseconds(),
		}).Info("redis_del")
	}()
	if client == nil {
		return errors.New("redis_del: client is nil")
	}
	_, err = client.Del(ctx, key).Result()
	return err
}
