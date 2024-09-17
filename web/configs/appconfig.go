package configs

import (
	"github.com/spf13/viper"
	"log"
)

type AppConfig struct {
	App     App     `mapstructure:"app"`
	Stong   Stong   `mapstructure:"stong"`
	Tencent Tencent `mapstructure:"tencent"`
}

type App struct {
	Listen string `mapstructure:"listen"`
}

type Stong struct {
	Reg string `mapstructure:"reg"`
	Pwd string `mapstructure:"pwd"`
}

type Tencent struct {
	SecretId    string `mapstructure:"secret_id"`
	SecretKey   string `mapstructure:"secret_key"`
	Region      string `mapstructure:"region"`
	SmsSdkAppId string `mapstructure:"sms_sdk_app_id"`
	TemplateId  string `mapstructure:"template_id"`
	SignName    string `mapstructure:"sign_name"`
}

func (app *AppConfig) Load() error {
	// 设置配置文件名和类型
	viper.SetConfigName("config") // 不带扩展名
	viper.SetConfigType("toml")   // 配置文件类型
	viper.AddConfigPath(".")      // 配置文件路径，"." 表示当前目录

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	err := viper.Unmarshal(&app)
	if err != nil {
		return err
	}
	// 输出配置项
	log.Printf("GetAppConfig: %+v\n", app)
	return nil
}
