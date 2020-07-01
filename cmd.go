package main

import (
	"code.coolops.cn/prometheus-alert-sms/adapter"
	"code.coolops.cn/prometheus-alert-sms/conf"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

var Settings conf.Config

func init() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "/app/conf/sms.yaml"
	}
	file, err := ioutil.ReadFile(configPath)
	dir, _ := os.Getwd()
	fmt.Println(dir)
	if err != nil {
		log.Println("加载配置文件失败")
		panic(err)
	}
	if err = yaml.Unmarshal(file, &Settings); err != nil {
		log.Println("配置文件反序列化失败")
		panic(err)
	}
}

func RunCmd(ctx *gin.Context) {
	// 获取body数据
	data, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("获取消息数据失败")
		panic(err)
	}
	log.Println("接受到的报警数据",string(data))
	// 对数据进行序列号
	var sendData map[string]interface{}
	_ = json.Unmarshal(data, &sendData)
	log.Println("转换后的报警数据",sendData)
	getAdapter := Settings.Adapter.AdapterName

	switch getAdapter {
	case "RongLianYun":
		baseUrl := Settings.RongLianYun.BaseUrl
		accountSid := Settings.RongLianYun.AccountSid
		appToken := Settings.RongLianYun.AppToken
		appId := Settings.RongLianYun.AppId
		templateId := Settings.RongLianYun.TemplateId
		phones := Settings.RongLianYun.Phones
		rly := adapter.InitRongLianYun(baseUrl, accountSid, appToken, appId, templateId, phones)
		rly.Cmd(sendData)
	case "AliYun":
		aliRegion := Settings.AliYun.AliRegion
		accessKeyId := Settings.AliYun.AccessKeyId
		accessSecret := Settings.AliYun.AccessSecret
		signName := Settings.AliYun.SignName
		templateCode := Settings.AliYun.TemplateCode
		phoneNumbers := Settings.AliYun.PhoneNumbers
		aly := adapter.InitAliYun(aliRegion, accessKeyId, accessSecret, signName, phoneNumbers, templateCode)
		aly.Cmd(sendData)
		log.Println("阿里云短信")
	case "TengXunYun":
		log.Println("腾讯云短信")
	default:
		log.Println("没有找到对应的adapter")
	}

}
