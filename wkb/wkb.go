package wkg

import (
	"encoding/binary"
	"errors"
	"math"
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
	Uint32Size    = int(unsafe.Sizeof(uint32(0)))
	Float64Size   = int(unsafe.Sizeof(float64(0)))
	PointSize     = int(unsafe.Sizeof(Point{}))
)

var (
	ErrInvalidStorage   = errors.New("Invalid storage type or size")
	ErrUnsupportedValue = errors.New("Unsupported value")
)

type Point struct {
	X, Y float64
}

type LineString []Point
type Polygon []LinearRing
type MultiPoint []Point
type MultiLineString []LineString
type MultiPolygon []Polygon
type Geometry struct {
	Kind  Kind
	Value interface{}
}
type GeometryCollection []Geometry

type LinearRing []Point

func (p *Point) Scan(src interface{}) error {
	b, dec, err := header(src, GeomPoint)
	if err != nil {
		return err
	}

	if len(b) < PointSize {
		return ErrInvalidStorage
	}

	b, p.X = readFloat64(b, dec)
	_, p.Y = readFloat64(b, dec)

	return nil
}

func (ls *LineString) Scan(src interface{}) error {
	b, dec, err := header(src, GeomLineString)
	if err != nil {
		return err
	}

	_, *ls, err = readPoints(b, dec)
	return err
}

func (p *Polygon) Scan(src interface{}) error {
	b, dec, err := header(src, GeomPolygon)
	if err != nil {
		return err
	}

	_, *p, err = readPolygon(b, dec)
	return err
}

func (mp *MultiPoint) Scan(src interface{}) error {
	b, dec, err := header(src, GeomMultiPoint)
	if err != nil {
		return err
	}

	_, *mp, err = readMultiPoint(b, dec)
	return err
}

func (mls *MultiLineString) Scan(src interface{}) error {
	b, dec, err := header(src, GeomMultiLineString)
	if err != nil {
		return err
	}

	_, *mls, err = readMultiLineString(b, dec)
	return err
}

func (mp *MultiPolygon) Scan(src interface{}) error {
	b, dec, err := header(src, GeomMultiPolygon)
	if err != nil {
		return err
	}

	_, *mp, err = readMultiPolygon(b, dec)
	return err
}

func (g *Geometry) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return ErrInvalidStorage
	}

	err := error(nil)
	_, *g, err = readGeometry(b)
	return err
}

func readUint32(b []byte, dec binary.ByteOrder) ([]byte, uint32) {
	return b[Uint32Size:], dec.Uint32(b)
}

func readCount(b []byte, dec binary.ByteOrder) ([]byte, int, error) {
	if len(b) < Uint32Size {
		return nil, 0, ErrInvalidStorage
	}
	b, n := readUint32(b, dec)
	return b, int(n), nil
}

func readFloat64(b []byte, dec binary.ByteOrder) ([]byte, float64) {
	return b[Float64Size:], math.Float64frombits(dec.Uint64(b))
}

func readPoint(b []byte, dec binary.ByteOrder) ([]byte, *Point, error) {
	if len(b) < PointSize {
		return nil, nil, ErrInvalidStorage
	}

	p := &Point{}
	b, p.X = readFloat64(b, dec)
	b, p.Y = readFloat64(b, dec)
	return b, p, nil
}

func readPoints(b []byte, dec binary.ByteOrder) ([]byte, []Point, error) {
	b, n, err := readCount(b, dec)
	if err != nil {
		return nil, nil, err
	}

	if len(b) < PointSize*n {
		return nil, nil, ErrInvalidStorage
	}

	p := make([]Point, n)
	for i := 0; i < n; i++ {
		b, p[i].X = readFloat64(b, dec)
		b, p[i].Y = readFloat64(b, dec)
	}

	return b, p, nil
}

func readLineString(b []byte, dec binary.ByteOrder) ([]byte, LineString, error) {
	return readPoints(b, dec)
}

