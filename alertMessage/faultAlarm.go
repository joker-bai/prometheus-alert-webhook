package alertMessage
// 故障报警发送的消息模板

type FaultAlarm struct {
	// 告警状态
	AlertStatus string
	// 告警类型
	AlertType string
	// 告警级别
	AlertLevel string
	// 告警详情
	AlertDetails string
	// 故障时间
	FaultTime string
	// 实例信息
	InstanceInfo string
	// 命名空间
	Namespace string
	// 实例名称
	InstanceName string
}