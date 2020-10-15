package grpc_proxy_router

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"src/gatewayProject/dao"
	"src/gatewayProject/grpc_proxy_middleware"
	"src/gatewayProject/proxy"
	"src/gatewayProject/reverse_proxy"
)

var grpcServerList []*warpGrpcServer

type warpGrpcServer struct {
	Addr string
	*grpc.Server
}

func GrpcServerRun() {
	serviceList := dao.ServiceManagerHandler.GetGrpcServiceList()
	for _, serviceItem := range serviceList {
		tempItem := serviceItem
		go func(serviceDetail *dao.ServiceDetail) {
			addr := fmt.Sprintf(":%d", serviceDetail.GRPCRule.Port)
			rb, err := dao.LoadBalancerHandler.GetLoadBalancer(serviceDetail)
			if err != nil {
				log.Fatalf(" [INFO] GetTcpLoadBalancer %v err:%v\n", addr, err)
				return
			}

			lis, err := net.Listen("tcp", addr)
			if err != nil {
				log.Fatalf(" [INFO] GrpcListen %v err:%v\n", addr, err)
			}

			grpcHandler := reverse_proxy.NewGrpcLoadBalanceHandler(rb) // proxy.TransparentHandler
			s := grpc.NewServer(
				grpc.ChainStreamInterceptor(
					grpc_proxy_middleware.GrpcFlowCountMiddleware(serviceDetail),
					grpc_proxy_middleware.GrpcFlowLimitMiddleware(serviceDetail),
					grpc_proxy_middleware.GrpcJwtAuthTokenMiddleware(serviceDetail),
					grpc_proxy_middleware.GrpcJwtFlowCountMiddleware(serviceDetail),
					grpc_proxy_middleware.GrpcJwtFlowLimitMiddleware(serviceDetail),
					grpc_proxy_middleware.GrpcWhiteListMiddleware(serviceDetail),
					grpc_proxy_middleware.GrpcBlackListMiddleware(serviceDetail),
					grpc_proxy_middleware.GrpcHeaderTransferMiddleware(serviceDetail),
				),
				grpc.CustomCodec(proxy.Codec()),
				grpc.UnknownServiceHandler(grpcHandler),
			)

			grpcServerList = append(grpcServerList, &warpGrpcServer{
				Addr:   addr,
				Server: s,
			})
			log.Printf(" [INFO] grpc_proxy_run %v\n", addr)

			// Serve accepts incoming connections on the listener lis,
			// creating a new ServerTransport and service goroutine for each.
			if err := s.Serve(lis); err != nil {
				log.Fatalf(" [INFO] grpc_proxy_run %v err:%v\n", addr, err)
			}
		}(tempItem)
	}
}

func GrpcServerStop() {
	for _, grpcServer := range grpcServerList {
		grpcServer.GracefulStop()
		log.Printf(" [INFO] grpc_proxy_stop %v stopped\n", grpcServer.Addr)
	}
}
