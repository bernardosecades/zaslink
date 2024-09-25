package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const DBNameTest = "testdb"

func setupMongoContainer(t *testing.T) (testcontainers.Container, *mongo.Client) {
	ctx := context.Background()

	// Set up MongoDB container
	mongoC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mongo:latest",
			ExposedPorts: []string{"27017/tcp"},
			WaitingFor:   wait.ForLog("Waiting for connections"),
		},
		Started: true,
	})
	if err != nil {
		t.Fatalf("Could not start container: %s", err)
	}

	// Get the host and port
	host, err := mongoC.Host(ctx)
	assert.NoError(t, err, "Could not get container host")

	port, err := mongoC.MappedPort(ctx, "27017")
	assert.NoError(t, err, "Could not get port for container")

	mongoURI := fmt.Sprintf("mongodb://%s:%s", host, port.Port())

	// Connect to MongoDB
	client, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	assert.NoError(t, err, "Could not connect to mongo")

	// Wait for the MongoDB container to be ready
	ctxTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err = client.Ping(ctxTimeout, nil)
	assert.NoError(t, err, "Could not ping mongo")

	return mongoC, client
}
