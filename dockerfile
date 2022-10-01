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

ARG MODE=dev

# 时区
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

COPY --from=builder /build/app /
COPY sql/init.sql /sql/

# nas-anime的web端口

ENV USERNAME=admin
ENV PASSWORD=admin
ENV MEDIA_PATH=/meido-anime
ENV SOURCE_PATH=/downloads

ENV SERVER_PORT=8081
ENV SERVER_TOKEN_EXPIRED_TIME=259200

ENV DB_PATH=./db/meido-anime.db
ENV DB_MAX_CONS=10

ENV QB_WEB_URL=http://localhost:9999
ENV QB_USERNAME=admin
ENV QB_PASSWORD=adminadmin
ENV QB_CATEGORY=meido-anime
ENV QB_DOWNLOAD_PATH=/downloads

ENTRYPOINT ["/app"]