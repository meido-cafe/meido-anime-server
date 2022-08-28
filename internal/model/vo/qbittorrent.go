package vo

import "meido-anime-server/internal/model"

type GetQBLogsResponseItem struct {
	model.QBLog
	Time string `json:"time"`
}
type GetQBLogsResponse struct {
	Items []GetQBLogsResponseItem `json:"items"`
	Total int64                   `json:"total"`
}
