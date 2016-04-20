package wkb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
}
