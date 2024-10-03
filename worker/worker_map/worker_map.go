package worker_map

import (
	"chatbot/logger"
	"chatbot/worker/worker_msg"
	"context"
	"sync"
)

var workers sync.Map

func Register(name string, ch chan worker_msg.Message) {
	workers.Store(name, ch)
}

func TransferTo(name string, msg worker_msg.Message) {
	if v, ok := workers.Load(name); ok {
		v.(chan worker_msg.Message) <- msg
	} else {
		logger.Errorf(context.Background(), "worker %s not exist, drop msg: %+v", name, msg)
	}
}
