package main

import (
	"context"
	"github.com/nats-io/stan.go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

type GCFRequest struct {
	Number1 int64 `json:"number_1"`
	Number2 int64 `json:"number_2"`
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

	databaseAdapter := NewMongoDBAdapter(mongoDBClient)

	gcfEventHandler := NewGCFEventHandler(databaseAdapter)

	sub, err := sc.Subscribe(subject, gcfEventHandler.Handle)
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

	gcfRESTHandler := NewGCFRestHandler(databaseAdapter)

	http.HandleFunc("/gcf", gcfRESTHandler.ListGCFResults)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
