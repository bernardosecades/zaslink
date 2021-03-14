package grpc

import (
	sharesecretgrpc "github.com/bernardosecades/sharesecret/build"
	"github.com/bernardosecades/sharesecret/repository"
	"github.com/bernardosecades/sharesecret/service"
	"github.com/joho/godotenv"

	"github.com/stretchr/testify/assert"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"context"
	"log"
	"net"
	"os"
	"testing"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Print("Not .env file found")
	}

	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()

	dbName := os.Getenv("DB_NAME")
	dbPass := os.Getenv("DB_PASS")
	dbUser := os.Getenv("DB_USER")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	secretKey := os.Getenv("SECRET_KEY")
	secretPassword := os.Getenv("SECRET_PASSWORD")

	secretRepository := repository.NewMySQLSecretRepository(dbName, dbUser, dbPass, dbHost, dbPort)
	secretService := service.NewSecretService(secretRepository, secretKey, secretPassword)

	sharesecretgrpc.RegisterSecretAppServer(s, NewShareSecretServer(secretService))
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

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
