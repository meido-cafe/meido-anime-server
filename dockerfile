FROM golang:alpine AS builder

RUN apk update && apk add --no-cache musl-dev gcc build-base

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn

WORKDIR /build
COPY . .

RUN go build -ldflags="-s -w" -trimpath -o app .


FROM alpine

# 时区
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

COPY --from=builder /build/app /
COPY sql/init.sql /sql/
COPY etc/common.yaml /etc/
COPY etc/dev.yaml /etc/
COPY etc/pro.yaml /etc/

# nas-anime的web端口
ENV QB_WEB_URL=http://localhost:8081
ENV QB_USERNAME=admin
ENV QB_PASSWORD=adminadmin
ENV QB_DOWNLOAD_PATH=/downloads
ENV QB_CATEGORY=meido-anime
ENV MEDIA_PATH=/meido-anime

ENTRYPOINT ["/app"]