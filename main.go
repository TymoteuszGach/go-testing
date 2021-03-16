package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nats-io/stan.go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	databaseName = "go-testing"
	gcfResultsCollection = "gcf"
)

type GCFRequest struct {
	Number1 int64 `json:"number_1"`
	Number2 int64 `json:"number_2"`
}

type GCFResult struct {
	Number1 int64
	Number2 int64
	GCF 	int64
}

func main(){
	clusterID := "test-cluster"
	clientID := "test-client"
	subject := "gcf"
	mongoDBURI := "mongodb://localhost:27017"

	sc, _ := stan.Connect(clusterID, clientID)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoDBClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoDBURI))
	if err != nil {
		log.Fatalf("cannot initialize mongoDB client: %v", err)
	}

	defer func() {
		if err := mongoDBClient.Disconnect(ctx); err != nil {
			log.Printf("cannot disconnect mongoDB client: %v\n", err)
		}
	}()

	sub, err := sc.Subscribe(subject, handleMessage(mongoDBClient))
	if err != nil {
		log.Fatalf("cannot subscribe to NATS: %v", err)
	}

	defer func() {
		if err := sub.Unsubscribe(); err != nil {
			log.Printf("cannot unsubscribe from NATS: %v\n", err)
		}
		if err := sc.Close(); err != nil {
			log.Printf("cannot close NATS connection: %v\n", err)
		}
	}()

	http.HandleFunc("/gcf", listGCFResults(mongoDBClient))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func listGCFResults(mongoDBClient *mongo.Client) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		collection := mongoDBClient.Database(databaseName).Collection(gcfResultsCollection)
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
}

func handleMessage(mongoDBClient *mongo.Client) func(m *stan.Msg) {
	return func(m *stan.Msg) {
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
		err = saveToDatabase(mongoDBClient, gcfResult)
		if err != nil {
			log.Printf("cannot save result to database: %v", err)
			return
		}
	}
}

func saveToDatabase(mongoDBClient *mongo.Client, result GCFResult) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := mongoDBClient.Database(databaseName).Collection(gcfResultsCollection)
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
