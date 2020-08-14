FROM golang:1.12.6 AS build-env
ADD . /opt
ENV GOPROXY=https://goproxy.cn
WORKDIR /opt
RUN go build

FROM alpine
RUN apk add --no-cache tzdata \
    && ln -snf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && mkdir /app
ENV TZ Asia/Shanghai
COPY --from=build-env /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
COPY --from=build-env /opt/conf/conf.yaml /app/conf/
COPY --from=build-env /opt/prometheus-alert-sms /app
WORKDIR /app
EXPOSE 9000
ENTRYPOINT ["/app/prometheus-alert-sms"]
