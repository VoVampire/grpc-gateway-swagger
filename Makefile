PROTOCHEAD=protoc -I/usr/local/include -I. -I${GOPATH}/src -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis

all: proto gw swag

proto:
	@$(PROTOCHEAD) --go_out=plugins=grpc:. ./protos/*.proto

gw:
	@$(PROTOCHEAD) --grpc-gateway_out=logtostderr=true:. ./protos/*.proto

swag:
	@$(PROTOCHEAD) --swagger_out=logtostderr=true:. ./protos/*.proto