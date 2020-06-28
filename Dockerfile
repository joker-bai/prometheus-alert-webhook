FROM alpine
RUN mkdir /app && \
    ln -sf /usr/share/zoneinfo/Asia/Shanghai  /etc/localtime
ADD PrometheusAlertSMS /app
EXPOSE 9000
CMD ["/app/PrometheusAlertSMS"]