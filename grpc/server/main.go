package server

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
	"log"
	"net"
	pb "src/gatewayProject/grpc/proto"
)

var port = flag.Int("port", 50055, "the port to serve on")

const (
	streamingCount = 10
)

// server 需要实现：
//	UnaryEcho(context.Context, *EchoRequest) (*EchoResponse, error)
//	ServerStreamingEcho(*EchoRequest, Echo_ServerStreamingEchoServer) error
//	ClientStreamingEcho(Echo_ClientStreamingEchoServer) error
//	BidirectionalStreamingEcho(Echo_BidirectionalStreamingEchoServer) error
type server struct{}

func (s *server) UnaryEcho(ctx context.Context, in *pb.EchoRequest) (*pb.EchoResponse, error) {
	fmt.Printf("--- UnaryEcho ---\n")

	// FromIncomingContext returns the incoming metadata in ctx if it exists.
	// The returned MD should not be modified.
	// Writing to it may cause races.
	// Modification should be made to copies of the returned MD.
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("miss metadata from context")
	}

	fmt.Println("md", md)
	fmt.Printf("request received: %v, sending echo\n", in)
	return &pb.EchoResponse{Message: in.Message}, nil
}

func (s *server) ClientStreamingEcho(stream pb.Echo_ClientStreamingEchoServer) error {
	fmt.Printf("--- ClientStreamingEcho ---\n")
	// Read requests and send responses.
	var message string
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			fmt.Printf("echo last received message\n")
			return stream.SendAndClose(&pb.EchoResponse{Message: message})
		}

		message = in.Message
		fmt.Printf("request received: %v, building echo\n", in)
		if err != nil {
			return err
		}
	}

}

func (s *server) ServerStreamingEcho(in *pb.EchoRequest, stream pb.Echo_ServerStreamingEchoServer) error {
	fmt.Printf("--- ServerStreamingEcho ---\n")
	fmt.Printf("request received: %v\n", in)
	// Read requests and send responses.
	for i := 0; i < streamingCount; i++ {
		fmt.Printf("echo message %v\n", in.Message)
		err := stream.Send(&pb.EchoResponse{Message: in.Message})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *server) BidirectionalStreamingEcho(stream pb.Echo_BidirectionalStreamingEchoServer) error {
	fmt.Printf("--- BidirectionalStreamingEcho ---\n")
	// Read requests and send responses.
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		fmt.Printf("request received %v, sending echo\n", in)
		if err := stream.Send(&pb.EchoResponse{Message: in.Message}); err != nil {
			return err
		}
	}
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("server listening at %v\n", lis.Addr())
	s := grpc.NewServer()
	// s.RegisterService(&_Echo_serviceDesc, srv)
	// RegisterService registers a service and its implementation to the gRPC server.
	pb.RegisterEchoServer(s, &server{})
	_ = s.Serve(lis)
}
