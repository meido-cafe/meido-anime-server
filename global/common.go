package global

import (
	"github.com/jmoiron/sqlx"
	"sync"
)

var DB *sqlx.DB
var DBOnce sync.Once
