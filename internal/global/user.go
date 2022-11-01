package global

import (
	"meido-anime-server/internal/model"
)

var TokenCache = make(map[string]*model.User)
var Salt = "me1d0@M01"
