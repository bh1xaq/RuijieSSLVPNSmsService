package web

import (
	"RuijieSSLVPNSmsService/libraries/tencent_sms"
	"RuijieSSLVPNSmsService/web/configs"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

var appConfig configs.AppConfig

func Start() {
	if err := appConfig.Load(); err != nil {
		log.Fatalf("Load Config: %s\n", err)
	}
	engine := gin.Default()

	SetRouter(engine)

	srv := &http.Server{
		Addr:    appConfig.App.Listen,
		Handler: engine,
	}

	log.Printf("Start Service: http://%s\n", srv.Addr)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Service Listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

func SetRouter(engine *gin.Engine) {
	engine.GET("/sdkhttp/sendsms.aspx", HookStongNetSms)
}

func HookStongNetSms(c *gin.Context) {

	reg := c.Query("reg")
	pwd := c.Query("pwd")
	phone := c.Query("phone")
	content := c.Query("content")
	content = strings.Replace(content, "【SMSHOOK】", "", -1)

	// 拦截缺失的参数
	if reg == "" || pwd == "" || phone == "" || content == "" {
		log.Println("[StongNetSms]发送短信失败，原因是：参数缺失")
		c.String(http.StatusOK, "result=1&message=参数缺失")
		return
	}
	if reg != appConfig.Stong.Reg || pwd != appConfig.Stong.Pwd {
		log.Println("[StongNetSms]发送短信失败，原因是：鉴权失败")
		c.String(http.StatusOK, "result=1&message=鉴权失败")
		return
	}

	switch appConfig.App.Sender {
	case "tencent":
		// 调用腾讯云发送短信
		app := tencent_sms.Tencent{
			SecretId:    appConfig.Tencent.SecretId,
			SecretKey:   appConfig.Tencent.SecretKey,
			Region:      appConfig.Tencent.Region,
			SmsSdkAppId: appConfig.Tencent.SmsSdkAppId,
			TemplateId:  appConfig.Tencent.TemplateId,
			SignName:    appConfig.Tencent.SignName,
		}
		resp, err := app.Send(phone, content)
		if err != nil {
			c.String(http.StatusOK, err.Error())
			return
		} else {
			c.String(http.StatusOK, resp)
			return
		}
	}
}
