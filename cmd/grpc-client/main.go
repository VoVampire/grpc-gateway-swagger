package main

import (
	"context"
	"encoding/json"
	"flag"
	"google.golang.org/grpc"
	"gopkg.in/resty.v1"
	hello "grpc-gateway-swagger/protos"
	"log"
)

var (
	grpcServerGateway  = flag.String("grpc-server-gateway", "localhost:8081", "gRPC server gateway")
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:9090", "gRPC server endpoint")
)

func invokeHTTP() error {
	r, err := resty.R().
		SetBody(&hello.SayGoodbyeRequest{Name: "Bob"}).
		Post("http://" + *grpcServerGateway + "/v1/examples/say_goodbye")
	if err != nil {
		return err
	}

	res := hello.SayGoodbyeResponse{}
	if err := json.Unmarshal(r.Body(), &res); err != nil {
		return err
	}

	log.Printf("Greeting: %s", res.Message)
	return nil
}

func invokeRPC() error {
	// new connect
	conn, err := grpc.Dial(*grpcServerEndpoint, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	// say hello
	r, err := hello.NewHelloServiceClient(conn).SayHello(context.Background(), &hello.SayHelloRequest{Name: "Alice"})
	if err != nil {
		return err
	}

	log.Printf("Greeting: %s", r.Message)
	return nil
}

func main() {
	flag.Parse()

	if err := invokeRPC(); err != nil {
		log.Fatalf("invoke rpc err: %v", err)
	}

	if err := invokeHTTP(); err != nil {
		log.Fatalf("invoke http err: %v", err)
	}
}
