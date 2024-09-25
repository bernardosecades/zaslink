//go:build integration

package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bernardosecades/sharesecret/internal/api/handler/secret"
	"github.com/bernardosecades/sharesecret/internal/component"
	"github.com/bernardosecades/sharesecret/internal/entity"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/testcontainers/testcontainers-go"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bernardosecades/sharesecret/internal/repository"
	"github.com/bernardosecades/sharesecret/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestShowSecretHandler(t *testing.T) {
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

	keyWithPwd := "@myPassword111111111111111111111"
	contentPlainText := "this is my secret"

	contentEncrypted, err := component.Encrypt([]byte(keyWithPwd), []byte(contentPlainText))
	assert.NoError(t, err)

	secretRepository := repository.NewMongoDbSecretRepository(ctx, client, DBNameTest)
	secretService := service.NewSecretService(secretRepository, defaultPwd, secretKey)
	secretHandler := secret.NewHandler(secretService)

	// load fixtures
	item := entity.Secret{
		ID:        "854d492d-038e-4900-ba1c-454346f16a61",
		Content:   contentEncrypted,
		CustomPwd: false,
		Viewed:    false,
		CreatedAt: time.Now(),
		ExpiredAt: time.Now(),
	}
	err = secretRepository.SaveSecret(ctx, item)
	assert.NoError(t, err)

	r := mux.NewRouter()
	r.HandleFunc("/secret/{id}", secretHandler.RetrieveSecret)
	req, err := http.NewRequest("GET", fmt.Sprintf("/secret/%s", item.ID), http.NoBody)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var data secret.RetrieveSecretResponse
	err = json.Unmarshal(rr.Body.Bytes(), &data)
	assert.NoError(t, err)
	assert.Equal(t, contentPlainText, data.Content)
}
