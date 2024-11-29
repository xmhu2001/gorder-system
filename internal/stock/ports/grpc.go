package ports

import (
	context "context"
	"github.com/xmhu2001/gorder-system/stock/app"

	"github.com/xmhu2001/gorder-system/common/genproto/stockpb"
)

type GRPCServer struct {
	app app.Application // app注入到GRPCServer，然后生成新的构造函数 NewGRPCServer
}

func NewGRPCServer(app app.Application) *GRPCServer {
	return &GRPCServer{app: app}
}

func (G GRPCServer) GetItems(ctx context.Context, request *stockpb.GetItemsRequest) (*stockpb.GetItemsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (G GRPCServer) CheckIfItemsInStock(ctx context.Context, request *stockpb.CheckIfItemsInStockRequest) (*stockpb.CheckIfItemsInStockResponse, error) {
	//TODO implement me
	panic("implement me")
}
