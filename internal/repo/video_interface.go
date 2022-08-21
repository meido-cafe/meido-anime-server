package repo

import (
	"context"
	"meido-anime-server/internal/model"
	"meido-anime-server/internal/model/vo"
)

type VideoInterface interface {
	InsertOne(ctx context.Context, video model.Video) (err error)                                               // 插入一条video数据
	SelectList(ctx context.Context, request vo.VideoGetListRequest) (ret []model.Video, total int64, err error) // 获取video列表
	SelectOne(ctx context.Context, id int64, bangumiId int64) (ret model.Video, err error)                      // 根据bangumiID获取video
	DeleteList(ctx context.Context, idList []int64) (err error)                                                 // 批量删除video

	UpdateRss(ctx context.Context, id int64, rss string) (err error) // 根据id更新rss链接
	DeleteRss(ctx context.Context, id int64) (err error)             // 删除rss链接
}
