package job

import (
	globalconfig "chatbot/config"
	"chatbot/job/config"
	"chatbot/logger"
	"chatbot/msg"
	"chatbot/utils/deep"
	"chatbot/utils/engine_pool"
	"chatbot/utils/luatool"
	"context"
	"errors"
	"fmt"
	"github.com/icyseptember2237/engine"
	"github.com/mitchellh/mapstructure"
	"github.com/robfig/cron"
	"strings"
	"sync"
	"time"
)

const MaxLogLength = 10

type Log struct {
	startTime int64
	endTime   int64
	err       error
}

type Job struct {
	conf config.JobConfig

	ready   bool
	running bool
	mu      sync.Mutex

	logs []Log

	stopSig chan bool
	corn    *cron.Cron
	engine  engine.Engine
	logger  logger.LoggerInterface
}

func NewJob(conf config.JobConfig) *Job {
	ctx := context.Background()
	if !conf.Enable {
		logger.Warnf(ctx, "job %s is disabled", conf.Name)
		return nil
	}

	if conf.Script == "" || conf.Handler == "" {
		logger.Warnf(ctx, "job %s script or handler is empty, script: %s, handler: %s", conf.Name, conf.Script, conf.Handler)
		return nil
	}

	if conf.Cron == "" {
		logger.Warnf(ctx, "job %s cron is empty", conf.Name)
		return nil
	}

	job := &Job{
		logs:    make([]Log, 0),
		stopSig: make(chan bool),
		ready:   false,
		running: false,
		logger: logger.WithFields(logger.Fields{
			"module": "job",
			"name":   conf.Name,
		}),
	}

	if err := deep.Copy(&job.conf, conf); err != nil {
		job.logger.Warnf(ctx, "job %s config copy error: %v", conf.Name, err)
	}

	return job
}

func (j *Job) preload() error {
	j.engine = engine_pool.NewEnginePool().GetRawEngine()
	if j.conf.Config != nil && len(j.conf.Config) > 0 {
		j.engine.RegisterObject("config", j.conf.Config)
	}

	if err := j.engine.ParseFile(j.conf.Script); err != nil {
		return err
	}

	j.ready = true
	return nil
}

func (j *Job) run(ts int64) {
	var err error
	defer func() {
		log := Log{
			startTime: ts,
			endTime:   time.Now().Unix(),
			err:       err,
		}
		j.logs = append(j.logs, log)
		if len(j.logs) > MaxLogLength {
			j.logs = j.logs[1:]
		}
	}()

	if !j.mu.TryLock() {
		err = errors.New(fmt.Sprintf("job %s mutex locked", j.conf.Name))
		j.logger.Errorf(context.Background(), "job %s run err: %v", j.conf.Name, err)
		return
	}
	defer j.mu.Unlock()
	var rets []interface{}
	if err, rets = j.engine.Call(j.conf.Handler, 1, ts); err != nil {
		logger.Errorf(context.Background(), "job %s run err: %v", j.conf.Name, err)
		return
	}

	retRes := luatool.ConvertLuaData(rets[0])
	if retRes == nil {
		return
	}

	if j.conf.WhiteGroup == nil || len(j.conf.WhiteGroup) <= 0 {
		err = errors.New(fmt.Sprintf("job %s white group is empty", j.conf.Name))
		j.logger.Errorf(context.Background(), "job %s run err: %v", j.conf.Name, err)
		return
	}

	var segments []*msg.Segment

	if err = mapstructure.Decode(retRes, &segments); err != nil {
		err = errors.New(fmt.Sprintf("job %s decode err: %v", j.conf.Name, err))
		return
	}

	for k, m := range segments {
		n := make(map[string]interface{})
		for k1, v := range m.Data {
			n[strings.ToLower(k1)] = v
		}
		segments[k].Data = n
	}

	for _, groupId := range j.conf.WhiteGroup {
		send := &msg.GroupMessage{
			GroupId:    groupId,
			Message:    segments,
			AutoEscape: false,
		}

		if _, err = send.Send(globalconfig.Get().Server.BotAddr, globalconfig.Get().Server.BotToken); err != nil {
			logger.Errorf(context.Background(), "job %s send err: %v", j.conf.Name, err)
		}
	}
}

func (j *Job) Start() error {
	if err := j.preload(); err != nil {
		logger.Errorf(context.Background(), "job %s preload error: %v", j.conf.Name, err)
		return err
	}

	j.corn = cron.New()
	err := j.corn.AddFunc(j.conf.Cron, func() {
		ts := time.Now().Unix()
		j.run(ts)
	})
	if err != nil {
		logger.Errorf(context.Background(), "start job %s error: %v", j.conf.Name, err)
		return err
	}
	j.corn.Start()
	j.running = true

	<-j.stopSig
	j.mu.Lock()
	defer j.mu.Unlock()
	j.engine.Close()
	j.corn.Stop()
	j.running = false

	return nil
}

func (j *Job) Stop() {
	if j.running {
		j.stopSig <- true
	}
}

func (j *Job) Execute() {
	ts := time.Now().Unix()
	j.run(ts)
}

func (j *Job) Running() bool {
	return j.running
}
