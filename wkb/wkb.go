package wkb

import (
	"bytes"
	"errors"
	"unsafe"
)

type ByteOrder byte

const (
	BigEndian = iota
	LittleEndian
)

type Kind uint32

const (
	_ = iota
	GeomPoint
	GeomLineString
	GeomPolygon
	GeomMultiPoint
	GeomMultiLineString
	GeomMultiPolygon
	GeomCollection
)

const (
	ByteOrderSize = int(unsafe.Sizeof(ByteOrder(0)))
	GeomTypeSize  = int(unsafe.Sizeof(Kind(0)))
	HeaderSize    = ByteOrderSize + GeomTypeSize
	CountSize     = int(unsafe.Sizeof(uint32(0)))
	Float64Size   = int(unsafe.Sizeof(float64(0)))
	PointSize     = int(unsafe.Sizeof(Point{}))
)

var (
	ErrInvalidStorage   = errors.New("Invalid storage type or size")
	ErrUnsupportedValue = errors.New("Unsupported value")
)

type Geometry interface {
	ByteSize() int
	Write(*bytes.Buffer)
}

type LineString Points
type Polygon []LinearRing
type MultiPoint Points
type MultiLineString []LineString
type MultiPolygon []Polygon
type GeometryCollection []Geometry

type LinearRing Points
type Points []Point
