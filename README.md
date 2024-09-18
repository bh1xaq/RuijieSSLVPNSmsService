# Ruijie SSLVPN 短信服务扩展组件

目前 Ruijie SSLVPN 支持的短信方式有限，本组件实现了基于腾讯云的短信扩展。

## 原理
劫持华兴软通的短信接口，转发至腾讯云短信服务。

## 接口文档

GET http://www.stongnet.com/sdkhttp/sendsms.aspx

请求：

| 参数        | 描述                                     |
|-----------|----------------------------------------|
| reg       | 注册号（由华兴软通指定），不可为空                      |
| pwd       | 密码（由华兴软通指定），不可为空                       |
| sourceadd | 子通道号（最长10位，可为空）                        |
| phone     | 手机号码（最多1000个），多个用英文逗号(,)隔开，不可为空        |
| content   | 短信内容（UTF-8编码）（最多600个字符，字母、标点都算字符），不可为空 |

响应：

result=0&message=短信发送成功&smsid=

## 配置说明
您需要配置 SSLVPN 的配置和 RuijieSMS 配置一致。
### SSLVPN 侧
配置短信策略管理，短信提供商为华兴软通。

注册码同 config.toml 中的 stong.reg 参数。

密码同 config.toml 中的 stong.pwd 参数。

发送模板为固定值：${VCode}【SMSHOOK】

### RuijieSms 侧

1. [app] 节
   - listen: 应用程序监听的地址和端口号。
     - 默认值: 127.0.0.1:8080
   - sender：发送短信的渠道
     - 选项：feishu，tencent
2. [stong] 节
   - reg: 华兴软通平台的注册号，用于认证。必须替换为实际的注册号。
     - 示例："your_registration_value"
   - pwd: 华兴软通平台的密码。必须替换为实际的密码。
     - 示例："your_password_value"
3. [tencent] 节
   - secret_id: 腾讯云账户的密钥 ID，用于 API 鉴权。需替换为实际的 Secret ID。
     - 示例："your_secret_id"
   - secret_key: 腾讯云账户的密钥 Key，用于 API 鉴权。需替换为实际的 Secret Key。
     - 示例："your_secret_key"
   - region: 腾讯云服务器的区域，指定短信服务使用的区域代码。
     - 示例："your_region"（例如 "ap-guangzhou" 表示广州区域）
   - sms_sdk_app_id: 腾讯云短信服务的 SDK App ID，用于指定发送短信的应用 ID。需替换为实际的 SMS SDK App ID。
     - 示例："your_sms_sdk_app_id"
   - template_id: 短信模板 ID，用于发送短信时选择指定的模板。需替换为实际的模板 ID。
     - 示例："your_template_id"
   - sign_name: 短信签名，用于标识短信来源。需替换为实际的签名名称。
     - 示例："your_sign_name"
4. [feishu] 节
   - app_id：飞书应用的 APPID
     - 示例："your_app_id"
   - app_secret：飞书应用的 Secret
     - 示例："your_app_secret"
   - template_id: 飞书卡片搭建工具创建的卡片 ID
     - 示例："your_template_id"
   - template_version_name：飞书卡片搭建工具创建的卡片版本
     - 示例："your_template_version_name"

## 编译服务

```shell
./build.sh
```

## 运行服务

根据您所在的平台，运行对应的程序。

## 劫持 DNS
对锐捷网关配置 DNS，将 www.stongnet.com 指向您的服务。

## 锐捷 SSLVPN 短信服务 Debug

在锐捷网关上执行以下命令，查看短信服务状态。
```shell
show sslvpn short-message stong
```

## 特殊说明

本程序仅供学习使用，请您在下载后 24小时内删除，如果本程序侵犯了您或公司的权益，请及时联系我，我会在第一时间处理。