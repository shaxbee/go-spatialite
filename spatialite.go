package spatialite

import (
	"database/sql"
	"database/sql/driver"

	"github.com/mattn/go-sqlite3"
)

func init() {
	sql.Register("spatialite", &sqlite3.SQLiteDriver{
		Extensions: []string{"mod_spatialite"},
		ConnectHook: func(c *sqlite3.SQLiteConn) error {
			_, err := c.Exec("SELECT InitSpatialMetadata()", []driver.Value{})
			return err
		},
	})
}
