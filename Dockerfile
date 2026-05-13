# syntax=docker/dockerfile:1
#
# 使用预编译二进制文件构建 RSK API 镜像
# 本地编译: make build 或 scripts/build-local.sh
# 然后将 bin/dujiao-api 提交到 git，VPS 直接 COPY

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata \
    && mkdir -p /app/db /app/uploads /app/logs

COPY bin/dujiao-api /app/dujiao-api
RUN chmod +x /app/dujiao-api
COPY config.yml.example /app/config.yml.example

EXPOSE 8080

CMD ["./dujiao-api"]
