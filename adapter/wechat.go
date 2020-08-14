package adapter

import (
	"bytes"
	"code.coolops.cn/prometheus-alert-sms/alertMessage"
	"code.coolops.cn/prometheus-alert-sms/utils"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	// 发送消息的url
	sendUrl = `https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=`
	// 获取token的url
	getTokenUrl = `https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=`
)

// 定义失败状态码
var requestError = errors.New("request error,please check it")

type accessToken struct {
	ErrorCode int `json:"errcode"`
	ErrorMsg string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIN int `json:"expires_in"`
}

// 定义消息文本格式
type wechatMsg struct {
	ToUser  string            `json:"touser"`
	ToParty string            `json:"toparty"`
	ToTag   string            `json:"totag"`
	MsgType string            `json:"msgtype"`
	AgentId string               `json:"agentid"`
	Text    map[string]string `json:"text"`
	Safe    int               `json:"safe"`
}

// 错误消息
type sendMsgError struct {
	ErrCode int    `json:"errcode`
	ErrMsg  string `json:"errmsg"`
}

type wechat struct {
	toUser string
	agentId string
	corpId string
	corpSecret string
	sendData string
}

func InitWeChat(toUser,agentId,corpId,corpSecret string) *wechat{
	return &wechat{
		toUser:     toUser,
		agentId:    agentId,
		corpId:     corpId,
		corpSecret: corpSecret,
	}
}

func (w wechat)Cmd(sendData alertMessage.AlertMessage){
	// 获取token
	token, err := w.getToken()
	if err != nil {
		log.Println("get token from wechat failed.")
		panic(err)
	}
	// 获取警报内容
	for _,data := range sendData.Alerts{
		w.sendData = utils.FormatData(data)
		content := w.formatData(w.sendData)
		var msg = wechatMsg{
			ToUser:  w.toUser,
			MsgType: "text",
			AgentId: w.agentId,
			Text:   map[string]string{"content": content},
		}
		dataBytes, err := json.Marshal(msg)
		if err != nil {
			log.Println("发送数据序列化失败")
			panic(err)
		}
		w.sendMsg(token.AccessToken,dataBytes)
	}
}

func (w wechat)getToken()(at accessToken,err error){
	url := getTokenUrl+w.corpId+"&corpsecret="+w.corpSecret
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = requestError
		return
	}
	buf, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(buf, &at)
	if err != nil {
		err = errors.New("corpid or corpsecret error")
		return
	}
	return
}

// 格式化需要发送的数据
func (w wechat)formatData(sendData string)(content string){
	var newData map[string]string
	err := json.Unmarshal([]byte(sendData), &newData)
	if err != nil {
		log.Println("反序列化需要发送的数据失败")
		panic(err)
	}
	content += "==========异常告警==========" + "\n"
	if value,ok := newData["AlertName"]; ok && value != "" {
		content += "告警类型：" + value + "\n"
	}
	if value,ok := newData["AlertStatus"]; ok && value != "" {
		content += "告警状态：" + value + "\n"
	}
	if value,ok := newData["AlertSeverity"];ok && value != ""{
		content += "告警级别：" + value + "\n"
	}
	if value,ok :=	newData["AlertSummary"];ok && value != ""{
		content += "告警主题：" + value + "\n"
	}
	if value,ok := newData["AlertDetails"];ok && value != ""{
		content += "告警详情：" + value + "\n"
	}
	if value,ok := newData["Instance"];ok && value != ""{
		content += "实例信息：" + value + "\n"
	}
	if value,ok := newData["Namespace"];ok && value != "" {
		content += "命名空间：" + value + "\n"
	}
	if value,ok := newData["PodName"];ok && value != "" {
		content += "实例名称：" + value + "\n"
	}
	if value,ok:=  newData["NodeName"];ok && value != ""{
		content += "节点信息：" + value + "\n"
	}
	if value,ok := newData["FaultTime"];ok && value != ""{
		content += "故障时间：" + value + "\n"
	}
	if value,ok := newData["RecoveryTime"];ok && value != ""{
		content += "恢复时间：" + value + "\n"
	}
	content += "============END============"
	return
}

// 发送消息
func (w wechat) sendMsg(accessToken string,msgBody []byte){
	buffer := bytes.NewBuffer(msgBody)
	url := sendUrl+accessToken
	post, err := http.Post(url, "application/json", buffer)
	if err != nil {
		panic(err)
	}
	defer post.Body.Close()
	buf, _ := ioutil.ReadAll(post.Body)
	var e sendMsgError
	err = json.Unmarshal(buf, &e)
	if err != nil {
		panic(err)
	}
	log.Println("告警信息发送到企业微信成功")
}