package grpc_proxy_middleware

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"src/gatewayProject/dao"
	"strings"
)

func GrpcHeaderTransferMiddleware(serviceDetail *dao.ServiceDetail) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// FromIncomingContext returns the incoming metadata in ctx if it exists.
		// The returned MD should not be modified. Writing to it may cause races.
		// Modification should be made to copies of the returned MD.
		md, ok := metadata.FromIncomingContext(ss.Context()) // type MD map[string][]string
		if !ok {
			return errors.New("miss metadata from context")
		}

		for _, item := range strings.Split(serviceDetail.GRPCRule.HeaderTransfor, ",") {
			items := strings.Split(item, " ")
			if len(items) != 3 {
				continue
			}
			if items[0] == "add" || items[0] == "edit" {
				md.Set(items[1], items[2])
			}
			if items[0] == "del" {
				delete(md, items[1])
			}
		}
		if err := ss.SetHeader(md); err != nil {
			return errors.WithMessage(err, "SetHeader")
		}

		if err := handler(srv, ss); err != nil {
			log.Printf("GrpcHeaderTransferMiddleware failed with error %v\n", err)
			return err
		}
		return nil
	}
}
