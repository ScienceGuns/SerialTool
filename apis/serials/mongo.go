/*
Copyright © 2025 Mad Scientist Research LLC
This file is part of Serial Tool.
*/

package serials

import (
	"context"
	"fmt"
	"github.com/ScienceGuns/SerialTool/apis/internal_funcs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

// Datastore is an interface that defines how we connect to our database for database operations.
type Datastore interface {
	GetSerialCount(config internal_funcs.SerialToolConfig, prefix string) (int, error)
	StoreSerialNumber(config internal_funcs.SerialToolConfig, serial string) error
	GetConfig() internal_funcs.SerialToolConfig
	Disconnect()
}

// Database is the struct that holds our live database connection and config.
// It implements the Datastore interface.
type Database struct {
	client *mongo.Client
	config internal_funcs.SerialToolConfig
	ctx    context.Context
}

// Disconnect closes the connection to the MongoDB server.
func (db *Database) Disconnect() {
	if err := db.client.Disconnect(db.ctx); err != nil {
		log.Printf("Failed to disconnect from the database: %v", err)
	}
	fmt.Println("Disconnected from the database.")
}

// GetConfig returns the configuration stored within the Database object.
func (db *Database) GetConfig() internal_funcs.SerialToolConfig {
	return db.config
}

// GetSerialCount queries MongoDB to find how many USNs with the given prefix exist. (Leave off CCCC!)
func (db *Database) GetSerialCount(config internal_funcs.SerialToolConfig, serialPrefix string) (int, error) {
	// Call the correct Database and collection pair
	coll := db.client.Database(config.MongoDB.Database).Collection(config.MongoDB.Collection)
	// Filter based on the serial_number field and find anything that starts with the prefix
	filter := bson.M{"serial_number": bson.M{"$regex": "^" + serialPrefix}}
	// Simply count the total entries for this USN and return the result
	count, err := coll.CountDocuments(db.ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count documents: %w", err)
	}
	return int(count), nil
}

// StoreSerialNumber inserts a new USN document into the database.
func (db *Database) StoreSerialNumber(config internal_funcs.SerialToolConfig, fullSerial string) error {
	// Call the correct Database and collection pair
	coll := db.client.Database(config.MongoDB.Database).Collection(config.MongoDB.Collection)
	// Create the correct record type
	record := serialRecord{
		SerialNumber: fullSerial,
		CreatedAt:    time.Now(),
	}
	// Attempt to insert the new USN into the database
	_, err := coll.InsertOne(db.ctx, record)
	if err != nil {
		return fmt.Errorf("failed to insert the new USN record: %w", err)
	}
	return nil
}

// InitDB initializes the connection to a MongoDB server.
func InitDB(config internal_funcs.SerialToolConfig) (Datastore, error) {
	ctx := context.Background()
	clientOptions := options.Client().ApplyURI(config.MongoDB.URI)
	// Connect a Mongo client
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping the database: %w", err)
	}

	fmt.Println("Successfully connected to the database.")
	return &Database{
		client: client,
		config: config,
		ctx:    ctx,
	}, nil
}
