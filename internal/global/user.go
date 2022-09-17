package global

import (
	"meido-anime-server/internal/model"
)

var TokenCache = make(map[string]*model.User)

const Salt = "me1d0@nlm1"
