package main

import (
	sharesecretgrpc "github.com/bernardosecades/sharesecret/build"
	_ "github.com/bernardosecades/sharesecret/cmd"
	"github.com/stretchr/testify/assert"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"context"
	"log"
	"net"
	"testing"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {

	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()

	sharesecretgrpc.RegisterSecretAppServer(s, App())

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

// READ: https://codegangsta.gitbooks.io/building-web-apps-with-go/content/testing/end_to_end/index.html
func TestCreateAndSeeSecretWithoutPassword(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := sharesecretgrpc.NewSecretAppClient(conn)
	resp1, err1 := client.CreateSecret(ctx, &sharesecretgrpc.CreateSecretRequest{Content: "This is my secret"})
	if err1 != nil {
		t.Fatalf("CreateSecret failed: %v", err1)
	}

	assert.Len(t, resp1.GetId(), 36)

	resp2, err2 := client.SeeSecret(ctx, &sharesecretgrpc.SeeSecretRequest{Id: resp1.GetId()})

	if err2 != nil {
		t.Fatalf("CreateSecret failed: %v", err1)
	}

	assert.Equal(t, "This is my secret", resp2.GetContent())
}

func TestCreateAndSeeSecretWithPassword(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := sharesecretgrpc.NewSecretAppClient(conn)
	resp1, err1 := client.CreateSecret(ctx, &sharesecretgrpc.CreateSecretRequest{Content: "This is my secret", Password: "1234"})
	if err1 != nil {
		t.Fatalf("CreateSecret failed: %v", err1)
	}

	assert.Len(t, resp1.GetId(), 36)

	resp2, err2 := client.SeeSecret(ctx, &sharesecretgrpc.SeeSecretRequest{Id: resp1.GetId(), Password: "1234"})

	if err2 != nil {
		t.Fatalf("CreateSecret failed: %v", err1)
	}

	assert.Equal(t, "This is my secret", resp2.GetContent())
}
