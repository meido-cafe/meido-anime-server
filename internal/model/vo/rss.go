package vo

import "meido-anime-server/internal/model"

type GetRssInfoMikanRequest struct {
	SubjectName string `form:"subject_name" json:"subject_name"`
}
type MikanGroup struct {
	Gid       int64  `json:"gid"`
	GroupName string `json:"group_name"`
}

type GetRssMiaknInfoResponse struct {
	Mid   int64        `json:"mid"`
	Group []MikanGroup `json:"group"`
}

type GetRssSearchRequest struct {
	SubjectName string `from:"subject_name" json:"subject_name"`
}

type GetRssSearchResponse struct {
	Url  string     `json:"url"`  // rss url
	Feed model.Feed `json:"feed"` // rss解析结果
}

type GetRssSubjectRequest struct {
	MikanId      int64 `form:"mikan_id" json:"mikan_id"`
	MikanGroupId int64 `form:"mikan_group_id" json:"mikan_group_id"`
}
type GetRssSubjectResponse struct {
	Url  string     `json:"url"`  // rss url
	Feed model.Feed `json:"feed"` // rss解析结果
}
