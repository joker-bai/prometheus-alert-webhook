package utils

import (
	"code.coolops.cn/prometheus-alert-sms/alertMessage"
	"encoding/json"
	"log"
)

// 格式化数据
func FormatData(alert alertMessage.Alerts) string {
	// 获取状态，根据状态判断是故障还是恢复
	status := alert.Status
	switch status {
	case "resolved":
		log.Println("故障恢复消息")
		data := formatResolvedData(alert)
		return data
	case "firing":
		log.Println("故障告警消息")
		data := formatFiringData(alert)
		return data
	default:
		log.Println("无效的消息")
		return ""
	}
	return ""
}

func formatFiringData(alert alertMessage.Alerts) string {
	var newData alertMessage.FaultAlarm
	newData.AlertType = alert.Labels.AlertName
	newData.AlertDetails = alert.Annotations.Message + alert.Annotations.Description + alert.Annotations.Summary
	newData.AlertLevel = alert.Labels.Severity
	newData.AlertStatus = alert.Status
	newData.FaultTime = alert.StartsAt
	newData.InstanceInfo = alert.Labels.Instance
	newData.InstanceName = alert.Labels.Pod
	newData.Namespace = alert.Labels.Namespace
	mData, err := json.Marshal(newData)
	if err != nil {
		log.Println("序列化数据失败")
		return ""
	}
	return string(mData)
}

func formatResolvedData(alert alertMessage.Alerts) string {
	var newData alertMessage.FaultRecovery
	newData.AlertType = alert.Labels.AlertName
	newData.AlertDetails = alert.Annotations.Message + alert.Annotations.Description + alert.Annotations.Summary
	newData.AlertLevel = alert.Labels.Severity
	newData.AlertStatus = alert.Status
	newData.FaultTime = alert.StartsAt
	newData.InstanceInfo = alert.Labels.Instance
	newData.InstanceName = alert.Labels.Pod
	newData.Namespace = alert.Labels.Namespace
	newData.RecoveryTime = alert.EndsAt
	mData, err := json.Marshal(newData)
	if err != nil {
		log.Println("序列化数据失败")
		return ""
	}
	return string(mData)
}
