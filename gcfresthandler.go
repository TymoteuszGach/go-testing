package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"strings"
	"time"
)

type GCFRestHandler struct {
	mongoDBClient *mongo.Client
}

func NewGCFRestHandler(mongoDBClient *mongo.Client) *GCFRestHandler {
	return &GCFRestHandler{mongoDBClient: mongoDBClient}
}

func (handler *GCFRestHandler) ListGCFResults(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	collection := handler.mongoDBClient.Database(databaseName).Collection(gcfResultsCollection)
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil { log.Fatal(err) }
	defer func() {
		if err := cur.Close(ctx); err != nil {
			log.Printf("cannot close cursor: %v\n", err)
		}
	}()
	results := make([]string, 0)
	for cur.Next(ctx) {
		var result bson.D
		err := cur.Decode(&result)
		if err != nil {
			log.Printf("cannot decode result: %v", err)
			return
		}
		results = append(results, fmt.Sprintf("%v", result))
	}
	if err := cur.Err(); err != nil {
		log.Printf("cursor error: %v", err)
		return
	}
	_, err = fmt.Fprintf(w, "Results:\n%s", strings.Join(results, "\n"))
	if err != nil {
		log.Printf("cannot return GCF results: %v\n", err)
	}
}