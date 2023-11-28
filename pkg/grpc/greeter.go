package grpc

import (
	"context"

	"github.com/kczapla/go-unit-testing-examples/api/protobufs"
)


type greeterServer struct {
	protobufs.UnimplementedGreeterServer
}

func (gs *greeterServer) UnaryHello(_ context.Context, req *protobufs.UnaryHelloRequest) (*protobufs.UnaryHelloResponse, error) {
	switch name := req.GetName(); name {
	    case "":
	        return &protobufs.UnaryHelloResponse{Code: 2, Message: "name can not be empty"}, nil
	    case "test1":
	        return &protobufs.UnaryHelloResponse{Code: 2, Message: "name can not be 'test1'"}, nil
	    case "test2":
	        return &protobufs.UnaryHelloResponse{Code: 2, Message: "name can not be 'test1'"}, nil
	    case "test3":
	        return &protobufs.UnaryHelloResponse{Code: 2, Message: "name can not be 'test1'"}, nil
	    case "test4":
	        return &protobufs.UnaryHelloResponse{Code: 1, Message: "Wow!"}, nil
	    default:
		return &protobufs.UnaryHelloResponse{Code: 1, Message: "ok"}, nil
	}
}

func (gs *greeterServer) ServerStreamingHello(*protobufs.ServerStreamingHelloRequest, protobufs.Greeter_ServerStreamingHelloServer) error {
	return nil
}

func (gs *greeterServer) ClientStreamingHello(protobufs.Greeter_ClientStreamingHelloServer) error {
	return nil
}

func (gs *greeterServer) BidirectionalStreamingHello(protobufs.Greeter_BidirectionalStreamingHelloServer) error {
	return nil
}

func NewServer() *greeterServer {
	return &greeterServer{}
}
