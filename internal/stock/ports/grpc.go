package ports

import (
	context "context"
	"github.com/sirupsen/logrus"
	"github.com/xmhu2001/gorder-system/common/genproto/orderpb"
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
	logrus.Info("rpc_request_in, stock.GetItems")
	return nil, nil
}

func (G GRPCServer) CheckIfItemsInStock(ctx context.Context, request *stockpb.CheckIfItemsInStockRequest) (*stockpb.CheckIfItemsInStockResponse, error) {
	logrus.Info("rpc_request_in, stock.CheckIfItemsInStock")
	defer func() {
		logrus.Info("rpc_request_out, stock.CheckIfItemsInStock")
	}()
	fakeData := []*orderpb.Item{
		{
			ID: "fake-id-CheckIfItemsInStock",
		},
	}
	return &stockpb.CheckIfItemsInStockResponse{Items: fakeData}, nil
}
