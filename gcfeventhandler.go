package main

import (
	"encoding/json"
	"github.com/nats-io/stan.go"
	"log"
)

type GCFEventHandler struct{
	dbAdapter DatabaseAdapter
}

func NewGCFEventHandler(dbAdapter DatabaseAdapter) *GCFEventHandler {
	return &GCFEventHandler{dbAdapter: dbAdapter}
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
	err = handler.dbAdapter.SaveGCFResult(gcfResult)
	if err != nil {
		log.Printf("cannot save result to database: %v", err)
		return
	}
}

func calculateGCF(number1 int64, number2 int64) int64 {
	if number1 == 0 {
		return number2
	}
	return calculateGCF(number2 % number1, number1)
}