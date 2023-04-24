package app

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"insomnia/src/pkg/config/middleware"
	_ "insomnia/src/pkg/log"
	grpc2 "insomnia/src/pkg/middleware_custom/grpc"
	"time"
)

func NewServer(protoFile string) *grpc.Server {
	// TODO: init basic components, such as trace
	//set timezone
	if loc, err := time.LoadLocation("Asia/Shanghai"); err == nil {
		time.Local = loc
	}
	middleware.Init()
	serverInterceptors := []grpc.UnaryServerInterceptor{
		grpc2.BasicInfoInterceptor,
		grpc2.RecordRequestInterceptor,
		grpc2.MetricInterceptor,
		grpc2.WarpServerErrorInterceptor,
		grpc2.HandlerPanicRecover,
	}
	interceptor := grpc2.ChainUnaryServer(serverInterceptors...)
	option := grpc.UnaryInterceptor(interceptor)
	server := grpc.NewServer(option)
	reflection.Register(server)
	return server
}
