package grpc

import (
	"context"
	"fmt"
	"github.com/json-iterator/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"runtime/debug"
	"time"
)

// RecordRequestInterceptor print request, reply content and the time cost of the request deal
func RecordRequestInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	startAt := time.Now().UnixNano()
	resp, err = handler(ctx, req)
	spanTime := fmt.Sprintf("%.2f", float64(time.Now().UnixNano()-startAt)/float64(time.Millisecond))
	logInfo := map[string]any{
		"method": info.FullMethod,
		"req":    req,
		"span":   spanTime,
		"reply":  resp,
	}
	if err != nil {
		if rpcErr, ok := status.FromError(err); ok {
			logInfo["err"] = fmt.Sprintf("code: %d, msg: %s", rpcErr.Code(), rpcErr.Message())
		} else {
			logInfo["err"] = err
		}
	}
	content, _ := jsoniter.MarshalToString(logInfo)
	// in production mode, the log should be saved in file
	log.Println(content)
	return resp, err
}

func MetricInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	resp, err = handler(ctx, req)
	// TODO: to make troubleshooting easier, should send the grpc code, cost time, method name to prometheus service
	return resp, err
}

func BasicInfoInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	resp, err = handler(ctx, req)
	// TODO: extract the basic info from request, such as user id, device id, etc. and put basic info into context
	return resp, err
}

func HandlerPanicRecover(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			log.Printf("panic err: %v, stack: %v", r, string(debug.Stack()))
			err = status.Errorf(codes.Internal, "panic err: %v", r)
		}
	}()
	return handler(ctx, req)
}

func WarpServerErrorInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if err == nil {
			return
		}

		// convert NeoError into grpc exception
		if rerr, ok := err.(CooError); ok {
			err = SoaErr(rerr)
		} else {
			err = status.Errorf(codes.Code(CodeCodeUnknownError), err.Error())
		}
	}()

	return handler(ctx, req)
}

// ChainUnaryServer is a helper function to chain multiple unary service interceptors into one.
// copy from github.com/grpc-ecosystem/go-grpc-middleware
func ChainUnaryServer(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	n := len(interceptors)

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		chainer := func(currentInter grpc.UnaryServerInterceptor, currentHandler grpc.UnaryHandler) grpc.UnaryHandler {
			return func(currentCtx context.Context, currentReq interface{}) (interface{}, error) {
				return currentInter(currentCtx, currentReq, info, currentHandler)
			}
		}

		chainedHandler := handler
		for i := n - 1; i >= 0; i-- {
			chainedHandler = chainer(interceptors[i], chainedHandler)
		}

		return chainedHandler(ctx, req)
	}
}
