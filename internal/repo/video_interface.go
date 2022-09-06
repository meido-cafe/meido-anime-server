package repo

import (
	"context"
	"meido-anime-server/internal/model"
	"meido-anime-server/internal/model/vo"
)

type VideoInterface interface {
	/* 番剧 */
	InsertOne(ctx context.Context, video model.Video, rule model.QBSetRule) (err error)                         // 插入一条番剧数据
	SelectList(ctx context.Context, request vo.VideoGetListRequest) (ret []model.Video, total int64, err error) // 获取番剧列表
	SelectOne(ctx context.Context, id int64, bangumiId int64) (ret model.Video, err error)                      // 根据bangumiID获取番剧
	DeleteList(ctx context.Context, idList []int64) (err error)                                                 // 批量删除番剧

	/* qbittorrent */
	UpdateRss(ctx context.Context, id int64, rss string) (err error) // 根据id更新rss链接
	GetQBLogs(ctx context.Context) (res []model.QBLog, err error)    // 获取qbittorrent的日志

	/* 视频 */
	InsertItemList(ctx context.Context, item []model.VideoItem) (cnt int64, err error)            // 插入视频数据
	SelectItemList(ctx context.Context, id int64) (ret []model.VideoItem, total int64, err error) // 根据番剧ID获取视频列表
	DeleteItemList(ctx context.Context, idList []int64) (err error)                               // 删除视频数据
}
