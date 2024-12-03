package ports

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/xmhu2001/gorder-system/order/app"
	"github.com/xmhu2001/gorder-system/order/app/command"
	"github.com/xmhu2001/gorder-system/order/app/query"
	domain "github.com/xmhu2001/gorder-system/order/domain/order"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/xmhu2001/gorder-system/common/genproto/orderpb"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ports 和外界交互的入口
// payment 的输入（更新order的指令）从这里进来
type GRPCServer struct {
	app app.Application
}

func NewGRPCServer(app app.Application) *GRPCServer {
	return &GRPCServer{app: app}
}

func (G GRPCServer) CreateOrder(ctx context.Context, request *orderpb.CreateOrderRequest) (*emptypb.Empty, error) {
	_, err := G.app.Commands.CreateOrder.Handle(ctx, command.CreateOrder{
		CustomerID: request.CustomerID,
		Items:      request.Items,
	})
	// grpc 中错误需要用 status.Error 库包一层
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (G GRPCServer) GetOrder(ctx context.Context, request *orderpb.GetOrderRequest) (*orderpb.Order, error) {
	o, err := G.app.Queries.GetCustomerOrder.Handle(ctx, query.GetCustomerOrder{
		CustomerId: request.CustomerID,
		OrderId:    request.OrderID,
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	// handle 返回的是 domain 的 order
	// 这里返回的是 protobuf 的 order
	return o.ToProto(), nil
}

func (G GRPCServer) UpdateOrder(ctx context.Context, request *orderpb.Order) (_ *emptypb.Empty, err error) {
	logrus.Infof("order_grpc||request_in||request=%+v", request)
	order, err := domain.NewOrder(request.ID, request.CustomerID, request.Status, request.PaymentLink, request.Items)
	if err != nil {
		err = status.Error(codes.Internal, err.Error())
		return
	}
	_, err = G.app.Commands.UpdateOrder.Handle(ctx, command.UpdateOrder{
		Order: order,
		UpdateFn: func(ctx context.Context, order *domain.Order) (*domain.Order, error) {
			return order, nil
		},
	})
	return
}
