package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	databaseName         = "go-testing"
	gcfResultsCollection = "gcf"
)

type MongoDBAdapter struct {
	client *mongo.Client
}

func NewMongoDBAdapter(client *mongo.Client) *MongoDBAdapter {
	return &MongoDBAdapter{client: client}
}

type GCFResult struct {
	Number1 int64 `bson:"number_1"`
	Number2 int64 `bson:"number_2"`
	GCF     int64 `bson:"gcf"`
}

func (result GCFResult) String() string {
	return fmt.Sprintf("number 1: %d, number 2: %d, result: %d", result.Number1, result.Number2, result.GCF)
}

func (adapter *MongoDBAdapter) ListGCFResults() ([]GCFResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	collection := adapter.client.Database(databaseName).Collection(gcfResultsCollection)
	var results []GCFResult
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("cannot find gcf results: %v", err)
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("cannot get all gcf results: %v", err)
	}

	return results, nil
}

func (adapter *MongoDBAdapter) SaveGCFResult(result GCFResult) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := adapter.client.Database(databaseName).Collection(gcfResultsCollection)
	_, err := collection.InsertOne(ctx, result)
	if err != nil {
		return fmt.Errorf("cannot insert result: %v", err)
	}
	return nil
}
