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
	g := Geometry{}
	if err := g.Scan(rawPoint); assert.NoError(t, err) {
		assert.Equal(t, Geometry{
			Kind:  GeomPoint,
			Value: Point{30, 10},
		}, g)
	}

	g = Geometry{}
	if err := g.Scan(rawMultiPoint); assert.NoError(t, err) {
		assert.Equal(t, Geometry{
			Kind:  GeomMultiPoint,
			Value: MultiPoint{{10, 40}, {40, 30}, {20, 20}, {30, 10}},
		}, g)
	}

	g = Geometry{}
	if err := g.Scan(rawLineString); assert.NoError(t, err) {
		assert.Equal(t, Geometry{
			Kind:  GeomLineString,
			Value: LineString{{30, 10}, {10, 30}, {40, 40}},
		}, g)
	}

	g = Geometry{}
	if err := g.Scan(rawMultiLineString); assert.NoError(t, err) {
		assert.Equal(t, Geometry{
			Kind: GeomMultiLineString,
			Value: MultiLineString{
				LineString{{10, 10}, {20, 20}, {10, 40}},
				LineString{{40, 40}, {30, 30}, {40, 20}, {30, 10}},
			},
		}, g)
	}

	g = Geometry{}
	if err := g.Scan(rawPolygon); assert.NoError(t, err) {
		assert.Equal(t, Geometry{
			Kind: GeomPolygon,
			Value: Polygon{
				LinearRing{{30, 10}, {40, 40}, {20, 40}, {10, 20}, {30, 10}},
			},
		}, g)
	}

	g = Geometry{}
	if err := g.Scan(rawMultiPolygon); assert.NoError(t, err) {
		assert.Equal(t, Geometry{
			Kind: GeomMultiPolygon,
			Value: MultiPolygon{
				Polygon{
					LinearRing{{30, 20}, {45, 40}, {10, 40}, {30, 20}},
				},
				Polygon{
					LinearRing{{15, 5}, {40, 10}, {10, 20}, {5, 10}, {15, 5}},
				},
			},
		}, g)
	}

	g = Geometry{}
	if err := g.Scan(rawGeometryCollection); assert.NoError(t, err) {
		assert.Equal(t, Geometry{
			Kind: GeomCollection,
			Value: GeometryCollection{
				Geometry{
					Kind:  GeomPoint,
					Value: Point{4, 6},
				},
				Geometry{
					Kind:  GeomLineString,
					Value: LineString{{4, 6}, {7, 10}},
				},
			},
		}, g)
	}
}

func TestGeometryCollection(t *testing.T) {
	gc := GeometryCollection{}
	if err := gc.Scan(rawGeometryCollection); assert.NoError(t, err) {
		assert.Equal(t, GeometryCollection{
			Geometry{
				Kind:  GeomPoint,
				Value: Point{4, 6},
			},
			Geometry{
				Kind:  GeomLineString,
				Value: LineString{{4, 6}, {7, 10}},
			},
		}, gc)
	}
}
