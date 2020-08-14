package main

import (
	"code.coolops.cn/prometheus-alert-sms/adapter"
	"code.coolops.cn/prometheus-alert-sms/alertMessage"
	"code.coolops.cn/prometheus-alert-sms/conf"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
)

var Settings conf.Config

//func init() {
//	configPath := os.Getenv("CONFIG_PATH")
//	if configPath == "" {
//		configPath = "/app/conf/conf.yaml"
//	}
//	file, err := ioutil.ReadFile(configPath)
//	dir, _ := os.Getwd()
//	fmt.Println(dir)
//	if err != nil {
//		log.Println("加载配置文件失败")
//		panic(err)
//	}
//	if err = yaml.Unmarshal(file, &Settings); err != nil {
//		log.Println("配置文件反序列化失败")
//		panic(err)
//	}
//}

func init() {
	viper.SetConfigName("conf")
	viper.AddConfigPath("conf")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("config file error: %s\n", err)
		return
	}
}
	//get := viper.Get("adapter")
//}

func RunCmd(ctx *gin.Context) {
	// 获取body数据
	data, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("获取消息数据失败")
		panic(err)
	}
	log.Println("接受到的报警数据",string(data))
	// 对数据进行序列号
	//var sendData map[string]interface{}
	var sendData alertMessage.AlertMessage
	_ = json.Unmarshal(data, &sendData)
	log.Println("转换后的报警数据",sendData)
	// 对数据进行格式化
	// getAdapter := Settings.Adapter.AdapterName

	// 从配置文件读取webhook配置
	adapters := viper.GetStringSlice("adapter")
	for _,myAdapter := range adapters{
		// 判断adapter的开关是否打开
		isEnabled := viper.GetBool(myAdapter+".enable")
		// 如果是打开的，则读取起配置，并发送消息
		if isEnabled{
			switch myAdapter {
			case "sms":
				// 判断是哪个短信平台
				smsAdapter := viper.GetString("sms.adapter_name")
				switch smsAdapter {
				case "RongLianYun":
					baseUrl := viper.GetString("sms.RongLianYun.baseUrl")
					accountSid := viper.GetString("sms.RongLianYun.accountSid")
					appToken := viper.GetString("sms.RongLianYun.appToken")
					appId := viper.GetString("sms.RongLianYun.appId")
					templateId := viper.GetString("sms.RongLianYun.templateId")
					phones := viper.GetStringSlice("sms.RongLianYun.phones")
					rly := adapter.InitRongLianYun(baseUrl, accountSid, appToken, appId, templateId, phones)
					rly.Cmd(sendData)
				case "AliYun":
					aliRegion := viper.GetString("sms.AliYun.aliRegion")
					accessKeyId := viper.GetString("sms.AliYun.accessKeyId")
					accessSecret := viper.GetString("sms.AliYun.accessSecret")
					signName := viper.GetString("sms.AliYun.signName")
					templateCode := viper.GetString("sms.AliYun.templateCode")
					phoneNumbers := viper.GetString("sms.AliYun.phoneNumbers")
					aly := adapter.InitAliYun(aliRegion, accessKeyId, accessSecret, signName, phoneNumbers, templateCode)
					aly.Cmd(sendData)
					log.Println("阿里云短信")
				case "TengXunYun":
					log.Println("腾讯云短信")
				default:
					log.Println("没有找到对应的adapter")
				}
			case "wechat":
				toUser := viper.GetString("wechat.toUser")
				agentId := viper.GetString("wechat.agentId")
				corpId := viper.GetString("wechat.corpid")
				corpSecret := viper.GetString("wechat.corpSecret")
				wc := adapter.InitWeChat(toUser,agentId,corpId,corpSecret)
				wc.Cmd(sendData)
			case "dingTalk":
				secret := viper.GetString("dingTalk.secret")
				accessToken := viper.GetString("dingTalk.access_token")
				dt := adapter.InitDingTalk(secret,accessToken)
				dt.Cmd(sendData)
			default:
				log.Println("请指定至少一个adapter")
			}
		}
	}
}

