package wkg

type ByteOrder uint8

const (
	BigEndian = iota
	LittleEndian
)

type Coord uint32

const (
	Coord2D = iota << 4
	CoordZ
	CoordM
	CoordZM
)

type GeomType uint32

const (
	_ = iota
	GeomPoint
	GeomLineString
	GeomPolygon
	GeomMultiPoint
	GeomMultiLineString
	GeomMultiPolygon
	GeomGeometryCollection
)
