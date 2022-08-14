package model

type Page struct {
	Page     int64 `json:"page"`
	PageSize int64 `json:"page_size"`
}

type Time struct {
	CreateTime int64 `json:"create_time"`
	UpdateTime int64 `json:"update_time"`
}
