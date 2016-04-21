package wkb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var rawGeometryCollection = []byte{
	0x01, 0x07, 0x00, 0x00, 0x00, // header
	0x02, 0x00, 0x00, 0x00, // numgeometry - 2
	0x01, 0x01, 0x00, 0x00, 0x00, // geometry 1 - point
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x40,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x18, 0x40,
	0x01, 0x02, 0x00, 0x00, 0x00, // geometry 2 - linestring
	0x02, 0x00, 0x00, 0x00, // numpoints - 2
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x40,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x18, 0x40,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x1c, 0x40,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40,
}

func TestGeometry(t *testing.T) {
	invalid := []struct {
		err error
		b   []byte
	}{
		// empty
		{
			ErrInvalidStorage,
			[]byte{},
		},
		// invalid byte order
		{
			ErrInvalidStorage,
			[]byte{0x02, 0x01, 0x00, 0x00, 0x00},
		},
		// no payload
		{
			ErrInvalidStorage,
			[]byte{0x01, 0x01, 0x00, 0x00, 0x00},
		},
		// invalid type
		{
			ErrUnsupportedValue,
			[]byte{0x01, 0x42, 0x00, 0x00, 0x00},
		},
	}

	for _, e := range invalid {
		if _, err := New(e.b); assert.Error(t, err) {
			assert.Exactly(t, e.err, err)
		}
	}

	if g, err := New(rawPoint); assert.NoError(t, err) {
		assert.Equal(t, Point{30, 10}, g)
	}

	if g, err := New(rawMultiPoint); assert.NoError(t, err) {
		assert.Equal(t, MultiPoint{{10, 40}, {40, 30}, {20, 20}, {30, 10}}, g)
	}

	if g, err := New(rawLineString); assert.NoError(t, err) {
		assert.Equal(t, LineString{{30, 10}, {10, 30}, {40, 40}}, g)
	}

	if g, err := New(rawMultiLineString); assert.NoError(t, err) {
		assert.Equal(t, MultiLineString{
			LineString{{10, 10}, {20, 20}, {10, 40}},
			LineString{{40, 40}, {30, 30}, {40, 20}, {30, 10}},
		}, g)
	}

	if g, err := New(rawPolygon); assert.NoError(t, err) {
		assert.Equal(t, Polygon{
			LinearRing{{30, 10}, {40, 40}, {20, 40}, {10, 20}, {30, 10}},
		}, g)
	}

	if g, err := New(rawMultiPolygon); assert.NoError(t, err) {
		assert.Equal(t, MultiPolygon{
			Polygon{
				LinearRing{{30, 20}, {45, 40}, {10, 40}, {30, 20}},
			},
			Polygon{
				LinearRing{{15, 5}, {40, 10}, {10, 20}, {5, 10}, {15, 5}},
			},
		}, g)
	}

	if g, err := New(rawGeometryCollection); assert.NoError(t, err) {
		assert.Equal(t, GeometryCollection{
			Point{4, 6},
			LineString{{4, 6}, {7, 10}},
		}, g)
	}
}

func TestGeometryCollection(t *testing.T) {
	invalid := []struct {
		err error
		b   []byte
	}{
		// invalid type
		{
			ErrUnsupportedValue,
			[]byte{0x01, 0x42, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			// no payload
			ErrInvalidStorage,
			[]byte{0x01, 0x07, 0x00, 0x00, 0x00},
		},
		// no element payload
		{
			ErrInvalidStorage,
			[]byte{
				0x01, 0x07, 0x00, 0x00, 0x00, // header
				0x02, 0x00, 0x00, 0x00, // numgeometry - 2
				0x01, 0x01, 0x00, 0x00, 0x00, // geometry 1 - point
			},
		},
	}

	for _, e := range invalid {
		gc := GeometryCollection{}
		if err := gc.Scan(e.b); assert.Error(t, err) {
			assert.Exactly(t, e.err, err)
		}
	}

	if err := (&GeometryCollection{}).Scan(""); assert.Error(t, err) {
		assert.Exactly(t, ErrInvalidStorage, err)
	}

	gc := GeometryCollection{}
	if err := gc.Scan(rawGeometryCollection); assert.NoError(t, err) {
		assert.Equal(t, GeometryCollection{
			Point{4, 6},
			LineString{{4, 6}, {7, 10}},
		}, gc)
	}

	if raw, err := gc.Value(); assert.NoError(t, err) {
		assert.Equal(t, rawGeometryCollection, raw)
	}
}
