package model

type Video struct {
	Id        int64  `db:"id"`
	BangumiId int64  `db:"bangumi_id"` // bangumi的ID
	Origin    int64  `db:"origin"`     // 来源方式 1:订阅 2: 手动
	Category  int64  `db:"category"`   // 类别 1:tv,2:剧场版,3:OVA
	Title     string `db:"title"`      // 标题
	Season    int64  `db:"season"`     // 第几季
	Cover     string `db:"cover"`      // 封面图
	Total     int64  `db:"total"`      // 总集数
	RssUrl    string `db:"rss_url"`    // rss链接
	PlayTime  int64  `db:"play_time"`  // 放送时间
	SourceDir string `db:"source_dir"` // 下载路径
	Time
}

type VideoItem struct {
	Id         int64  `db:"id"`
	Pid        int64  `db:"pid"`         // video表的ID
	Episode    int64  `db:"episode"`     // 第几集
	SourcePath string `db:"source_path"` // 下载路径
	LinkPath   string `db:"link_path"`   // 硬链接路径
	Time
}
