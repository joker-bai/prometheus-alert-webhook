package adapter

import (
	"code.coolops.cn/prometheus-alert-sms/alertMessage"
	"code.coolops.cn/prometheus-alert-sms/utils"
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"log"
)

type aliyun struct {
	aliRegion string
	accessKeyId string
	accessSecret string
	signName string
	templateCode string
	phoneNumbers string
	sendData string
}

func InitAliYun(aliRegion,accessKeyId,accessSecret,signName,phoneNumbers,templateCode string)*aliyun{
	return &aliyun{
		aliRegion: aliRegion,
		accessKeyId: accessKeyId,
		accessSecret: accessSecret,
		signName: signName,
		phoneNumbers: phoneNumbers,
		templateCode: templateCode,
	}
}

func (a aliyun)Cmd(sendData alertMessage.AlertMessage) {
	client, err := dysmsapi.NewClientWithAccessKey(a.aliRegion,a.accessKeyId,a.accessSecret)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	// 手机号码
	request.PhoneNumbers = a.phoneNumbers
	// 签名
	request.SignName = a.signName
	// 模板ID
	request.TemplateCode = a.templateCode
	// 需要发送的数据
	for _, alert := range sendData.Alerts{
		a.sendData = utils.FormatData(alert)
		request.TemplateParam = a.sendData

		response, err := client.SendSms(request)
		if err != nil {
			fmt.Print(err.Error())
		}
		fmt.Printf("response is %#v\n", response)
	}
}

func (a aliyun)formatData(sendData alertMessage.AlertMessage)string{
	//alterType := sendData["告警类型"].(string)
	//alterHost := sendData["实例名称"].(string)
	//alterTime := sendData["故障时间"].(string)
	//alterDetails := sendData["告警详情"].(string)
	marshal, err := json.Marshal(sendData)
	if err != nil {
		log.Println("待发送数据转换失败")
		panic(err)
	}
	return string(marshal)
}