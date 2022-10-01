package vo

import "meido-anime-server/internal/model"

type VideoGetListRequest struct {
	Title         string `form:"title"`           // 名称模糊搜索
	Origin        int64  `form:"origin"`          // 来源搜索 1: 订阅rss, 2: 种子
	Category      int64  `form:"category"`        // 类别搜索 1:tv 2:剧场版 3:OVA
	PlayStartTime int64  `form:"play_start_time"` // 放送时间范围的开始时间
	PlayEndTime   int64  `form:"play_end_time"`   // 放送时间范围的结束时间
	AddStartTime  int64  `form:"add_start_time"`  // 添加时间范围的开始时间
	AddEndTime    int64  `form:"add_end_time"`    // 添加时间范围的结束时间
	model.Page
}

type VideoGetListResponse struct {
	Items []model.Video `json:"items"`
	Total int64         `json:"total"`
}

type VideoGetOneRequest struct {
	Id        int64  `form:"id" json:"id"`
	BangumiId int64  `form:"bangumi_id" json:"bangumi_id"`
	Title     string `form:"title" json:"title"`
}

type VideoGetOneResponse struct {
	Video model.Video `json:"video"`
}
type VideoSubscribeRequest struct {
	BangumiId int64  `form:"bangumi_id" json:"bangumi_id"` // bangumi的ID
	Title     string `form:"title" json:"title"`           // 标题
	Category  int64  `form:"category" json:"category"`     // 类别 1:tv 2:剧场版 3:OVA
	Season    int64  `form:"season" json:"season"`         // 第几季
	Cover     string `form:"cover" json:"cover"`           // 封面图
	Total     int64  `form:"total" json:"total"`           // 总集数
	PlayTime  int64  `form:"play_time" json:"play_time"`   // 放送时间
	RssUrl    string `form:"rss_url" json:"rss_url"`       // RSS 链接

	MustContain    string `form:"must_contain" json:"must_contain"`         // 必须包含
	MustNotContain string `form:"must_not_contain" json:"must_not_contain"` // 必须不包含
	UseRegex       bool   `form:"use_regex" json:"use_regex"`               // 是否开启正则
	EpisodeFilter  string `form:"episode_filter" json:"episode_filter"`     // 剧集过滤
	SmartFilter    bool   `form:"smart_filter" json:"smart_filter"`         // 是否开启智能剧集过滤
}

type VideoAdd struct {
	VideoSubscribeRequest
	Mode        int      `json:"mode" form:"mode"`         // 订阅:1 种子:2
	TorrentList []string `json:"torrent_list" form:"mode"` // 种子列表
}

type DeleteVideoRequest struct {
	Id         int64 `form:"id" json:"id"`
	DeleteFile bool  `form:"delete_file" json:"delete_file"`
}
type UpdateRssRequest struct {
	Id  int64  `form:"id" json:"id"`
	Rss string `form:"rss" json:"rss"`
}

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

type UpdateVideoCategoryRequest struct {
	Ids      []int64 `form:"ids" json:"ids"`
	Category int64   `form:"category" json:"category"`
}

type DeleteCategoryRequest struct {
	Id int64 `form:"id" json:"id"`
}
type UpdateCategoryNameRequest struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}
