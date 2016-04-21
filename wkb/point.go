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
	buf := bytes.NewBuffer(make([]byte, 0, p.byteSize()))
	p.write(buf)
	return buf.Bytes(), nil
}

func (p *Point) Scan(src interface{}) error {
	b, dec, err := header(src, GeomPoint)
	if err != nil {
		return err
	}

	if len(b) < PointSize {
		return ErrInvalidStorage
	}

	_, *p = readPoint(b, dec)
	return nil
}

func (mp *MultiPoint) Scan(src interface{}) error {
	b, dec, err := header(src, GeomMultiPoint)
	if err != nil {
		return err
	}

	_, *mp, err = readMultiPoint(b, dec)
	return err
}

func (mp MultiPoint) Value() (driver.Value, error) {
	buf := bytes.NewBuffer(make([]byte, 0, mp.byteSize()))
	mp.write(buf)
	return buf.Bytes(), nil
}

func readPoint(b []byte, dec binary.ByteOrder) ([]byte, Point) {
	p := Point{}
	b, p.X = readFloat64(b, dec)
	b, p.Y = readFloat64(b, dec)
	return b, p
}

func (p *Point) byteSize() int {
	return HeaderSize + PointSize
}

func (p *Point) write(buf *bytes.Buffer) {
	writeHeader(buf, GeomPoint)
	writeFloat64(buf, p.X)
	writeFloat64(buf, p.Y)
}

func readMultiPoint(b []byte, dec binary.ByteOrder) ([]byte, MultiPoint, error) {
	b, n, err := readCount(b, dec)
	if err != nil {
		return nil, nil, err
	}

	if len(b) < (HeaderSize+PointSize)*n {
		return nil, nil, ErrInvalidStorage
	}

	mp := make([]Point, n)
	for i := 0; i < n; i++ {
		b, dec, err = byteHeader(b, GeomPoint)
		if err != nil {
			return nil, nil, err
		}
		b, mp[i] = readPoint(b, dec)
	}

	return b, mp, nil
}

func (mp MultiPoint) byteSize() int {
	return HeaderSize + Uint32Size + len(mp)*(HeaderSize+PointSize)
}

func (mp MultiPoint) write(buf *bytes.Buffer) {
	writeHeader(buf, GeomMultiPoint)
	writeCount(buf, len(mp))
	for _, p := range mp {
		p.write(buf)
	}
}

func readPoints(b []byte, dec binary.ByteOrder) ([]byte, Points, error) {
	b, n, err := readCount(b, dec)
	if err != nil {
		return nil, nil, err
	}

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
	return Uint32Size + len(pts)*PointSize
}

func (pts Points) write(buf *bytes.Buffer) {
	writeCount(buf, len(pts))
	for _, p := range pts {
		writeFloat64(buf, p.X)
		writeFloat64(buf, p.Y)
	}
}
