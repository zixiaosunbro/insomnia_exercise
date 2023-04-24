package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"insomnia/src/pkg/app"
	"insomnia/src/services/rule"
	"insomnia/src/services/rule/handler"
	pb "insomnia/src/services/rule/proto/rule"
	"net"
)

var (
	port = flag.Int("port", 50051, "The service port")
)

func main() {
	server := app.NewServer("proto/rule.proto")
	pb.RegisterRuleServiceServer(server, handler.NewRuleService())
	rule.ResourceInit()
	log.Info("server start success")
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		panic("failed to listen: " + err.Error())
	}
	if err = server.Serve(lis); err != nil {
		panic("failed to serve: " + err.Error())
	}
}
