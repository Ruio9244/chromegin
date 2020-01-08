package system

import (
	"net/http"
	"os"

	"git.code.oa.com/fip-team/fiutils/ip"

	"github.com/gin-gonic/gin"
)

// Healthz 健康状况显示
func Healthz(ctx *gin.Context) {
	host, _ := os.Hostname()
	internalIp, _ := ip.GetIntranetIP()
	ctx.JSON(http.StatusOK, &gin.H{
		"host": host,
		"ip":   internalIp,
	})
}
