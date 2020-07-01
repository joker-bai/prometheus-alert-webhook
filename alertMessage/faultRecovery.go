package alertMessage
// 故障恢复发送消息模板
type FaultRecovery struct {
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
	// 恢复时间
	RecoveryTime string
	// 实例信息
	InstanceInfo string
	// 命名空间
	Namespace string
	// 实例名称
	InstanceName string
}
