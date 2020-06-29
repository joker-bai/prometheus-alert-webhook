FROM alpine
RUN mkdir /app && \
    ln -sf /usr/share/zoneinfo/Asia/Shanghai  /etc/localtime
COPY conf/sms.yaml /app/conf/
ADD prometheus-alert-sms /usr/local/bin/
EXPOSE 9000
CMD ["/usr/local/bin/prometheus-alert-sms"]