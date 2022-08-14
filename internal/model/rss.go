package model

type Feed struct {
	Title string     `json:"title"` // rss标题
	Desc  string     `json:"desc"`  // rss描述
	Items []FeedItem `json:"items"` // 种子
}

type FeedItem struct {
	Title   string `json:"title"`    // 种子标题
	Desc    string `json:"desc"`     // 种子描述
	PubDate string `json:"pub_date"` // 发布日期
	Url     string `json:"url"`      // 种子链接
	Length  int    `json:"length"`   // 内容大小, 单位B
}
