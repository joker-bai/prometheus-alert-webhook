package alertMessage

// 接受到的报警消息结构体
type AlertMessage struct {
	Receiver string `json:"receiver"`
	Status string `json:"status"`
	Alerts []Alerts
	GroupLabels GroupLabels
	CommonLabels CommonLabels
	CommonAnnotations map[string]interface{} `json:"commonAnnotations"`
	ExternalURL string `json:"externalURL"`
	Version string `json:"version"`
	GroupKey string `json:"groupKey"`
}

type Alerts struct {
	Status string `json:"status"`
	Labels Labels
	Annotations Annotations
	StartsAt string `json:"startsAt"`
	EndsAt string `json:"endsAt"`
	GeneratorURL string `json:"generatorURL"`

}

type Labels struct {
	AlertName string `json:"alertname"`
	Deployment string `json:"deployment"`
	Instance string `json:"instance"`
	Job string `json:"job"`
	Namespace string `json:"namespace"`
	Pod string `json:"pod"`
	Prometheus string `json:"prometheus"`
	Severity string `json:"severity"`
}

type Annotations struct {
	Message string `json:"message"`
	Description string `json:"description"`
	Summary string `json:"summary"`
	RunBookURL string `json:"runbook_url"`
}

type GroupLabels struct {
	Job string `json:"job"`
}

type CommonLabels struct {
	Instance string `json:"instance"`
	Job string `json:"job"`
	Namespace string `json:"namespace"`
	Prometheus string `json:"prometheus"`
	Severity string `json:"severity"`
}