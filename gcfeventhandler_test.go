package main

import (
	"github.com/nats-io/stan.go"
	"github.com/nats-io/stan.go/pb"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestGCFEventHandler_Handle_InvalidMessage(t *testing.T) {
	// arrange
	messageBytes := []byte("{invalid message}")
	message := &stan.Msg{
		MsgProto: pb.MsgProto{
			Data: messageBytes,
		},
		Sub: nil,
	}
	onGCFRequestCalled := false
	onGCFRequest := func(GCFRequest) error {
		onGCFRequestCalled = true
		return nil
	}
	handler := NewGCFEventHandler(onGCFRequest)

	// act
	handler.Handle(message)

	// assert
	assert.False(t, onGCFRequestCalled)
}

func TestGCFEventHandler_Handle_ValidMessage(t *testing.T) {
	// arrange
	messageBytes := []byte(`{"number_1": 12, "number_2": 45}`)
	message := &stan.Msg{
		MsgProto: pb.MsgProto{
			Data: messageBytes,
		},
		Sub: nil,
	}
	var actualRequest GCFRequest
	onGCFRequest := func(request GCFRequest) error {
		actualRequest = request
		return nil
	}
	handler := NewGCFEventHandler(onGCFRequest)

	// act
	handler.Handle(message)

	// assert
	expectedRequest := GCFRequest{
		Number1: 12,
		Number2: 45,
	}
	assert.True(t, reflect.DeepEqual(expectedRequest, actualRequest))
}
