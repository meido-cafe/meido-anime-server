package global

import (
	"meido-anime-server/internal/model"
)

var TokenCache = make(map[string]*model.User)
