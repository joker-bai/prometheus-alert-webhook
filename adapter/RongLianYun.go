package adapter

import (
	"bytes"
	"code.coolops.cn/prometheus-alert-sms/alertMessage"
	"code.coolops.cn/prometheus-alert-sms/utils"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type rongLianYun struct {
	baseUrl    string
	accountSid string
	appToken   string
	appId      string
	templateId string
	phones     []string
	timestamp  string
}

func InitRongLianYun(baseUrl, accountSid, appToken, appId, templateId string, phones []string) *rongLianYun {
	return &rongLianYun{
		baseUrl:    baseUrl,
		accountSid: accountSid,
		appToken:   appToken,
		appId:      appId,
		templateId: templateId,
		phones:     phones,
	}
}

func (r rongLianYun) Cmd(sendData alertMessage.AlertMessage) {
	//newData := r.formatData(sendData)
	// 获取时间戳
	r.timestamp = time.Now().Format("20060102150405")

	// 生成签名
	sign := r.md5Sign(r.timestamp)
	upperSign := strings.ToUpper(sign)
	// 构造请求的url
	requestUrl := r.baseUrl + "/2013-12-26/Accounts/" + r.accountSid + "/SMS/TemplateSMS?sig=" + upperSign
	// 把报警信息进行聚合去重，
	for _, phone := range r.phones {
		for _, alert := range sendData.Alerts {
			newData := utils.FormatData(alert)
			sendNewData := r.formatData(newData)
			r.sendSMS(requestUrl, phone, sendNewData)
		}

	}
}

func (r rongLianYun) formatData(sendData string) []string {
	// 通知类型，主机，故障，时间
	var formatData = make([]string, 0, 10)
	var newData map[string]string
	alertHost := ""
	err := json.Unmarshal([]byte(sendData), &newData)
	if err != nil {
		log.Println("反序列化需要发送的数据失败")
		return nil
	}
	alertType := newData["AlertType"]
	if value, ok := newData["Instance"]; ok && value != "" {
		alertHost = value
	} else {
		alertHost = newData["PodName"]
	}
	alertTime := newData["FaultTime"]
	alertDetails := newData["AlertDetails"]
	formatData = append(formatData, alertType, alertHost, alertDetails, alertTime)
	return formatData
}

func (r rongLianYun) sendSMS(requestUrl, phone string, newdata []string) {
	client := &http.Client{Timeout: time.Second}
	body := map[string]interface{}{
		"to":         phone,
		"appId":      r.appId,
		"templateId": r.templateId,
		"datas":      newdata,
	}
	marshal, _ := json.Marshal(body)
	buffer := bytes.NewBuffer(marshal)
	request, err := http.NewRequest(http.MethodPost, requestUrl, buffer)
	if err != nil {
		panic(err)
	}
	// 获取认证
	authString := r.accountSid + ":" + r.timestamp
	fmt.Println(authString)
	auth := base64.StdEncoding.EncodeToString([]byte(r.accountSid + ":" + r.timestamp))
	fmt.Println("auth", auth)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json;charset=utf-8")
	request.Header.Set("Authorization", auth)
	do, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	if do.StatusCode == http.StatusOK {
		log.Println("发送报警信息到" + phone + "成功！！！")
	} else {
		log.Println("发送报警信息到" + phone + "失败！！！")
		all, _ := ioutil.ReadAll(do.Body)
		fmt.Println(string(all))
	}
}

func (r rongLianYun) md5Sign(nowTime string) string {
	m := md5.New()
	m.Write([]byte(r.accountSid + r.appToken + nowTime))
	return hex.EncodeToString(m.Sum(nil))
}
