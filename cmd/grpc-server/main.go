package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	hello "grpc-swagger/protos"
	"log"
	"net"
)

var (
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:9090", "gRPC server endpoint")
)

type HelloService struct {
}

func (*HelloService) SayHello(c context.Context, r *hello.SayHelloRequest) (*hello.SayHelloResponse, error) {
	return &hello.SayHelloResponse{Message: fmt.Sprintf("Hello %s", r.Name)}, nil
}

func (*HelloService) SayGoodbye(c context.Context, r *hello.SayGoodbyeRequest) (*hello.SayGoodbyeResponse, error) {
	return &hello.SayGoodbyeResponse{Message: fmt.Sprintf("Goodbye %s", r.Name)}, nil
}

func run() error {
	lis, err := net.Listen("tcp", *grpcServerEndpoint)
	if err != nil {
		return err
	}

	grpcService := grpc.NewServer()

	hello.RegisterHelloServiceServer(grpcService, &HelloService{})

	return grpcService.Serve(lis)
}

func main() {
	flag.Parse()

	log.Println(fmt.Sprintf("listen %s", *grpcServerEndpoint))

	if err := run(); err != nil {
		log.Fatalf("failed to run: %v", err)
	}
}
