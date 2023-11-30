package grpc

import (
	"context"
	"io"

	"github.com/kczapla/go-unit-testing-examples/api/protobufs"
)


type greeterServer struct {
	protobufs.UnimplementedGreeterServer
	savedMessages []string
}

func (gs *greeterServer) UnaryHello(_ context.Context, req *protobufs.UnaryHelloRequest) (*protobufs.UnaryHelloResponse, error) {
	switch name := req.GetName(); name {
	    case "":
	        return &protobufs.UnaryHelloResponse{Code: 2, Message: "name can not be empty"}, nil
	    case "test1":
	        return &protobufs.UnaryHelloResponse{Code: 2, Message: "name can not be 'test1'"}, nil
	    case "test2":
	        return &protobufs.UnaryHelloResponse{Code: 2, Message: "name can not be 'test2'"}, nil
	    case "test3":
	        return &protobufs.UnaryHelloResponse{Code: 2, Message: "name can not be 'test3'"}, nil
	    case "test4":
	        return &protobufs.UnaryHelloResponse{Code: 1, Message: "Wow!"}, nil
	    default:
		return &protobufs.UnaryHelloResponse{Code: 1, Message: "ok"}, nil
	}
}

func (gs *greeterServer) ServerStreamingHello(
	_ *protobufs.ServerStreamingHelloRequest,
	stream protobufs.Greeter_ServerStreamingHelloServer) error {

	for _, savedMessage := range gs.savedMessages {
		err := stream.Send(&protobufs.ServerStreamingHelloResponse{Message: savedMessage})
		if err != nil {
			return err
		}
	}
	return nil
}

func (gs *greeterServer) ClientStreamingHello(stream protobufs.Greeter_ClientStreamingHelloServer) error {
	for {
		_, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&protobufs.ClientStreamingHelloResponse{Message: "ok"})
		}
		if err != nil {
			return err
		}
	}
}

func (gs *greeterServer) BidirectionalStreamingHello(stream protobufs.Greeter_BidirectionalStreamingHelloServer) error {
	for {
		_, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		if err := stream.Send(&protobufs.BidirectionalStreamingHelloResponse{Message: "ok"}); err != nil {
			return err
		}
	}
}

func NewServer() *greeterServer {
	return &greeterServer{
		savedMessages: []string{
			"test1",
			"test2",
			"test3",
			"test4",
			"test5",
		},
	}
}
