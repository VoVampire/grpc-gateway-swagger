package main

import (
	"context"
	"flag"
	"fmt"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	hello "grpc-gateway-swagger/protos"
	"grpc-gateway-swagger/swagger"
	"log"
	"net/http"
)

var (
	grpcServerGateway  = flag.String("grpc-server-gateway", ":8081", "gRPC server gateway")
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:9090", "gRPC server endpoint")
)

func serveSwaggerJson(mux *http.ServeMux) {
	fileServer := http.FileServer(http.Dir("/home/qydev/"))
	prefix := "/swagger.json/"
	mux.Handle(prefix, http.StripPrefix(prefix, fileServer))
}

func serveSwaggerUI(mux *http.ServeMux) {
	fileServer := http.FileServer(&assetfs.AssetFS{
		Asset:    swagger.Asset,
		AssetDir: swagger.AssetDir,
		Prefix:   "/home/qydev/go/src/grpc-gateway-swagger/swagger/swagger-ui",
	})
	prefix := "/swagger-ui/"
	mux.Handle(prefix, http.StripPrefix(prefix, fileServer))
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
