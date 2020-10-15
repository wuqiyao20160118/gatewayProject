package grpc_proxy_middleware

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"src/gatewayProject/dao"
	"src/gatewayProject/public"
)

func GrpcJwtFlowCountMiddleware(serviceDetail *dao.ServiceDetail) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			return errors.New("miss metadata from context")
		}

		// 获取租户信息
		appInfos := md.Get("app")
		if len(appInfos) == 0 {
			if err := handler(srv, ss); err != nil {
				log.Printf("GrpcJwtFlowCountMiddleware failed with error %v\n", err)
				return err
			}
			return nil
		}

		appInfo := &dao.App{}
		// json格式转为字符流格式，grpc metadata为json格式数据
		if err := json.Unmarshal([]byte(appInfos[0]), appInfo); err != nil {
			return err
		}

		appCounter, err := public.FlowCounterHandler.GetFlowCounter(public.FlowAppPrefix + appInfo.AppID)
		if err != nil {
			return err
		}
		appCounter.Increase()

		if appInfo.Qpd > 0 && appCounter.TotalCount > appInfo.Qpd {
			return errors.New(fmt.Sprintf("User QPD exceeds the limit:%v current:%v", appInfo.Qpd, appCounter.TotalCount))
		}

		if err := handler(srv, ss); err != nil {
			log.Printf("GrpcJwtFlowCountMiddleware failed with error %v\n", err)
			return err
		}
		return nil
	}
}
