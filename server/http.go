package server

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gotemplate/handler"
	"gotemplate/server/httpserver"
	"net/http"
)

type HttpServer struct {
	Server httpserver.Server
}

var defaultHTTP = &HttpServer{}

// New creates a new HTTP server.
func (h *HttpServer) New(address string) Service {
	return &HttpServer{
		Server: newRouter(address),
	}
}

// 在此注册路由
var httphandlers = []handler.Handler{}

func newRouter(address string) httpserver.Server {
	httpServer := httpserver.NewServer(
		httpserver.WithName("template"),
		httpserver.WithAddress(address),
		httpserver.WithMiddles(),
	)

	router := httpServer.GetKernel()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/alive", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, nil)
	})

	v1 := router.Group("/v1")
	{
		for _, handler := range httphandlers {
			handler.Init(v1)
		}
	}

	return httpServer
}
