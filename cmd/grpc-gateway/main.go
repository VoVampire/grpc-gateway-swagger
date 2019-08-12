package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	hello "grpc-gateway-swagger/protos"
	"log"
	"net/http"
	"path"
)

var (
	grpcServerGateway     = flag.String("grpc-server-gateway", ":8081", "gRPC server gateway")
	grpcServerEndpoint    = flag.String("grpc-server-endpoint", "localhost:9090", "gRPC server endpoint")
	swaggerFileServerPath = flag.String("swagger-file-server-path", "/home/qydev/go/src/grpc-gateway-swagger/swagger/", "swagger file server path")
)

func serveSwaggerJson(mux *http.ServeMux) {
	prefix := "/swagger-json/"
	p := path.Join(*swaggerFileServerPath, "swagger-json")
	mux.Handle(prefix, http.StripPrefix(prefix, http.FileServer(http.Dir(p))))
}

func serveSwaggerUI(mux *http.ServeMux) {
	prefix := "/swagger-ui/"
	p := path.Join(*swaggerFileServerPath, "swagger-ui")
	mux.Handle(prefix, http.StripPrefix(prefix, http.FileServer(http.Dir(p))))
}

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := hello.RegisterHelloServiceHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
	if err != nil {
		return err
	}

	httpMux := http.NewServeMux()
	httpMux.Handle("/", mux)
	serveSwaggerUI(httpMux)
	serveSwaggerJson(httpMux)

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(*grpcServerGateway, httpMux)
}

func main() {
	flag.Parse()

	log.Println(fmt.Sprintf("listen %s", *grpcServerGateway))

	if err := run(); err != nil {
		log.Fatalf("failed to run: %v", err)
	}
}
