package system

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

// Host 获取当前主机名
func Host(ctx *gin.Context) {
	host, _ := os.Hostname()
	ctx.String(http.StatusOK, host)
}
