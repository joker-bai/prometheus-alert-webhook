package adapter

import (
	"encoding/json"
	"log"
)

// 格式化需要发送的数据
func formatData(sendData string)(content string){
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
