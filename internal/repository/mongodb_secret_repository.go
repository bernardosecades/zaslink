package repository

import (
	"context"
	"time"

	"github.com/bernardosecades/sharesecret/internal/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const SecretCollectionName = "secrets"

type MongoDbSecretRepository struct {
	database *mongo.Database
}

func NewMongoDbSecretRepository(ctx context.Context, client *mongo.Client, dbName string) *MongoDbSecretRepository {
	err := client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}

	return &MongoDbSecretRepository{database: client.Database(dbName)}
}

// GetSecret return handler if exist and viewed is false
func (r *MongoDbSecretRepository) GetSecret(ctx context.Context, id string) (entity.Secret, error) {
	var result entity.Secret
	filter := bson.M{
		"_id":    id,
		"viewed": false,
		"expiredAt": bson.M{
			"$gte": time.Now().UTC(),
		},
	}
	err := r.database.Collection(SecretCollectionName).FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return entity.Secret{}, err
	}
	return result, nil
}

// SaveSecret insert a handler if not exist or update it if exist
func (r *MongoDbSecretRepository) SaveSecret(ctx context.Context, secret entity.Secret) error {
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"_id": secret.ID}
	update := bson.M{"$set": secret}
	_, err := r.database.Collection(SecretCollectionName).UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}
	return nil
}