func readMultiPoint(b []byte, dec binary.ByteOrder) ([]byte, MultiPoint, error) {
	return readPoints(b, dec)
}

func readMultiLineString(b []byte, dec binary.ByteOrder) ([]byte, MultiLineString, error) {
	b, n, err := readCount(b, dec)
	if err != nil {
		return nil, nil, err
	}

	mls := make([]LineString, n)
	for i := 0; i < n; i++ {
		b, mls[i], err = readLineString(b, dec)
		if err != nil {
			return nil, nil, err
		}
	}
	return b, mls, err
}

func readPolygon(b []byte, dec binary.ByteOrder) ([]byte, Polygon, error) {
	if len(b) < Uint32Size {
		return nil, nil, ErrInvalidStorage
	}

	b, n, err := readCount(b, dec)
	if err != nil {
		return nil, nil, err
	}

	lr := make([]LinearRing, n)
	for i := 0; i < n; i++ {
		b, lr[i], err = readPoints(b, dec)
		if err != nil {
			return nil, nil, err
		}
	}
	return b, lr, nil
}

func readMultiPolygon(b []byte, dec binary.ByteOrder) ([]byte, MultiPolygon, error) {
	b, n, err := readCount(b, dec)
	if err != nil {
		return nil, nil, err
	}

	mp := make([]Polygon, n)
	for i := 0; i < n; i++ {
		b, mp[i], err = readPolygon(b, dec)
		if err != nil {
			return nil, nil, err
		}
	}

	return b, mp, nil
}

func readGeometry(b []byte) ([]byte, Geometry, error) {
	g := Geometry{}
	if len(b) < HeaderSize {
		return nil, g, ErrInvalidStorage
	}

	dec := byteOrder(b[0])
	if dec == nil {
		return nil, g, ErrInvalidStorage
	}

	err := error(nil)

	b, kind := readUint32(b[ByteOrderSize:], dec)

	switch kind {
	case GeomPoint:
		b, g.Value, err = readPoint(b, dec)
	case GeomLineString:
		b, g.Value, err = readLineString(b, dec)
	case GeomPolygon:
		b, g.Value, err = readPolygon(b, dec)
	case GeomMultiPoint:
		b, g.Value, err = readMultiPoint(b, dec)
	case GeomMultiLineString:
		b, g.Value, err = readMultiLineString(b, dec)
	case GeomMultiPolygon:
		b, g.Value, err = readMultiPolygon(b, dec)
	case GeomCollection:
		b, g.Value, err = readGeometryCollection(b, dec)
	default:
		err = ErrUnsupportedValue
	}

	if err != nil {
		return nil, g, err
	}

	return b, g, nil
}

func readGeometryCollection(b []byte, dec binary.ByteOrder) ([]byte, GeometryCollection, error) {
	b, n, err := readCount(b, dec)
	if err != nil {
		return nil, nil, err
	}

	gc := make([]Geometry, n)
	for i := 0; i < n; i++ {
		b, gc[i], err = readGeometry(b)
		if err != nil {
			return nil, nil, err
		}
	}

	return b, gc, nil
}

func header(src interface{}, tpe Kind) ([]byte, binary.ByteOrder, error) {
	b, ok := src.([]byte)
	if !ok {
		return nil, nil, ErrInvalidStorage
	}

	if len(b) < HeaderSize {
		return nil, nil, ErrInvalidStorage
	}

	dec := byteOrder(b[0])
	if dec == nil {
		return nil, nil, ErrUnsupportedValue
	}

	b, kind := readUint32(b[ByteOrderSize:], dec)
	if tpe != Kind(kind) {
		return nil, nil, ErrUnsupportedValue
	}

	return b, dec, nil
}

func byteOrder(b byte) binary.ByteOrder {
	switch b {
	case BigEndian:
		return binary.BigEndian
	case LittleEndian:
		return binary.LittleEndian
	default:
		return nil
	}
}
