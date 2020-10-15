package grpc_proxy_middleware

import (
	"google.golang.org/grpc"
	"log"
	"src/gatewayProject/dao"
	"src/gatewayProject/public"
)

// 流量统计
func GrpcFlowCountMiddleware(serviceDetail *dao.ServiceDetail) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		counter, err := public.FlowCounterHandler.GetFlowCounter(public.FlowTotal)
		if err != nil {
			return err
		}
		counter.Increase()

		serviceCounter, err := public.FlowCounterHandler.GetFlowCounter(public.FlowServicePrefix + serviceDetail.Info.ServiceName)
		if err != nil {
			return err
		}
		serviceCounter.Increase()

		if err := handler(srv, ss); err != nil {
			log.Printf("GrpcFlowCountMiddleware failed with error %v\n", err)
			return err
		}
		return nil
	}
}
