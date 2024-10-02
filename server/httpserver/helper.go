package httpserver

import (
	"chatbot/server/httpserver/middles"
	"github.com/gin-gonic/gin"
)

func NewHandlerFuncFrom(method interface{}, opt ...middles.Option) gin.HandlerFunc {
	return middles.NewHandlerFuncFrom(method, opt...)
}
