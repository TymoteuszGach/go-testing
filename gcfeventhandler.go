package main

import (
	"encoding/json"
	"github.com/nats-io/stan.go"
	"log"
)

type GCFRequest struct {
	Number1 int64 `json:"number_1"`
	Number2 int64 `json:"number_2"`
}

type OnGCFRequest func(GCFRequest) error

type GCFEventHandler struct {
	onGCFRequest OnGCFRequest
}

func NewGCFEventHandler(onGCFRequest OnGCFRequest) *GCFEventHandler {
	return &GCFEventHandler{onGCFRequest: onGCFRequest}
}

func (handler *GCFEventHandler) Handle(m *stan.Msg) {
	var request GCFRequest
	if err := json.Unmarshal(m.Data, &request); err != nil {
		log.Printf("cannot parse message: %s\n", m.Data)
		return
	}

	if err := handler.onGCFRequest(request); err != nil {
		log.Printf("cannot process gcf request: %v\n", err)
		return
	}
}
