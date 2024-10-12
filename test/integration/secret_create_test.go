//go:build integration

package integration

import (
	"context"
	"encoding/json"
	"github.com/bernardosecades/zaslink/pkg/events"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/testcontainers/testcontainers-go"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/bernardosecades/zaslink/internal/api/handler/secret"
	"github.com/bernardosecades/zaslink/internal/repository"
	"github.com/bernardosecades/zaslink/internal/service"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestCreateSecretHandler(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	container, client := setupMongoContainer(t)
	defer func(container testcontainers.Container, ctx context.Context) {
		err := container.Terminate(ctx)
		if err != nil {
			log.Printf("Error terminating from container: %s", err)
		}
	}(container, ctx)
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			log.Printf("Error disconnecting from MongoDB: %s", err)
		}
	}(client, ctx)

	defaultPwd := "@myPassword"
	secretKey := "11111111111111111111111111111111"

	secretRepository := repository.NewMongoDbSecretRepository(ctx, client, DBNameTest)
	secretService := service.NewSecretService(secretRepository, events.NewInMemoryPublisher(), defaultPwd, secretKey)
	secretHandler := secret.NewHandler(secretService)

	r := mux.NewRouter()
	r.HandleFunc("/secret", secretHandler.CreateSecret)

	secretBody := `{"content": "this is my secret", "pwd": ""}` // #nosec (skip linter G101: Potential hardcoded credentials)
	req, err := http.NewRequest("POST", "/secret", http.NoBody)
	req.Body = io.NopCloser(strings.NewReader(secretBody))

	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var data secret.CreateSecretResponse
	err = json.Unmarshal(rr.Body.Bytes(), &data)
	assert.NoError(t, err)

	assert.NotEmpty(t, data.ID)
	assert.NotEmpty(t, data.ExpiredAt)

	filter := bson.M{"_id": data.ID}
	totalDocuments, err := client.Database(DBNameTest).Collection(repository.SecretCollectionName).CountDocuments(ctx, filter)

	assert.NoError(t, err)
	assert.EqualValues(t, 1, totalDocuments)
}
