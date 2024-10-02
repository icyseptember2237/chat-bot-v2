package handler

import "github.com/gin-gonic/gin"

type Handler interface {
	Init(ginRouter *gin.RouterGroup)
}
