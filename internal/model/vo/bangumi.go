package vo

type GetIndexRequest struct {
	Page     int `json:"page"`      // 页码
	PageSize int `json:"page_size"` // 分页size
	GetIndexRequestBody
}

type GetIndexRequestBody struct {
	Keyword string `json:"keyword"` // 关键词
	Sort    string `json:"sort"`    // 排序
	Filter  struct {
		Type    []int64  `json:"type"`     // subject类型 列表 或关系
		Tag     []string `json:"tag"`      // subject标签 且关系
		AirDate []string `json:"air_date"` //播出日期/发售日期。且 关系。    [">=2020-07-01","<2022-10-01"]
		Rating  []string `json:"rating"`   // 用于搜索指定评分的条目。且 关系。 [">=6","<8"]
		Rank    []string `json:"rank"`     // 用于搜索指定排名的条目。且 关系。 [">10","<=18"]
		Nsfw    bool     `json:"nsfw"`     //  使用 include 包含NSFW搜索结果。默认排除搜索NSFW条目。无权限情况下忽略此选项，不会返回NSFW条目。
	} `json:"filter"`
}

type GetSubjectRequest struct {
	Id int `form:"id" json:"id"`
}
type SearchRequest struct {
	Name  string `form:"name" json:"name"`
	Class int    `form:"class" json:"class"`
}

type GetSubjectCharactersRequest struct {
	Id int `form:"id" json:"id"`
}
