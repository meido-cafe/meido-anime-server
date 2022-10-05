package vo

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
