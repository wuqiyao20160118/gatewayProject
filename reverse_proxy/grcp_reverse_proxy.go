package reverse_proxy

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"src/gatewayProject/proxy"
	"src/gatewayProject/reverse_proxy/load_balance"
)

func NewGrpcLoadBalanceHandler(lb load_balance.LoadBalance) grpc.StreamHandler {
	return func() grpc.StreamHandler {
		nextAddr, err := lb.Get("")
		if err != nil {
			log.Fatal("get next addr fail")
		}
		director := func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {
			c, err := grpc.DialContext(ctx, nextAddr, grpc.WithCodec(proxy.Codec()), grpc.WithInsecure())

			// FromIncomingContext returns the incoming metadata in ctx if it exists.
			// The returned MD should not be modified. Writing to it may cause races.
			// Modification should be made to copies of the returned MD.
			md, _ := metadata.FromIncomingContext(ctx)

			// WithCancel returns a copy of parent with a new Done channel.
			// The returned context's Done channel is closed when the returned cancel function is called
			// or when the parent context's Done channel is closed, whichever happens first.
			// Canceling this context releases resources associated with it
			// so code should call cancel as soon as the operations running in this Context complete.
			outCtx, _ := context.WithCancel(ctx)
			outCtx = metadata.NewOutgoingContext(outCtx, md.Copy())
			return outCtx, c, err
		}

		// TransparentHandler returns a handler that attempts to proxy all requests that are not registered in the server.
		// The indented use here is as a transparent proxy
		// where the server doesn't know about the services implemented by the backends.
		// It should be used as a `grpc.UnknownServiceHandler`.
		return proxy.TransparentHandler(director)
	}()
}
