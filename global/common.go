package global

import (
	"database/sql"
	"sync"
)

var DB *sql.DB
var DBOnce sync.Once
