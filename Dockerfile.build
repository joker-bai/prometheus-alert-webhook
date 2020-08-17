FROM alpine
RUN ln -snf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && mkdir /app
ENV TZ Asia/Shanghai
COPY conf/conf.yaml /app/conf/
ADD prometheus-alert-webhook /app
WORKDIR /app
EXPOSE 9000
ENTRYPOINT ["/app/prometheus-alert-webhook"]