package model

type QBLog struct {
	Id        int64  `form:"id" json:"id"`
	Message   string `form:"message" json:"message"`
	Timestamp int64  `form:"timestamp" json:"timestamp"`
	Type      int64  `form:"type" json:"type"`
}

type QBRule struct {
	Enabled                   bool          `json:"enabled"`
	MustContain               string        `json:"mustContain"`               // 必须包含
	MustNotContain            string        `json:"mustNotContain"`            // 必须不包含
	UseRegex                  bool          `json:"useRegex"`                  // 启用正则
	EpisodeFilter             string        `json:"episodeFilter"`             // 剧集过滤
	SmartFilter               bool          `json:"smartFilter"`               // 是否开启智能剧集过滤
	PreviouslyMatchedEpisodes []interface{} `json:"previouslyMatchedEpisodes"` // ??
	AffectedFeeds             []string      `json:"affectedFeeds"`             // 应用rss
	IgnoreDays                int           `json:"ignoreDays"`                // 忽略多少天前的
	LastMatch                 string        `json:"lastMatch"`                 // 最后匹配
	AddPaused                 bool          `json:"addPaused"`                 // true: 添加后不会立即下载
	AssignedCategory          string        `json:"assignedCategory"`          // 下载分类
	SavePath                  string        `json:"savePath"`                  // 下载后的保存路径
}

type QBSetRule struct {
	RuleName string `json:"ruleName"`
	RuleDef  string `json:"ruleDef"`
}
