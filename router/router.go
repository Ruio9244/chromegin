package router

import (
	"github.com/gin-gonic/gin"
)

// UseRouters 定义路由
func UseRouters(eng *gin.Engine) {
	hostPath := "api"

	bizGroup := eng.Group(hostPath)
	sysGroup := eng.Group(hostPath)

	// 系统级通用路由
	useSystemRouters(sysGroup)

	// 业务自定义路由
	useBizRouters(bizGroup)
}
