package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
)

func ChromedpShotTest(c *gin.Context) {
	//var err error
	originUrl := "https://mojotv.cn/go/chromedp-example"
	/*	//url decode 参数
		originUrl, err = url.QueryUnescape(originUrl)
		if handleError(c, err) {
			return
		}
		if !strings.HasPrefix(originUrl, "http") {
			c.JSON(200, gin.H{"msg": originUrl + " 地址无效"})
			return
		}*/

	//md5 url 和时间信息一起拼接成截图名称.png
	fileName := fmt.Sprintf("%s.png", md5Encode(originUrl))
	//imagePath := path.Join(os.TempDir(), fileName)
	imagePath := path.Join("./data", fileName)
	if _, err := os.Stat(imagePath); os.IsExist(err) {
		//如果图片存在就直接gin response 图片
		c.File(imagePath)
		return
	}

	if err := runChromedp(originUrl, imagePath); err != nil {
		log.WithField("URL", originUrl).WithError(err)
		c.JSON(200, gin.H{"msg": err.Error()})
		return
	}
	c.File(imagePath)
}
