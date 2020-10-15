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

func GrpcWhiteListMiddleware(serviceDetail *dao.ServiceDetail) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// Peer contains the information of the peer for an RPC,
		// such as the address and authentication information.
		// FromContext returns the peer information in ctx if it exists.
		peerCtx, ok := peer.FromContext(ss.Context())
		if !ok {
			return errors.New("peer not found with context")
		}

		peerAddr := peerCtx.Addr.String()
		clientIP := peerAddr[0:strings.LastIndex(peerAddr, ":")]

		var ipList []string
		if serviceDetail.AccessControl.WhiteList != "" {
			// strings.Split() 在src字符串为空时会返回[""]，其长度为1，故需要首先判断字符串是否为空
			ipList = strings.Split(serviceDetail.AccessControl.WhiteList, ",")
		}

		if serviceDetail.AccessControl.OpenAuth == 1 && len(ipList) > 0 {
			if !public.InStringSlice(ipList, clientIP) {
				return errors.New(fmt.Sprintf("%s not in white ip list", clientIP))
			}
		}

		if err := handler(srv, ss); err != nil {
			log.Printf("GrpcWhiteListMiddleware failed with error %v\n", err)
			return err
		}
		return nil
	}
}
