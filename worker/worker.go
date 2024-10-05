package worker

import (
	"chatbot/logger"
	"chatbot/utils/deep"
	"chatbot/utils/engine_pool"
	"chatbot/utils/luatool"
	"chatbot/worker/config"
	"chatbot/worker/dest"
	"chatbot/worker/source"
	"chatbot/worker/worker_map"
	"chatbot/worker/worker_msg"
	"context"
	"fmt"
	"github.com/icyseptember2237/engine"
	"golang.org/x/time/rate"
	"math"
	"reflect"
	"sync"
	"time"
)

const (
	defaultRateLimit    = 600
	defaultBufferLength = 1024
)

type Worker struct {
	conf config.WorkerConfig

	inCh  chan worker_msg.Message
	outCh chan worker_msg.Message

	source source.Source
	dest   dest.Dest
	engine []engine.Engine

	ready   bool
	running bool
	stopSig chan bool

	rl     *rate.Limiter
	mu     sync.Mutex
	logger logger.LoggerInterface
}

func New(conf config.WorkerConfig) *Worker {
	if !conf.Enable {
		logger.Warnf(context.Background(), "worker %s is disabled", conf.Name)
		return nil
	}

	if conf.Script == "" || conf.Handler == "" {
		logger.Errorf(context.Background(), "empty script or handler, script: %+v, handler: %v", conf.Script, conf.Handler)
		return nil
	}

	worker := &Worker{
		ready:   false,
		running: false,
		stopSig: make(chan bool),
		logger: logger.WithFields(logger.Fields{
			"module": "worker",
			"name":   conf.Name,
		}),
	}

	if err := deep.Copy(&worker.conf, conf); err != nil {
		worker.logger.Errorf(context.Background(), "deep.Copy err : %+v", err.Error())
		return nil
	}

	if worker.conf.RateLimit <= 0 {
		worker.conf.RateLimit = defaultRateLimit
	}

	return worker
}

func (w *Worker) Name() string {
	return w.conf.Name
}

func (w *Worker) IsRunning() bool {
	return w.running
}

func (w *Worker) preload() error {
	if w.conf.Source != nil {
		w.source = source.NewSource(w.conf.Name, w.conf.Source.Type, w.conf.Source.Config)
	} else {
		logger.Warnf(context.Background(), "worker %s's source is empty", w.conf.Name)
	}

	if w.conf.Dest != nil {
		w.dest = dest.NewDest(w.conf.Name, w.conf.Dest.Type, w.conf.Dest.Config)
	} else {
		logger.Warnf(context.Background(), "worker %s's dest is empty", w.conf.Name)
	}

	if w.conf.Source != nil {
		w.inCh = make(chan worker_msg.Message, int(math.Max(defaultBufferLength, float64(w.conf.Source.BufferLength))))
	} else {
		w.inCh = make(chan worker_msg.Message, defaultBufferLength)
	}
	if w.conf.Dest != nil {
		w.outCh = make(chan worker_msg.Message, int(math.Max(defaultBufferLength, float64(w.conf.Dest.BufferLength))))
	}

	if w.conf.RateLimit == 0 {
		w.conf.RateLimit = defaultRateLimit
	}
	w.rl = rate.NewLimiter(rate.Every(time.Second*time.Duration(60/w.conf.RateLimit)), w.conf.RateLimit)

	if w.conf.Num == 0 {
		w.conf.Num = 1
	}
	ep := engine_pool.GetDefaultEnginePool()
	for i := 0; i < w.conf.Num; i++ {
		eng := ep.GetRawEngine()
		if w.conf.Config != nil && len(w.conf.Config) > 0 {
			eng.RegisterObject("config", w.conf.Config)
		}
		if err := eng.ParseFile(w.conf.Script); err != nil {
			w.logger.Errorf(context.Background(), fmt.Sprintf("engine.ParseFile %s error: %v", w.conf.Name, err))
			return err
		}
		eng.SetReady()
		w.engine = append(w.engine, eng)
	}
	w.running = true
	return nil
}

func (w *Worker) Start() error {
	if !w.conf.Enable {
		logger.Warnf(context.Background(), "worker %s is disabled", w.conf.Name)
		return nil
	}

	if err := w.preload(); err != nil {
		return err
	}

	sourceCtx, sourceCancel := context.WithCancel(context.Background())
	workerCtx, workerCancel := context.WithCancel(context.Background())
	destCtx, destCancel := context.WithCancel(context.Background())

	if w.source != nil {
		go w.source.Receive(sourceCtx, w.inCh)
	}

	for k := range w.engine {
		go w.waitMsg(k, workerCtx)
	}

	if w.dest != nil {
		go w.dest.Send(destCtx, w.outCh)
	}

	worker_map.Register(w.conf.Name, w.inCh)

	w.running = true
	<-w.stopSig

	if w.source != nil {
		sourceCancel()
		close(w.inCh)
	}
	time.Sleep(1 * time.Second)
	workerCancel()
	if w.dest != nil {
		time.Sleep(1 * time.Second)
		destCancel()
		close(w.outCh)
	}

	w.running = false
	return nil
}

func (w *Worker) waitMsg(k int, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-w.inCh:
			if len(msg.Content.Data) <= 0 {
				continue
			}
			if err := w.rl.Wait(ctx); err != nil {
				continue
			}

			w.process(ctx, k, msg)
		}
	}
}

func (w *Worker) process(ctx context.Context, k int, msg worker_msg.Message) {
	eng := w.engine[k]
	if eng == nil {
		logger.Errorf(context.Background(), "worker %s engine %d is nil", w.conf.Name, k)
		return
	}

	msg.Info.ProcessTime = time.Now().Unix()
	err, ret := eng.Call(w.conf.Handler, 1, msg.Content.Data, msg.Info)
	if err != nil {
		w.logger.Errorf(ctx, "worker %s engine %d handle error : %v", w.conf.Name, k, err)
		return
	}

	if w.dest == nil || ret == nil {
		return
	}

	res := luatool.ConvertLuaData(ret)
	if reflect.TypeOf(res) == reflect.TypeOf(map[string]interface{}{}) {
		w.outCh <- worker_msg.Message{
			Info:    msg.Info,
			Content: worker_msg.Content{Data: res.(map[string]interface{})},
		}
	} else if reflect.TypeOf(res) == reflect.TypeOf([]interface{}{}) {
		logger.Debugf(context.Background(), "worker stop %v", res)
	} else {
		logger.Errorf(context.Background(), "emit invalid type data = %v res = %v", reflect.TypeOf(res), reflect.TypeOf(res))
	}
}

func (w *Worker) Stop() {
	if w.running {
		w.stopSig <- true
	}
}

func (w *Worker) Running() bool {
	return w.running
}
