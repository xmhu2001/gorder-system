package decorator

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type MetricsClient interface {
	Inc(key string, value int)
}

type queryMetricsDecorator[C, R any] struct {
	base   QueryHandler[C, R]
	client MetricsClient
}

func (q queryMetricsDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	start := time.Now()
	actionName := strings.ToLower(generateActionName(cmd))
	defer func() {
		end := time.Since(start)
		q.client.Inc(fmt.Sprintf("query %s duration", actionName), int(end.Seconds()))
	}()
	return q.base.Handle(ctx, cmd)
}

type commandMetricsDecorator[C, R any] struct {
	base   CommandHandler[C, R]
	client MetricsClient
}

func (q commandMetricsDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	start := time.Now()
	actionName := strings.ToLower(generateActionName(cmd))
	defer func() {
		end := time.Since(start)
		q.client.Inc(fmt.Sprintf("command %s duration", actionName), int(end.Seconds()))
	}()
	return q.base.Handle(ctx, cmd)
}
