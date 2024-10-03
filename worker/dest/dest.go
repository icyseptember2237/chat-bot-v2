package dest

import (
	"chatbot/worker/worker_msg"
	"context"
)

type Dest interface {
	Send(ctx context.Context, eventCh <-chan worker_msg.Message)
}

func NewDest(name, typeName string, config map[string]interface{}) Dest {
	switch typeName {
	default:
		return nil
	}
}
