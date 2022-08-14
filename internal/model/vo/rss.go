package vo

type GetInfoMikanRequest struct {
	SubjectName string `form:"subject_name" json:"subject_name"`
}
type MikanGroup struct {
	Gid       int64  `json:"gid"`
	GroupName string `json:"group_name"`
}

type GetMiaknInfoResponse struct {
	Mid   int64        `json:"mid"`
	Group []MikanGroup `json:"group"`
}
