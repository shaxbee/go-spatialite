package spatialite

import (
	"database/sql"
	"testing"

	"github.com/shaxbee/go-spatialite/wkb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTable(t *testing.T) {
	db, err := sql.Open("spatialite", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec("CREATE TABLE poi (title TEXT, loc ST_Point)")
	assert.NoError(t, err)
}

func TestPoint(t *testing.T) {
	db, err := sql.Open("spatialite", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec("CREATE TABLE poi (title TEXT, loc ST_Point)")
	require.NoError(t, err)

	p1 := wkb.Point{10, 10}
	_, err = db.Exec("INSERT INTO poi(title, loc) VALUES (?, ST_PointFromWKB(?))", "foo", p1)
	assert.NoError(t, err)

	p2 := wkb.Point{}
	r := db.QueryRow("SELECT AsBinary(loc) AS loc FROM poi WHERE title=?", "foo")
	if err := r.Scan(&p2); assert.NoError(t, err) {
		assert.Equal(t, p1, p2)
	}
}
