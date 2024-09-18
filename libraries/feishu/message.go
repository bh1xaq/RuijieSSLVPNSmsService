package feishu

import (
	"context"
	"errors"
	"github.com/larksuite/oapi-sdk-go/v3"
	larkcontact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
	"github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"log"
	"strings"
)

type Feishu struct {
	AppId               string `mapstructure:"app_id"`
	AppSecret           string `mapstructure:"app_secret"`
	TemplateId          string `mapstructure:"template_id"`
	TemplateVersionName string `mapstructure:"template_version_name"`
	Client              *lark.Client
}

func (feishu *Feishu) NewClient() {
	feishu.Client = lark.NewClient(feishu.AppId, feishu.AppSecret)
}

func (feishu *Feishu) Send(phone string, content string) (string, error) {
	// 验证码主体内容
	content = strings.Replace(content, "【SMSHOOK】", "", -1)
	contentMap := strings.Split(content, ",")
	feishuContent := `{"type":"template","data":{"template_id":"` + feishu.TemplateId + `","template_version_name":"` + feishu.TemplateVersionName + `","template_variable":{"v_func":"` + contentMap[0] + `","v_code":"` + contentMap[1] + `","v_name":"` + contentMap[2] + `"}}}`
	//log.Println("[Feishu]验证码内容：", feishuContent)
	userId, err := feishu.getUseridByPhone(phone)
	if err != nil {
		log.Println("[Feishu]获取用户ID失败，原因是：", err.Error())
		return "", errors.New("result=1&message=发送失败")
	}
	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(`user_id`).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(userId).
			MsgType(`interactive`).
			Content(feishuContent).
			Build()).
		Build()
	resp, err := feishu.Client.Im.Message.Create(context.Background(), req)
	if err != nil {
		log.Println("[Feishu]验证码内容解析失败，原因是：", err.Error())
		return "", errors.New("result=1&message=发送失败")
	}

	if !resp.Success() {
		log.Println("[Feishu]发送短信失败，原因是：", resp.Code, resp.Msg)
		return "", errors.New("result=1&message=发送失败")
	}
	log.Println("[Feishu]发送短信成功，手机号：", phone)
	return "result=0&message=发送成功", nil
}

func (feishu *Feishu) getUseridByPhone(phone string) (string, error) {
	req := larkcontact.NewBatchGetIdUserReqBuilder().
		UserIdType(`user_id`).
		Body(larkcontact.NewBatchGetIdUserReqBodyBuilder().
			Mobiles([]string{phone}).
			IncludeResigned(true).
			Build()).
		Build()
	resp, err := feishu.Client.Contact.User.BatchGetId(context.Background(), req)
	if err != nil {
		return "", err
	}
	if !resp.Success() {
		return "", errors.New("获取用户ID失败，接口返回：" + resp.Msg)
	}
	if resp.Data.UserList == nil || len(resp.Data.UserList) == 0 {
		return "", errors.New("用户不存在")
	}
	return *resp.Data.UserList[0].UserId, nil
}
