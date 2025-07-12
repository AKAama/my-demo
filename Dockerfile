FROM alpine:3.18

LABEL maintainer="yhma <yhma@sudytech.cn>"
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add -U tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY bin/"myapi" /app/
COPY etc /app/etc
WORKDIR /app
CMD ["./myapi","-c","etc/config.yaml"]