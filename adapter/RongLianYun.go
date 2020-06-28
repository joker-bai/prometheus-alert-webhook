package adapter

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	timestamp string
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

func (r rongLianYun) Cmd(sendData map[string]interface{}) {
	newData := r.formatData(sendData)
	// 获取时间戳
	r.timestamp = time.Now().Format("20060102150405")

	// 生成签名
	sign := r.md5Sign(r.timestamp)
	upperSign := strings.ToUpper(sign)
	// 构造请求的url
	requestUrl := r.baseUrl + "/2013-12-26/Accounts/" + r.accountSid + "/SMS/TemplateSMS?sig=" + upperSign
	for _, phone := range r.phones  {
		r.sendSMS(requestUrl, phone, newData)
	}
}

func (r rongLianYun)formatData(sendData map[string]interface{})[]string{
	// 通知类型，主机，故障，时间
	var formatData = make([]string, 0, 10)
	alterType := sendData["告警类型"].(string)
	alterHost := sendData["实例名称"].(string)
	alterTime := sendData["故障时间"].(string)
	alterDetails := sendData["告警详情"].(string)
	formatData = append(formatData, alterType, alterHost, alterDetails, alterTime)
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
	all, _ := ioutil.ReadAll(do.Body)
	fmt.Println(string(all))
}

func (r rongLianYun) md5Sign(nowTime string) string {
	m := md5.New()
	m.Write([]byte(r.accountSid + r.appToken + nowTime))
	return hex.EncodeToString(m.Sum(nil))
}
