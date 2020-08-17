package adapter

import (
	"bytes"
	"code.coolops.cn/prometheus-alert-sms/alertMessage"
	"code.coolops.cn/prometheus-alert-sms/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/robfig/go-cache"
	"io/ioutil"
	"log"
	"net/http"
	"time"
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
	var token string
	memCache:=cache.New(7200 * time.Second,7100*time.Second)
	cacheFromMem, ok := memCache.Get("wechatAccessToken")
	if ok{
		token = cacheFromMem.(string)
		log.Println("获取微信token，命中缓存")
	}else{
		log.Println("未命中缓存，从服务端获取access token")
		getToken, err := w.getToken()
		if err != nil {
			log.Println("get token from wechat failed.")
			panic(err)
		}
		token = getToken.AccessToken
		memCache.Set("wechatAccessToken",token,7200*time.Second)
	}
	fmt.Print(token)
	// 获取警报内容
	for _,data := range sendData.Alerts{
		w.sendData = utils.FormatData(data)
		content := formatData(w.sendData)
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
		w.sendMsg(token,dataBytes)
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