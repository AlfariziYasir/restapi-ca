package validation

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"restapi/internal/web"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/xkeyideal/captcha/pool"
)

var captchaPool = pool.NewCaptchaPool(240, 80, 6, 2, 2, 2)
var cacheBuffer *pool.CaptchaBody

func CheckCaptchaSolver(valueSolution string, session sessions.Session) (bool, string) {
	v := session.Get("captcha")
	now := time.Now()

	captchaTime := session.Get("captcha_time")
	t := time.Unix(captchaTime.(int64), 0)
	if t.Add(time.Minute).Before(now) {
		return false, "captcha timeout"
	}

	flag := false
	str := fmt.Sprintf("%v", v)
	if v == nil {
		return flag, "captcha cannot empty"
	}

	if flag = str == valueSolution; !flag {
		return flag, "captcha not match"
	} else {
		return flag, "success"
	}
}

func CaptchaHandler(ctx *gin.Context) {
	var (
		result gin.H
	)
	base_url := GenerateCaptcha()
	result = gin.H{
		"image":        base_url,
		"captcha_code": string(cacheBuffer.Val),
	}
	session := sessions.Default(ctx)
	session.Set("captcha", string(cacheBuffer.Val))
	session.Set("captcha_time", time.Now().Unix())
	session.Save()

	web.MarshalPayload(ctx, http.StatusOK, "ok", result)
}

func GenerateCaptcha() string {
	cacheBuffer = captchaPool.GetImage()
	base_url := base64.StdEncoding.EncodeToString(cacheBuffer.Data.Bytes())
	base_url = "data:image/png;base64," + base_url
	return base_url
}
