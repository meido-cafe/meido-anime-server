.PHONY: wire
wire:
	@cd factory && wire

.PHONY: local
local:
	@\
	USERNAME=admin \
	PASSWORD=admin \
	SERVER_PORT=8081 \
    SERVER_TOKEN_EXPIRED_TIME=259200 \
    MEDIA_PATH="D:\\home\\meido-cafe\\meido-anime-server\\docker\\jellyfin\\media" \
    SOURCE_PATH="D:\\home\\meido-cafe\\meido-anime-server\\docker\\qbittorrent\\downloads" \
    DB_PATH=./meido-anime.db \
    DB_MAX_CONS=10 \
    QB_USERNAME=admin \
    QB_PASSWORD=adminadmin \
    QB_WEB_URL=http://localhost:9999 \
    QB_CATEGORY=meido-anime \
    QB_DOWNLOAD_PATH=/downloads \
	go run main.go


.PHONY: dev
dev:
	@\
	USERNAME=admin \
	PASSWORD=admin \
	SERVER_PORT=8081 \
    SERVER_TOKEN_EXPIRED_TIME=259200 \
    MEDIA_PATH=D:\home\meido-cafe\meido-anime-server\docker\jellyfin\media \
    SOURCE_PATH=D:\home\meido-cafe\meido-anime-server\docker\qbittorrent\downloads \
    DB_PATH=./meido-anime.db \
    DB_MAX_CONS=10 \
    QB_USERNAME=admin \
    QB_PASSWORD=adminadmin \
    QB_WEB_URL=http://localhost:9999 \
    QB_CATEGORY=meido-anime \
    QB_DOWNLOAD_PATH=/downloads \
	go run main.go