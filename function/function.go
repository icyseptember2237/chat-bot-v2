package function

import (
	"bytes"
	"chatbot/function/config"
	"chatbot/logger"
	"chatbot/msg"
	"chatbot/server/httpserver"
	"chatbot/utils/crypto"
	"chatbot/utils/deep"
	engine_pool "chatbot/utils/engine_pool"
	"chatbot/utils/luatool"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/icyseptember2237/engine"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/time/rate"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	defaultRateLimit = 300
)

type Server struct {
	conf       config.Server
	handlerMap sync.Map
	stopSig    chan bool
	logger     logger.LoggerInterface
}

var f *Server

func New(conf config.Server) *Server {
	if conf.ServerPort == "" {
		return nil
	}

	f = &Server{
		stopSig: make(chan bool),
		logger: logger.WithFields(
			logger.Fields{
				"module": "function",
			}),
	}
	if err := deep.Copy(&f.conf, conf); err != nil {
		f.logger.Fatal(context.Background(), "init server config error")
	}
	return f
}

func GetFunctionServer() *Server {
	return f
}

func (f *Server) preLoad(ep *engine_pool.EnginePool, handlerCfg map[string]interface{}, script string) {
	if f.conf.PreloadCnt != 0 {
		eArr := make([]engine.Engine, 0)
		for i := 0; i < f.conf.PreloadCnt; i++ {
			eng, err, _ := ep.GetEngine(script)
			if err == nil {
				eng.RegisterObject("env", f.conf.Env)
				eng.RegisterObject("config", handlerCfg)

				eArr = append(eArr, eng)
			}
		}
		for _, eng := range eArr {
			ep.PutEngine(eng)
		}
	}
}

func (f *Server) initFunction(functions []config.Function) {
	for _, function := range functions {
		ep := engine_pool.NewEnginePool()
		for _, handler := range function.Handlers {
			if handler.RateLimit <= 0 {
				handler.RateLimit = defaultRateLimit
			}
			cfg := &config.Handler{
				Command:     handler.Command,
				Description: handler.Description,
				Script:      handler.Script,
				Handler:     handler.Handler,
				RateLimit:   handler.RateLimit,
			}
			if cfg.Script == "" {
				cfg.Script = function.Script
				cfg.SetEnginePool(ep)
			} else {
				handlerEp := engine_pool.NewEnginePool()
				cfg.SetEnginePool(handlerEp)
				f.preLoad(ep, function.Config, cfg.Script)
			}
			cfg.SetLimiter(rate.NewLimiter(rate.Every(time.Second*time.Duration(60/cfg.RateLimit)), handler.RateLimit))

			key := function.Name + "/" + handler.Command
			if !strings.HasPrefix(key, "/") {
				key = "/" + key
			}
			f.handlerMap.Store(key, cfg)
		}
		f.preLoad(ep, function.Config, function.Script)
	}
}

func (f *Server) Help(receive *msg.ReceiveMessage) {
	var help strings.Builder
	for _, function := range f.conf.Functions {
		help.WriteString(fmt.Sprintf("%s: %s\n", function.Name, function.Description))
		for _, handler := range function.Handlers {
			if handler.Command != "" {
				help.WriteString(fmt.Sprintf("%s:\n/%s -%s arg\n", handler.Description, function.Name, handler.Command))
			} else {
				help.WriteString(fmt.Sprintf("%s:\n/%s arg\n", handler.Description, function.Name))
			}
		}
		help.WriteString("\n")
	}
	msg.NewGroupMessage(
		receive.GroupId,
		msg.NewReplySegment(receive.MessageId),
		msg.NewTextSegment(help.String()),
	).Send(f.conf.BotAddr, f.conf.BotToken)
}

