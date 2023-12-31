package grpc_test

import (
	"context"
	"fmt"
	"io"
	"net"
	"testing"

	"github.com/kczapla/go-unit-testing-examples/api/protobufs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	grpcGreeter "github.com/kczapla/go-unit-testing-examples/pkg/grpc"
)

type CleanUpFunc func()

func setupGRPC() (protobufs.GreeterClient, CleanUpFunc) {
	grpcServer := grpc.NewServer()
	serviceServer := grpcGreeter.NewServer()
	protobufs.RegisterGreeterServer(grpcServer, serviceServer)

	listener := bufconn.Listen(1024 * 1024)
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			panic("server did not shutdown gracefully")
		}
	}()

	dialOptions := []grpc.DialOption{
		grpc.WithContextDialer(
			func(context.Context, string) (net.Conn, error) {
				return listener.Dial()
			},
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	grpcServerConnection, err := grpc.Dial(
		"",
		dialOptions...,
	)

	if err != nil {
		panic(fmt.Sprintf("failure on DialContext: %s", err))
	}

	grpcClinet := protobufs.NewGreeterClient(grpcServerConnection)

	return grpcClinet, func() {
		grpcServerConnection.Close()
		grpcServer.Stop()
		listener.Close()
	}
}

func TestUnaryHelloOkResponse(t *testing.T) {
	grpcClient, cleanUp := setupGRPC()
	defer cleanUp()

	resp, err := grpcClient.UnaryHello(context.Background(), &protobufs.UnaryHelloRequest{Name: "foo"})

	if err != nil {
		t.Fatalf("unary hello returned error: %s", err)
	}

	if resp.Code != 1 {
		t.Errorf("expected code 1 but got %d", resp.Code)
	}

	if resp.Message != "ok" {
		t.Fatalf("expected message 'ok' but got %s", resp.Message)
	}
}

func TestUnaryHelloValidateResponse(t *testing.T) {
	testCases := []struct {
		name     string
		errormsg string
	}{
		{
			"",
			"name can not be empty",
		},
		{
			"test1",
			"name can not be 'test1'",
		},
		{
			"test2",
			"name can not be 'test2'",
		},
		{
			"test3",
			"name can not be 'test3'",
		},
	}

	grpcClient, cleanUp := setupGRPC()
	defer cleanUp()

	for _, testCase := range testCases {
		resp, err := grpcClient.UnaryHello(
			context.Background(),
			&protobufs.UnaryHelloRequest{Name: testCase.name},
		)

		if err != nil {
			t.Fatalf("unary hello returned error: %s", err)
		}

		if resp.Code != 2 {
			t.Errorf("expected code 1 but got %d", resp.Code)
		}

		if resp.Message != testCase.errormsg {
			t.Fatalf("expected message 'ok' but got %s", resp.Message)
		}
	}
}

func TestServerStreamingHello(t *testing.T) {
	grpcClient, cleanUp := setupGRPC()
	defer cleanUp()

	stream, err := grpcClient.ServerStreamingHello(context.Background(), &protobufs.ServerStreamingHelloRequest{Name: "test"})
	if err != nil {
		t.Fatalf("server streaming hello returned error: %s", err)
	}

	testCases := []string{
		"test1",
		"test2",
		"test3",
		"test4",
		"test5",
	}

	recvMessages := make([]string, 0)
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			t.Fatalf("%v.ServerStreamingHello(_) = _, %v", grpcClient, err)
		}

		recvMessages = append(recvMessages, resp.Message)
	}

	if len(testCases) != len(recvMessages) {
		t.Fatalf("expected to recv %d messages but got %d", len(testCases), len(recvMessages))
	}

	for i := range testCases {
		if testCases[i] != recvMessages[i] {
			t.Fatalf("expected message %s but got %s", testCases[i], recvMessages[i])
		}
	}
}

func TestClientStreamingHello(t *testing.T) {
	grpcClient, cleanUp := setupGRPC()
	defer cleanUp()

	stream, err := grpcClient.ClientStreamingHello(context.Background())
	if err != nil {
		t.Fatalf("client streaming hello returned error: %s", err)
	}

	testCases := []string{
		"test1",
		"test2",
		"test3",
		"test4",
		"test5",
	}

	for _, testCase := range testCases {
		req := &protobufs.ClientStreamingHelloRequest{Name: testCase}
		if err := stream.Send(req); err != nil {
			t.Fatalf("%v.Send(%v) = %v", stream, req, err)
		}
	}

	reply, err := stream.CloseAndRecv()
	if err != nil {
		t.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
	}

	if reply.GetMessage() != "ok" {
		t.Fatalf("expected reply message to be 'ok' but got '%s'", reply.GetMessage())
	}
}

func TestBidirectionalStreamingHello(t *testing.T) {
	grpcClient, cleanUp := setupGRPC()
	defer cleanUp()

	stream, err := grpcClient.BidirectionalStreamingHello(context.Background())
	if err != nil {
		t.Fatalf("%v.BidrectionalStreamingHello = _, %v", grpcClient, err)
	}

	names := []string{
		"test1",
		"test2",
		"test3",
		"test4",
	}

	for _, name := range names {
		req := protobufs.BidirectionalStreamingHelloRequest{Name: name}
		if err := stream.Send(&req); err != nil {
			t.Fatalf("%v.Send(%v) = %v", stream, name, err)
		}

		resp, err := stream.Recv()
		if err == io.EOF {
			return
		}

		if err != nil {
			t.Fatalf("%v.Recv() = %v", stream, err)
		}

		if resp.Message != "ok" {
			t.Fatalf("resp message was '%s' but expected 'ok'", resp.Message)
		}
	}

	stream.CloseSend()

	_, err = stream.Recv()
	if err != io.EOF {
		t.Fatalf("%v.Recv() = %v", stream, err)
	}
}
