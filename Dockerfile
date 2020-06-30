FROM alpine
RUN mkdir /app && \
    ln -sf /usr/share/zoneinfo/Asia/Shanghai  /etc/localtime
COPY conf/sms.yaml /app/conf/
ADD prometheus-alert-sms /app
WORKDIR /app
EXPOSE 9000
CMD ["/app/prometheus-alert-sms"]