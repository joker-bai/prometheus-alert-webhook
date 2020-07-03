# prometheus-alert-sms
----
用于prometheus使用alertmanager发送短信。
目前支持：
- 容联云
- 阿里云
## 部署
1、下载代码
```shell script
git clone https://github.com/cool-ops/prometheus-alert-sms.git
```
2、编译代码
```shell script
cd prometheus-alert-sms/
sh build.sh
```
3、打包镜像
```shell script
docker build -t registry.cn-hangzhou.aliyuncs.com/rookieops/prometheus-alert-sms:v0.0.7 .
```
注：镜像地址更换成自己的仓库地址  
4、推送镜像到镜像仓库
```shell script
docker push registry.cn-hangzhou.aliyuncs.com/rookieops/prometheus-alert-sms:v0.0.7
```
5、修改项目目录下的prometheus-alert-sms.yaml
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: sms-conf
  namespace: monitoring
data:
  sms.yaml: |
    adapter:
      adapter_name: "RongLianYun"
    RongLianYun:
      baseUrl : "https://app.cloopen.com:8883"
      accountSid : "xxxxxx"
      appToken   : "xxxxxx"
      appId      : "xxxxx"
      templateId : "xxx"
      phones : ["11111111111","22222222222"]

    AliYun:
      aliRegion: "cn-hangzhou"
      accessKeyId: "xxxx"
      accessSecret: "xxxx"
      phoneNumbers: "11111111111,22222222222"
      signName: "xxxx"
      templateCode: "xxxx"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: prometheus-alert-sms
  namespace: monitoring
spec:
  replicas: 1
  selector:
    matchLabels:
      app: prometheus-alert-sms
  template:
    metadata:
      labels:
        app: prometheus-alert-sms
    spec:
      containers:
        - name: prometheus-alert-sms
          image: registry.cn-hangzhou.aliyuncs.com/rookieops/prometheus-alert-sms:v0.0.7
          imagePullPolicy: IfNotPresent
          livenessProbe:
            failureThreshold: 3
            httpGet:
              path: /healthCheck
              port: tcp-9000
              scheme: HTTP
            initialDelaySeconds: 30
            periodSeconds: 10
            successThreshold: 1
            timeoutSeconds: 2
          readinessProbe:
              failureThreshold: 3
              httpGet:
                path: /healthCheck
                port: tcp-9000
                scheme: HTTP
              initialDelaySeconds: 30
              periodSeconds: 10
              successThreshold: 1
              timeoutSeconds: 2
          env:
            - name: CONFIG_PATH
              value: /app/conf/sms.yaml
          ports:
            - name: app-port
              containerPort: 9000
              protocol: TCP
          resources:
            limits:
              cpu: 500m
              memory: 1Gi
            requests:
              cpu: 500m
              memory: 1Gi
          volumeMounts:
            - name: sms-conf
              mountPath: /app/conf/sms.yaml
              subPath: sms.yaml
      volumes:
        - name: sms-conf
          configMap:
            name: sms-conf
---
apiVersion: v1
kind: Service
metadata:
  name: prometheus-alter-sms
  namespace: monitoring
spec:
  selector:
    app: prometheus-alert-sms
  ports:
    - name: app-port
      port: 9000
      targetPort: 9000
      protocol: TCP
```
到自己购买的短信服务获取对应的信息。
7、部署yaml文件
```shell script
kubectl apply -f prometheus-alert-sms.yaml
```
8、修改alertmanager的报警媒介
```shell script
 ......
      - receiver: sms 
        group_wait: 10s
        match:
          filesystem: node
    receivers:
    - name: 'sms'
      webhook_configs:
      - url: "http://prometheus-alter-sms.monitoring.svc:9000"
        send_resolved: true
......
```

9、模板示例
```shell script

{{ define "wechat.default.message" }}
{{- if gt (len .Alerts.Firing) 0 -}}
{{- range $index, $alert := .Alerts -}}
{{- if eq $index 0 }}
==========异常告警==========
告警类型: {{ $alert.Labels.alertname }}
告警级别: {{ $alert.Labels.severity }}
告警详情: {{ $alert.Annotations.message }}{{ $alert.Annotations.description}};{{$alert.Annotations.summary}}
故障时间: {{ ($alert.StartsAt.Add 28800e9).Format "2006-01-02 15:04:05" }}
{{- if gt (len $alert.Labels.instance) 0 }}
实例信息: {{ $alert.Labels.instance }}
{{- end }}
{{- if gt (len $alert.Labels.namespace) 0 }}
命名空间: {{ $alert.Labels.namespace }}
{{- end }}
{{- if gt (len $alert.Labels.node) 0 }}
节点信息: {{ $alert.Labels.node }}
{{- end }}
{{- if gt (len $alert.Labels.pod) 0 }}
实例名称: {{ $alert.Labels.pod }}
{{- end }}
============END============
{{- end }}
{{- end }}
{{- end }}
{{- if gt (len .Alerts.Resolved) 0 -}}
{{- range $index, $alert := .Alerts -}}
{{- if eq $index 0 }}
==========异常恢复==========
告警类型: {{ $alert.Labels.alertname }}
告警级别: {{ $alert.Labels.severity }}
告警详情: {{ $alert.Annotations.message }}{{ $alert.Annotations.description}};{{$alert.Annotations.summary}}
故障时间: {{ ($alert.StartsAt.Add 28800e9).Format "2006-01-02 15:04:05" }}
恢复时间: {{ ($alert.EndsAt.Add 28800e9).Format "2006-01-02 15:04:05" }}
{{- if gt (len $alert.Labels.instance) 0 }}
实例信息: {{ $alert.Labels.instance }}
{{- end }}
{{- if gt (len $alert.Labels.namespace) 0 }}
命名空间: {{ $alert.Labels.namespace }}
{{- end }}
{{- if gt (len $alert.Labels.node) 0 }}
节点信息: {{ $alert.Labels.node }}
{{- end }}
{{- if gt (len $alert.Labels.pod) 0 }}
实例名称: {{ $alert.Labels.pod }}
{{- end }}
============END============
{{- end }}
{{- end }}
{{- end }}
{{- end }}
```