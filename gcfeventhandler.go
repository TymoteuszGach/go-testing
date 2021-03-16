package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nats-io/stan.go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type GCFEventHandler struct{
	mongoDBClient *mongo.Client
}

func NewGCFEventHandler(mongoDBClient *mongo.Client) *GCFEventHandler {
	return &GCFEventHandler{mongoDBClient: mongoDBClient}
}

func (handler *GCFEventHandler) Handle(m *stan.Msg) {
	var request GCFRequest
	err := json.Unmarshal(m.Data, &request)
	if err != nil {
		log.Printf("cannot parse message %s\n", m.Data)
		return
	}

	gcf := calculateGCF(request.Number1, request.Number2)

	gcfResult := GCFResult{
		Number1: request.Number1,
		Number2: request.Number2,
		GCF:     gcf,
	}
	err = handler.saveToDatabase(gcfResult)
	if err != nil {
		log.Printf("cannot save result to database: %v", err)
		return
	}
}

func (handler *GCFEventHandler) saveToDatabase(result GCFResult) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := handler.mongoDBClient.Database(databaseName).Collection(gcfResultsCollection)
	_, err := collection.InsertOne(ctx, bson.D{
		{"number1", result.Number1},
		{"number2", result.Number2},
		{"result", result.GCF},
	})
	if err != nil {
		return fmt.Errorf("cannot insert result: %v", err)
	}
	return nil
}

func calculateGCF(number1 int64, number2 int64) int64 {
	if number1 == 0 {
		return number2
	}
	return calculateGCF(number2 % number1, number1)
}