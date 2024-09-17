package web

import (
	"RuijieSSLVPNSmsService/web/configs"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
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
	code := strings.Replace(content, "【SMSHOOK】", "", -1)

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
	credential := common.NewCredential(
		appConfig.Tencent.SecretId,
		appConfig.Tencent.SecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"
	client, _ := sms.NewClient(credential, appConfig.Tencent.Region, cpf)
	request := sms.NewSendSmsRequest()
	request.PhoneNumberSet = common.StringPtrs([]string{phone})
	request.SmsSdkAppId = common.StringPtr(appConfig.Tencent.SmsSdkAppId)
	request.TemplateId = common.StringPtr(appConfig.Tencent.TemplateId)
	request.SignName = common.StringPtr(appConfig.Tencent.SignName)
	// 验证码内容主体
	request.TemplateParamSet = common.StringPtrs([]string{code})
	response, err := client.SendSms(request)
	if err != nil {
		log.Println("[TencentSms]发送短信失败，原因是：", err.Error())
		c.String(http.StatusOK, "result=1&message=发送失败")
		return
	}
	requestId := response.Response.RequestId
	status := response.Response.SendStatusSet[0]
	if status.Code == common.StringPtr("Ok") {
		log.Println("[TencentSms]请求 ID：", requestId, " 发送短信成功，手机号：", phone, " 响应信息：", status)
		c.String(http.StatusOK, "result=0&message=发送成功&smsid="+*requestId)
		return
	} else {
		log.Println("[TencentSms]请求 ID：", requestId, " 发送短信失败，响应信息：", response.ToJsonString())
		c.String(http.StatusOK, "result=1&message=发送失败")
		return
	}
}
