package vo

type GetIndexRequest struct {
	Year  int    `form:"year" json:"year"`
	Month int    `form:"month" json:"month"`
	Page  int    `form:"page" json:"page"`
	Type  string `form:"type" json:"type"`
	Sort  string `form:"sort" json:"sort"`
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
