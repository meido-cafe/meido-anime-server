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
	Id        int64 `form:"id" json:"id"`
	BangumiId int64 `form:"bangumi_id" json:"bangumi_id"`
}

type VideoGetOneResponse struct {
	Video model.Video `json:"video"`
}
type VideoSubscribeRequest struct {
	BangumiId int64  `json:"bangumi_id"` // bangumi的ID
	Title     string `json:"title"`      // 标题
	Category  int64  `json:"category"`   // 类别 1:tv 2:剧场版 3:OVA
	Season    int64  `json:"season"`     // 第几季
	Cover     string `json:"cover"`      // 封面图
	Total     int64  `json:"total"`      // 总集数
	PlayTime  int64  `json:"play_time"`  // 放送时间

	RssUrl string `json:"rss_url"` // rss链接
}

type VideoAddTorrent struct {
	BangumiId int64  `json:"bangumi_id"` // bangumi的ID
	Title     string `json:"title"`      // 标题
	Category  int64  `json:"category"`   // 类别 1:tv 2:剧场版 3:OVA
	Season    int64  `json:"season"`     // 第几季
	Cover     string `json:"cover"`      // 封面图
	Total     int64  `json:"total"`      // 总集数
	PlayTime  int64  `json:"play_time"`  // 放送时间

	TorrentList []string `json:"torrent_list"` // 种子列表
}
