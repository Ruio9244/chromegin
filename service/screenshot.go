package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"math"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func ChromedpShot(c *gin.Context) {
	var err error
	originUrl := c.Query("url")
	//url decode 参数
	originUrl, err = url.QueryUnescape(originUrl)
	if handleError(c, err) {
		return
	}
	if !strings.HasPrefix(originUrl, "http") {
		c.JSON(200, gin.H{"msg": originUrl + " 地址无效"})
		return
	}

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

func runChromedp(targetUrl, imagePath string) error {
	// create context
	// timeout 90 秒
	timeContext, cancelFunc := context.WithTimeout(context.Background(), time.Second*90)
	defer cancelFunc()

	ctx, cancel := chromedp.NewContext(timeContext)
	defer cancel()

	// capture screenshot of an element
	var buf []byte
	// capture entire browser viewport, returning png with quality=50
	if err := chromedp.Run(ctx, fullScreenshot(targetUrl, 90, &buf)); err != nil {
		return err
	}
	imagePath = filepath.FromSlash(filepath.Join(RootDir(), imagePath))
	return ioutil.WriteFile(imagePath, buf, 0644)
}

func fullScreenshot(urlstr string, quality int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Emulate(device.IPad),
		chromedp.EmulateViewport(1024, 2048, chromedp.EmulateScale(1)),
		chromedp.Navigate(urlstr),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// get layout metrics
			_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}

			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

			// force viewport emulation
			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			if err != nil {
				return err
			}
			// capture screenshot
			*res, err = page.CaptureScreenshot().
				WithQuality(quality).
				WithClip(&page.Viewport{
					X:      contentSize.X,
					Y:      contentSize.Y,
					Width:  contentSize.Width,
					Height: contentSize.Height,
					Scale:  1,
				}).Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	}
}

func handleError(c *gin.Context, err error) bool {
	if err != nil {
		//logrus.WithError(err).Error("gin context http handler error")
		c.JSON(200, gin.H{"msg": err.Error()})
		return true
	}
	return false
}

func md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func RootDir() string {
	wd, _ := os.Getwd()
	return wd
}
