package processor

import (
	"context"
	"github.com/xmhu2001/gorder-system/common/genproto/orderpb"
)

// stub: 项目开始时为了快速跑通项目的实现
type InmemProcessor struct{}

func NewInmemProcessor() *InmemProcessor {
	return &InmemProcessor{}
}

func (i InmemProcessor) CreatePaymentLink(ctx context.Context, order *orderpb.Order) (string, error) {
	return "inmem-payment-link", nil
}
