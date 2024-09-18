package tencent_sms

import (
	"encoding/json"
	"errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"log"
)

type Tencent struct {
	SecretId    string `mapstructure:"secret_id"`
	SecretKey   string `mapstructure:"secret_key"`
	Region      string `mapstructure:"region"`
	SmsSdkAppId string `mapstructure:"sms_sdk_app_id"`
	TemplateId  string `mapstructure:"template_id"`
	SignName    string `mapstructure:"sign_name"`
}

func (tencent Tencent) Send(phone string, content string) (string, error) {
	credential := common.NewCredential(
		tencent.SecretId,
		tencent.SecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"
	client, _ := sms.NewClient(credential, tencent.Region, cpf)
	request := sms.NewSendSmsRequest()
	request.PhoneNumberSet = common.StringPtrs([]string{phone})
	request.SmsSdkAppId = common.StringPtr(tencent.SmsSdkAppId)
	request.TemplateId = common.StringPtr(tencent.TemplateId)
	request.SignName = common.StringPtr(tencent.SignName)
	// 验证码内容主体
	var jsonContent map[string]interface{}
	err := json.Unmarshal([]byte(content), &jsonContent)
	if err != nil {
		log.Println("[TencentSms]验证码内容解析失败，原因是：", err.Error())
		return "", errors.New("result=1&message=发送失败")
	}
	request.TemplateParamSet = common.StringPtrs([]string{jsonContent["name"].(string), jsonContent["func"].(string), jsonContent["code"].(string)})
	response, err := client.SendSms(request)
	if err != nil {
		log.Println("[TencentSms]发送短信失败，原因是：", err.Error())
		return "", errors.New("result=1&message=发送失败")
	}
	requestId := response.Response.RequestId
	status := response.Response.SendStatusSet[0]
	if status.Code == common.StringPtr("Ok") {
		log.Println("[TencentSms]请求 ID：", requestId, " 发送短信成功，手机号：", phone, " 响应信息：", status)
		return "result=0&message=发送成功&smsid=" + *requestId, nil
	} else {
		log.Println("[TencentSms]请求 ID：", requestId, " 发送短信失败，响应信息：", response.ToJsonString())
		return "", errors.New("result=1&message=发送失败")
	}
}
