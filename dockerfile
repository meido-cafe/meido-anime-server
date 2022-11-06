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

# 用户名密码
ENV USERNAME=admin
ENV PASSWORD=admin
# 服务端口
ENV SERVER_PORT=8081
# token过期时间
ENV SERVER_TOKEN_EXPIRED_TIME=259200

# 数据库文件存储路径
ENV DB_PATH=./db/meido-anime.db
# 数据库最大连接数
ENV DB_MAX_CONS=10
# QB的web url
ENV QB_WEB_URL=http://localhost:9999
# QB的用户名
ENV QB_USERNAME=admin
# QB的密码
ENV QB_PASSWORD=adminadmin
# 番剧在QB中的分类
ENV QB_CATEGORY=meido-anime
# 指定QB下载番剧路径
ENV QB_DOWNLOAD_PATH=/downloads
# QB下载后的番剧路径
ENV SOURCE_PATH=/downloads
# 媒体库的番剧路径 (硬链接路径)
ENV MEDIA_PATH=/meido-anime

ENTRYPOINT ["/app"]