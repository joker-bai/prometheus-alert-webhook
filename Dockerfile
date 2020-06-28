FROM alpine
RUN mkdir /app && \
    ln -sf /usr/share/zoneinfo/Asia/Shanghai  /etc/localtime
ADD prometheus-alert-sms /app
EXPOSE 9000
CMD ["/app/prometheus-alert-sms"]