package router

import (
	"chromegin/api/system"
	"github.com/gin-gonic/gin"
)

// 系统级通用路由
func useSystemRouters(rg *gin.RouterGroup) {
	rg.GET("/host", system.Host)
}
