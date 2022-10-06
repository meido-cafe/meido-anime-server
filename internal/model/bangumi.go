package model

type BangumiCalendar struct {
	Weekday struct {
		En string `json:"en"`
		Cn string `json:"cn"`
		Ja string `json:"ja"`
		Id int    `json:"id"`
	} `json:"weekday"`
	Items []struct {
		Id         int    `json:"id"`
		Url        string `json:"url"`
		Type       int    `json:"type"`
		Name       string `json:"name"`
		NameCn     string `json:"name_cn"`
		Summary    string `json:"summary"`
		AirDate    string `json:"air_date"`
		AirWeekday int    `json:"air_weekday"`
		Rating     struct {
			Total int `json:"total"`
			Count struct {
				Field1  int `json:"1"`
				Field2  int `json:"2"`
				Field3  int `json:"3"`
				Field4  int `json:"4"`
				Field5  int `json:"5"`
				Field6  int `json:"6"`
				Field7  int `json:"7"`
				Field8  int `json:"8"`
				Field9  int `json:"9"`
				Field10 int `json:"10"`
			} `json:"count"`
			Score float64 `json:"score"`
		} `json:"rating"`
		Images struct {
			Large  string `json:"large"`
			Common string `json:"common"`
			Medium string `json:"medium"`
			Small  string `json:"small"`
			Grid   string `json:"grid"`
		} `json:"images"`
		Collection struct {
			Doing int `json:"doing"`
		} `json:"collection"`
		Rank int `json:"rank,omitempty"`
	} `json:"items"`
}

type BangumiSubject struct {
	Date     string `json:"date"`
	Platform string `json:"platform"`
	Images   struct {
		Small  string `json:"small"`
		Grid   string `json:"grid"`
		Large  string `json:"large"`
		Medium string `json:"medium"`
		Common string `json:"common"`
	} `json:"images"`
	Summary string `json:"summary"`
	Name    string `json:"name"`
	NameCn  string `json:"name_cn"`
	Tags    []struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	} `json:"tags"`
	Infobox []struct {
		Key   string      `json:"key"`
		Value interface{} `json:"value"`
	} `json:"infobox"`
	Rating struct {
		Rank  int `json:"rank"`
		Total int `json:"total"`
		Count struct {
			Field1  int `json:"1"`
			Field2  int `json:"2"`
			Field3  int `json:"3"`
			Field4  int `json:"4"`
			Field5  int `json:"5"`
			Field6  int `json:"6"`
			Field7  int `json:"7"`
			Field8  int `json:"8"`
			Field9  int `json:"9"`
			Field10 int `json:"10"`
		} `json:"count"`
		Score float64 `json:"score"`
	} `json:"rating"`
	TotalEpisodes int `json:"total_episodes"`
	Collection    struct {
		OnHold  int `json:"on_hold"`
		Dropped int `json:"dropped"`
		Wish    int `json:"wish"`
		Collect int `json:"collect"`
		Doing   int `json:"doing"`
	} `json:"collection"`
	Id      int  `json:"id"`
	Eps     int  `json:"eps"`
	Volumes int  `json:"volumes"`
	Locked  bool `json:"locked"`
	Nsfw    bool `json:"nsfw"`
	Type    int  `json:"type"`
}

type BangumiSearchSubjectItem struct {
	Id         int    `json:"id"`
	Url        string `json:"url"`
	Type       int    `json:"type"`
	Name       string `json:"name"`
	NameCn     string `json:"name_cn"`
	Summary    string `json:"summary"`
	AirDate    string `json:"air_date"`
	AirWeekday int    `json:"air_weekday"`
	Images     struct {
		Large  string `json:"large"`
		Common string `json:"common"`
		Medium string `json:"medium"`
		Small  string `json:"small"`
		Grid   string `json:"grid"`
	} `json:"images"`
}
type BangumiSearchSubjectResponse struct {
	Results int                        `json:"results"`
	List    []BangumiSearchSubjectItem `json:"list"`
}

type BangumiSubjectCharacter struct {
	Images struct {
		Small  string `json:"small"`
		Grid   string `json:"grid"`
		Large  string `json:"large"`
		Medium string `json:"medium"`
	} `json:"images"`
	Name     string `json:"name"`
	Relation string `json:"relation"`
	Actors   []struct {
		Images struct {
			Small  string `json:"small"`
			Grid   string `json:"grid"`
			Large  string `json:"large"`
			Medium string `json:"medium"`
		} `json:"images"`
		Name         string   `json:"name"`
		ShortSummary string   `json:"short_summary"`
		Career       []string `json:"career"`
		Id           int      `json:"id"`
		Type         int      `json:"type"`
		Locked       bool     `json:"locked"`
	} `json:"actors"`
	Type int `json:"type"`
	Id   int `json:"id"`
}

type BangumiIndexItem struct {
	Total    int    // 总集数
	Rank     int    // 排名
	Rat      int    // 评分
	PlayTime int    // 放送时间
	Title    string // 标题
	Cover    string // 封面图
}
type BangumiIndexResponse struct {
	Items     BangumiIndexItem `json:"items"`      // 番剧
	TotalPage int              `json:"total_page"` // 总页数
}
