package repo

import (
	"context"
)

type TorrentClientInterface interface {
	DeleteTorrent(ctx context.Context, hashes []string, deleteFile bool) (err error) // 删除种子

	SetRule(ctx context.Context, data map[string]string) (err error) // 设置规则
	DeleteRule(ctx context.Context, ruleName string) (err error)     // 删除下载规则

	AddRss(ctx context.Context, rss string, name string) (err error) // rss订阅
	DeleteRss(ctx context.Context, rss string) (err error)           // 删除rss

	GetLog(ctx context.Context) (res []map[string]any, err error) // 获取日志
}
