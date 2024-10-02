package httpserver

import (
	"github.com/gin-gonic/gin"
	"gotemplate/server/httpserver/middles"
)

func NewHandlerFuncFrom(method interface{}, opt ...middles.Option) gin.HandlerFunc {
	return middles.NewHandlerFuncFrom(method, opt...)
}
