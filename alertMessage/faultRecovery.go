package alertMessage
// 故障恢复发送消息模板
type FaultRecovery struct {
	// 告警状态
	AlertStatus string
	// 告警类型
	AlertName string
	// 告警级别
	AlertSeverity string
	// 告警主题
	AlertSummary string
	// 告警详情
	AlertDetails string
	// 故障时间
	FaultTime string
	// 恢复时间
	RecoveryTime string
	// 实例信息
	Instance string
	// 命名空间
	Namespace string
	// 实例名称
	PodName string
	// 节点信息
	NodeName string
}
