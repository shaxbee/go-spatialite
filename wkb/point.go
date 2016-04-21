package wkb

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
)

type Point struct {
	X, Y float64
}

func (p Point) Equal(other Point) bool {
	return p.X == other.X && p.Y == other.Y
}

func (p Point) Value() (driver.Value, error) {
	buf := bytes.NewBuffer(make([]byte, 0, p.ByteSize()))
	p.Write(buf)
	return buf.Bytes(), nil
}

func (p *Point) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return ErrInvalidStorage
	}

	_, tmp, err := ReadPoint(b)
	if err != nil {
		return err
	}

	*p = tmp
	return nil
}

func (p Point) ByteSize() int {
	return HeaderSize + PointSize
}

func (p Point) Write(buf *bytes.Buffer) {
	writeHeader(buf, GeomPoint)
	writeFloat64(buf, p.X)
	writeFloat64(buf, p.Y)
}

func ReadPoint(b []byte) ([]byte, Point, error) {
	p := Point{}
	if len(b) < HeaderSize+PointSize {
		return nil, p, ErrInvalidStorage
	}

	b, dec, err := header(b, GeomPoint)
	if err != nil {
		return nil, p, err
	}

	b, p.X = readFloat64(b, dec)
	b, p.Y = readFloat64(b, dec)
	return b, p, nil
}

func (mp *MultiPoint) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return ErrInvalidStorage
	}

	_, tmp, err := ReadMultiPoint(b)
	if err != nil {
		return err
	}

	*mp = tmp
	return nil
}

func (mp MultiPoint) Value() (driver.Value, error) {
	buf := bytes.NewBuffer(make([]byte, 0, mp.ByteSize()))
	mp.Write(buf)
	return buf.Bytes(), nil
}

func ReadMultiPoint(b []byte) ([]byte, MultiPoint, error) {
	if len(b) < HeaderSize+CountSize {
		return nil, nil, ErrInvalidStorage
	}

	b, dec, err := header(b, GeomMultiPoint)
	if err != nil {
		return nil, nil, err
	}

	b, n := readCount(b, dec)

	mp := make([]Point, n)
	for i := 0; i < n; i++ {
		b, mp[i], err = ReadPoint(b)
		if err != nil {
			return nil, nil, err
		}
	}

	return b, mp, nil
}

func (mp MultiPoint) ByteSize() int {
	return HeaderSize + Points(mp).byteSize()
}

func (mp MultiPoint) Write(buf *bytes.Buffer) {
	writeHeader(buf, GeomMultiPoint)
	writeCount(buf, len(mp))
	for _, p := range mp {
		p.Write(buf)
	}
}

func readPoint(b []byte, dec binary.ByteOrder) ([]byte, Point) {
	p := Point{}
	b, p.X = readFloat64(b, dec)
	b, p.Y = readFloat64(b, dec)
	return b, p
}

func readPoints(b []byte, dec binary.ByteOrder) ([]byte, Points, error) {
	b, n := readCount(b, dec)

	if len(b) < PointSize*n {
		return nil, nil, ErrInvalidStorage
	}

	p := make([]Point, n)
	for i := 0; i < n; i++ {
		b, p[i] = readPoint(b, dec)
	}

	return b, p, nil
}

func (pts Points) byteSize() int {
	return CountSize + len(pts)*PointSize
}

func (pts Points) write(buf *bytes.Buffer) {
	writeCount(buf, len(pts))
	for _, p := range pts {
		writeFloat64(buf, p.X)
		writeFloat64(buf, p.Y)
	}
}
