package grpc_test

import (
	"context"
	"fmt"
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

	listener := bufconn.Listen(1024*1024)
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			panic("server did not shutdown gracefully")
		}
	}()

	ctx := context.Background()

	dialOptions := []grpc.DialOption{
		grpc.WithContextDialer(
			func(context.Context, string) (net.Conn, error) {
				return listener.Dial()
			},
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	}

	grpcServerConnection, err := grpc.DialContext(
		ctx,
		"",
		dialOptions...
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
