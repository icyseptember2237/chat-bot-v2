package source

import (
	"chatbot/worker/worker_msg"
	"context"
)

type Source interface {
	Receive(ctx context.Context, eventCh chan<- worker_msg.Message)
}

func NewSource(name, typeName string, config map[string]interface{}) Source {
	switch typeName {
	default:
		return nil
	}
}
