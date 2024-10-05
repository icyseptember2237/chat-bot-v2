package transfer

import (
	"chatbot/logger"
	"chatbot/utils/engine_pool"
	"chatbot/utils/luatool"
	"chatbot/worker/worker_map"
	"chatbot/worker/worker_msg"
	"context"
	"github.com/icyseptember2237/engine"
	"reflect"
	"time"
)

func init() {
	engine_pool.RegisterModule(func(i engine.Engine) {
		i.RegisterModule(moduleName, moduleMethods)
	})
}

const moduleName = "transfer"

var moduleMethods = map[string]interface{}{
	"send": send,
}

func send(workerName string, msg interface{}) bool {
	messages := luatool.ConvertLuaData(msg)
	if reflect.TypeOf(messages) != reflect.TypeOf(map[string]interface{}{}) {
		logger.Errorf(context.Background(), "invalid mes type %v", reflect.TypeOf(messages))
		return false
	}

	m := worker_msg.Message{
		Info: worker_msg.Info{
			To:          workerName,
			Time:        time.Now().Unix(),
			ProcessTime: 0,
			FinishTime:  0,
		},
		Content: worker_msg.Content{
			Data: messages.(map[string]interface{}),
		},
	}
	worker_map.TransferTo(workerName, m)
	return true
}
