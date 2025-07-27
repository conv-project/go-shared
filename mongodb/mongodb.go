package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"github.com/conv-project/conversion-service/pkg/logger"
)

// Database represents MongoDB database connection.
type Database struct {
	client   *mongo.Client
	database *mongo.Database
}

// New creates a new MongoDB connection.
func New(dsn string, dbName string) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(dsn)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database to verify connection
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	logger.Info("Connected to MongoDB", zap.String("dsn", dsn), zap.String("database", dbName))

	return &Database{
		client:   client,
		database: client.Database(dbName),
	}, nil
}

// Close disconnects from MongoDB.
func (db *Database) Close(ctx context.Context) error {
	if db.client != nil {
		if err := db.client.Disconnect(ctx); err != nil {
			return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
		}
		logger.Info("Disconnected from MongoDB")
	}
	return nil
}

// GetCollection returns a MongoDB collection.
func (db *Database) GetCollection(name string) *mongo.Collection {
	return db.database.Collection(name)
}

// GetDatabase returns the MongoDB database.
func (db *Database) GetDatabase() *mongo.Database {
	return db.database
}

func (db *Database) GetClient() *mongo.Client {
	return db.client
}
