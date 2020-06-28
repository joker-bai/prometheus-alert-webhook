package conf

type Config struct {
	RongLianYun RongLianYun `yaml:"rongLianYun"`
	Adapter
	AliYun AliYun `yaml:"AliYun"`
}

type Adapter struct {
	AdapterName string `yaml:"adapter_name"`
}

type RongLianYun struct {
	BaseUrl string `yaml:"baseUrl"`
	AccountSid string `yaml:"accountSid"`
	AppToken string `yaml:"appToken"`
	AppId string `yaml:"appId"`
	TemplateId string `yaml:"templateId"`
	Phones []string `yaml:"phones"`
}

type AliYun struct{
	AliRegion string `yaml:"aliRegion"`
	AccessKeyId string `yaml:"accessKeyId"`
	AccessSecret string `yaml:"accessSecret"`
	PhoneNumbers string `yaml:"phoneNumbers"`
	SignName string `yaml:"signName"`
	TemplateCode string `yaml:"templateCode"`
}