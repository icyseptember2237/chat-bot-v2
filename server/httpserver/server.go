package httpserver

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"gotemplate/logger"
	"gotemplate/server/httpserver/middles"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Server interface {
	Name() (name string)
	Run(ctx context.Context) (err error)
	AddMiddles(ms ...middles.Middle)
	GetKernel() (kernel *gin.Engine)
	RegisterOnShutdown(f func())
}

type server struct {
	kernel     *gin.Engine
	options    Options
	onShutdown []func()
}

func NewServer(opts ...Option) Server {
	return NewServerWithOptions(newOptions(opts...))
}

func NewServerWithOptions(options Options) Server {
	kernel := gin.New()
	// default Use middles
	kernel.Use(
		middles.Recovery(),
		middles.FlowControlTag(),
		gin.Logger(),
	)

	kernel.Use(options.Middles...) // user set middles
	s := &server{kernel: kernel, options: options}
	return s
}

func (s *server) Name() string {
	return s.options.Name
}

func (s *server) Run(ctx context.Context) error {
	srv := &http.Server{Handler: s.kernel, Addr: s.options.Address}

	for _, f := range s.onShutdown {
		srv.RegisterOnShutdown(f)
	}

	// server listenAndServe
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("httpServer ListenAndServe err:%v", err)
			}
		}
	}()

	// handle signal
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	logger.Infof(ctx, "got signal %v, exit\n", <-ch)
	srv.Shutdown(context.Background())
	return nil
}

func (s *server) RegisterOnShutdown(f func()) {
	s.onShutdown = append(s.onShutdown, f)
}

func (s *server) AddMiddles(ms ...middles.Middle) {
	s.kernel.Use(ms...)
}

func (s *server) GetKernel() *gin.Engine {
	return s.kernel
}
