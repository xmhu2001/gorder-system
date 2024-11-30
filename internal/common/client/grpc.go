package client

import (
	"context"
	"github.com/spf13/viper"
	"github.com/xmhu2001/gorder-system/common/genproto/stockpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// close是关闭连接的函数
func NewStockGRPCClient(ctx context.Context) (client stockpb.StockServiceClient, close func() error, err error) {
	grpcAddr := viper.GetString("stock.grpc-addr")
	opts, err := grpcDialOpts(grpcAddr)
	if err != nil {
		return nil, nil, err
	}
	conn, err := grpc.NewClient(grpcAddr, opts...)
	if err != nil {
		return nil, nil, err
	}
	return stockpb.NewStockServiceClient(conn), conn.Close, nil
}

func grpcDialOpts(addr string) ([]grpc.DialOption, error) {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}, nil
}
