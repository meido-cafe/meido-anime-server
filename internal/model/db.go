package model

type Video struct {
	Id        int64  `json:"id" db:"id"`
	BangumiId int64  `json:"bangumi_id" db:"bangumi_id"` // bangumi的ID
	Origin    int64  `json:"origin" db:"origin"`         // 来源方式 1:订阅 2: 手动
	Category  int64  `json:"category" db:"category"`     // 类别 1:tv,2:剧场版,3:OVA
	Title     string `json:"title" db:"title"`           // 标题
	Season    int64  `json:"season" db:"season"`         // 第几季
	Cover     string `json:"cover" db:"cover"`           // 封面图
	Total     int64  `json:"total" db:"total"`           // 总集数
	RssUrl    string `json:"rss_url" db:"rss_url"`       // rss链接
	RuleName  string `json:"rule_name" db:"rule_name"`   // 规则名称
	PlayTime  int64  `json:"play_time" db:"play_time"`   // 放送时间
	SourceDir string `json:"source_dir" db:"source_dir"` // 下载目录
	LinkDir   string `json:"link_dir" db:"link_dir"`     // 硬链接目录
	Time
}

type VideoItem struct {
	Id         int64  `json:"id" db:"id"`
	Pid        int64  `json:"pid" db:"pid"`                 // video表的ID
	Episode    int64  `json:"episode" db:"episode"`         // 第几集
	SourcePath string `json:"source_path" db:"source_path"` // 下载路径
	LinkPath   string `json:"link_path" db:"link_path"`     // 硬链接路径
	Time
}

type Cagetory struct {
	Id     int64  `json:"id" db:"id"`
	Name   string `json:"name" db:"name"`     // 分类名称
	Origin int64  `json:"origin" db:"origin"` // 分类来源 1为内置 2为自定义
	Time
}
