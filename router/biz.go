package router

import (
	"chromegin/service"
	"github.com/gin-gonic/gin"
)

func useBizRouters(rg *gin.RouterGroup) {
	cdp := rg.Group("chromedp")
	{
		cdp.GET("/python/ss", service.ChromedpShot)
		cdp.GET("/screenshot", service.ChromedpShot)
		cdp.GET("/screenshot-test", service.ChromedpShotTest)
	}
}
