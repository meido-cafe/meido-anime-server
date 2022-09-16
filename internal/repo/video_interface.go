package repo

import (
	"context"
	"meido-anime-server/internal/model"
	"meido-anime-server/internal/model/vo"
)

type VideoInterface interface {
	/* 番剧 */
	VideoInsertOne(ctx context.Context, video model.Video) (err error)                                               // 插入一条番剧数据
	VideoSelectList(ctx context.Context, request vo.VideoGetListRequest) (ret []model.Video, total int64, err error) // 获取番剧列表
	VideoSelectListById(ctx context.Context, idList []int64) (ret []model.Video, err error)                          // 获取番剧列表
	VideoSelectOne(ctx context.Context, id int64, bangumiId int64) (ret model.Video, err error)                      // 根据bangumiID获取番剧
	VideoDeleteList(ctx context.Context, idList []int64) (err error)                                                 // 批量删除番剧
	VideoUpdateCategoryList(ctx context.Context, videoList []model.Video) (err error)                                // 批量更新分类与硬链接目录路径

	/* 分类 */
	CategoryInsertOne(ctx context.Context, cagetory model.Cagetory) (err error)               // 创建一个分类
	CategorySelectOne(ctx context.Context, id int64) (ret model.Cagetory, err error)          // 获取一个分类
	CategorySelectOneByName(ctx context.Context, name string) (ret model.Cagetory, err error) // 根据分类名获取一个分类
	CategorySelectList(ctx context.Context) (ret []model.Cagetory, err error)                 // 获取分类列表
	CategoryDeleteOne(ctx context.Context, id int64) (err error)                              // 删除一个分类
	CategoryUpdateName(ctx context.Context, id int64, name string) (err error)                // 更新一个分类名称

	/* 视频 */
	VideoItemInsertList(ctx context.Context, videoItemList []model.VideoItem) (cnt int64, err error)      // 插入视频数据
	VideoItemSelectList(ctx context.Context, pid []int64) (ret []model.VideoItem, total int64, err error) // 根据番剧ID获取视频列表
	VideoItemDeleteList(ctx context.Context, idList []int64) (err error)                                  // 删除视频数据
	VideoItemDeleteListByPid(ctx context.Context, pid []int64) (err error)                                // 根据PID删除视频
	VideoItemUpdateLinkPathList(ctx context.Context, videoItemList []model.VideoItem) (err error)         // 批量更新硬链接路径
}