func (f *Server) Start() {
	port := ":8080"
	addr := "/msg"
	if f.conf.ServerPort != "" {
		port = f.conf.ServerPort
	}
	if f.conf.ServerAddr != "" {
		addr = f.conf.ServerAddr
	}

	server := httpserver.NewServer(httpserver.WithAddress(port))

	f.initFunction(f.conf.Functions)

	if f.conf.ServerToken == "" {
		server.GetKernel().POST(addr, f.receiveMessage)
	} else {
		group := server.GetKernel().Group("")
		group.Use(func(c *gin.Context) {
			sign := c.GetHeader("X-Signature")
			reqBody, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
			if sign == "" || "sha1="+crypto.HMACSha1(f.conf.ServerToken, string(reqBody)) != sign {
				fmt.Println(sign, "sha1="+crypto.HMACSha1(f.conf.ServerToken, string(reqBody)))
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			c.Next()
		})
		group.POST(addr, f.receiveMessage)
	}

	if f.conf.Static != nil {
		if f.conf.Static.Path == "/" {
			panic("static.path must not be /")
		}
		server.GetKernel().GET("/", func(c *gin.Context) {
			c.Redirect(http.StatusPermanentRedirect, f.conf.Static.Path)
		})
		server.GetKernel().Static(f.conf.Static.Path, f.conf.Static.Root)
	}

	ctx, stop := context.WithCancel(context.Background())
	go func() {
		<-f.stopSig
		stop()
	}()

	_ = server.Run(ctx)
}

func (f *Server) Reload(conf config.Server) {
	_ = deep.Copy(&f.conf, conf)

}

func (f *Server) Stop() {
	f.stopSig <- true
}

func (f *Server) receiveMessage(ctx *gin.Context) {
	var received msg.ReceiveMessage

	if err := ctx.BindJSON(&received); err != nil {
		logger.Errorf(ctx, "ctx.Bind err %v", err)
		ctx.Status(http.StatusNoContent)
		return
	}

	if !received.IsGroupMessage() {
		logger.Infof(ctx, "msg isn't group message")
		ctx.Status(http.StatusNoContent)
		return
	}

	if !received.CheckSource(f.conf.WhiteGroup, f.conf.BanGroup) {
		logger.Infof(ctx, "msg source %v is not on whitelist or baned", received.GroupId)
		ctx.Status(http.StatusNoContent)
		return
	}

	if !received.CheckFormat() {
		logger.Infof(ctx, "msg format wrong %+v", received.Message)
		ctx.Status(http.StatusNoContent)
		return
	}

	entry, command, ok := received.SplitMessage()
	if !ok {
		logger.Infof(ctx, "msg split err %+v, result: %s %s %s", received.Message, entry, command, received.Text)
		ctx.Status(http.StatusNoContent)
		return
	}

	key := entry + "/" + command
	asteriskKey := entry + "/*"
	if strings.HasPrefix(key, "/help/") {
		logger.Infof(ctx, "help msg")
		f.Help(&received)
		ctx.Status(http.StatusNoContent)
		return
	}
	var handler interface{}
	if handler, ok = f.handlerMap.Load(key); !ok {
		if handler, ok = f.handlerMap.Load(asteriskKey); !ok {
			logger.Infof(ctx, "function %s not found", key)
			ctx.Status(http.StatusNoContent)
			return
		}
	}

	if handler.(*config.Handler).ReachLimit() {
		logger.Errorf(ctx, "request ratelimit reached: %+v", key)
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}

	var (
		err  error
		rets []interface{}
		eng  engine.Engine
	)

	eng, err, _ = handler.(*config.Handler).GetEnginePool().GetEngine(handler.(*config.Handler).Script)
	defer handler.(*config.Handler).GetEnginePool().PutEngine(eng)

	if err != nil {
		f.logger.Errorf(ctx, "engine_pool.GetEngine error %+v", err)
		ctx.Status(http.StatusNoContent)
		return
	}

	req := make(map[string]interface{})
	jstring, err := json.Marshal(&received)
	if err != nil {
		f.logger.Errorf(ctx, "json.Marshal error %+v", err)
		ctx.Status(http.StatusNoContent)
		return
	}

	err = json.Unmarshal(jstring, &req)
	if err != nil {
		f.logger.Errorf(ctx, "json.Unmarshal error %+v", err)
		ctx.Status(http.StatusNoContent)
		return
	}

	err, rets = eng.Call(handler.(*config.Handler).Handler, 2, req)
	if err != nil {
		f.logger.Errorf(ctx, "handleMessage error: %+v", err)
		_ = ctx.AbortWithError(http.StatusNoContent, err)
		return
	}

	var (
		retRes interface{}
		retErr interface{}
	)

	if len(rets) == 1 {
		retRes = luatool.ConvertLuaData(rets[0])
	} else if len(rets) == 2 {
		retRes = luatool.ConvertLuaData(rets[0])
		retErr = luatool.ConvertLuaData(rets[1])
	}

	if retRes == nil {
		logger.Infof(ctx, "no rets")
		ctx.Status(http.StatusNoContent)
		return
	}

	if retErr != nil {
		logger.Errorf(ctx, "lua script error: %+v", retErr)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	var segments []*msg.Segment

	if err := mapstructure.Decode(retRes, &segments); err != nil {
		logger.Errorf(ctx, "mapstructure.Decode error: %+v", err)
		ctx.Status(http.StatusNoContent)
		return
	}

	for k, m := range segments {
		n := make(map[string]interface{})
		for k1, v := range m.Data {
			n[strings.ToLower(k1)] = v
		}
		segments[k].Data = n
	}

	send := &msg.GroupMessage{
		GroupId:    received.GroupId,
		Message:    segments,
		AutoEscape: false,
	}

	if _, err := send.Send(f.conf.BotAddr, f.conf.BotToken); err != nil {
		logger.Errorf(ctx, "sendMessage error: %+v", retErr)
		ctx.Status(http.StatusNoContent)
		return
	}

	ctx.Status(http.StatusNoContent)
	return
}
