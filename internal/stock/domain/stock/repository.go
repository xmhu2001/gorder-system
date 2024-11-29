package stock

import (
	"context"
	"fmt"
	"github.com/xmhu2001/gorder-system/common/genproto/orderpb"
	"strings"
)

type Repository interface {
	GetItems(ctx context.Context, ids []string) ([]*orderpb.Item, error)
}

type NotFoundError struct {
	Missing []string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("stock item not found: %s", strings.Join(e.Missing, ","))
}
