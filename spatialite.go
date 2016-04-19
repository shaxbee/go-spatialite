package spatialite

import (
	"database/sql"

	"github.com/mattn/go-sqlite3"
)

func init() {
	sql.Register("spatialite", &sqlite3.SQLiteDriver{
		Extensions: []string{"spatialite"},
	})
}
