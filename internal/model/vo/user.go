package vo

type LoginRequest struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}

type UpdateUsernameRequest struct {
	Username string `form:"username" json:"username"`
}
