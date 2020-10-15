package grpc_proxy_middleware

import (
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"log"
	"src/gatewayProject/dao"
	"src/gatewayProject/public"
	"strings"
)

// 流量统计
func GrpcFlowLimitMiddleware(serviceDetail *dao.ServiceDetail) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if serviceDetail.AccessControl.ServiceFlowLimit != 0 {
			serviceLimiter, err := public.FlowLimiterHandler.GetLimiter(
				public.FlowServicePrefix+serviceDetail.Info.ServiceName,
				float64(serviceDetail.AccessControl.ServiceFlowLimit))
			if err != nil {
				return err
			}

			// Allow is shorthand for AllowN(time.Now(), 1).
			// AllowN reports whether n events may happen at time now.
			// Use this method if you intend to drop / skip events that exceed the rate limit.
			// Otherwise use Reserve or Wait.
			if !serviceLimiter.Allow() {
				return errors.New(fmt.Sprintf("service flow limit %v", serviceDetail.AccessControl.ServiceFlowLimit))
			}
		}

		// Peer contains the information of the peer for an RPC,
		// such as the address and authentication information.
		// FromContext returns the peer information in ctx if it exists.
		peerCtx, ok := peer.FromContext(ss.Context())
		if !ok {
			return errors.New("peer not found with context")
		}

		peerAddr := peerCtx.Addr.String()
		clientIP := peerAddr[0:strings.LastIndex(peerAddr, ":")]

		if serviceDetail.AccessControl.ClientIPFlowLimit > 0 {
			clientLimiter, err := public.FlowLimiterHandler.GetLimiter(
				public.FlowServicePrefix+serviceDetail.Info.ServiceName+"_"+clientIP,
				float64(serviceDetail.AccessControl.ClientIPFlowLimit))
			if err != nil {
				return err
			}
			if !clientLimiter.Allow() {
				return errors.New(fmt.Sprintf("%v flow limit %v", clientIP, serviceDetail.AccessControl.ClientIPFlowLimit))
			}
		}

		if err := handler(srv, ss); err != nil {
			log.Printf("GrpcFlowLimitMiddleware failed with error %v\n", err)
			return err
		}
		return nil
	}
}
