package wkb

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
)

func (p *Polygon) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return ErrInvalidStorage
	}

	_, tmp, err := ReadPolygon(b)
	if err != nil {
		return err
	}

	*p = tmp
	return err
}

func (p Polygon) Value() (driver.Value, error) {
	buf := bytes.NewBuffer(make([]byte, 0, p.ByteSize()))
	p.Write(buf)
	return buf.Bytes(), nil
}

func ReadPolygon(b []byte) ([]byte, Polygon, error) {
	if len(b) < HeaderSize+CountSize {
		return nil, nil, ErrInvalidStorage
	}

	b, dec, err := header(b, GeomPolygon)
	if err != nil {
		return nil, nil, err
	}

	b, n := readCount(b, dec)

	p := make([]LinearRing, n)
	for i := 0; i < n; i++ {
		if len(b) < CountSize {
			return nil, nil, ErrInvalidStorage
		}

		b, p[i], err = readLinearRing(b, dec)
		if err != nil {
			return nil, nil, err
		}
	}

	return b, p, nil
}

func (p Polygon) ByteSize() int {
	size := HeaderSize + CountSize
	for _, lr := range p {
		size += lr.byteSize()
	}
	return size
}

func (p Polygon) Write(buf *bytes.Buffer) {
	writeHeader(buf, GeomPolygon)
	writeCount(buf, len(p))
	for _, lr := range p {
		lr.write(buf)
	}
}

func (mp *MultiPolygon) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return ErrInvalidStorage
	}

	_, tmp, err := ReadMultiPolygon(b)
	if err != nil {
		return err
	}

	*mp = tmp
	return nil
}

func (mp MultiPolygon) Value() (driver.Value, error) {
	buf := bytes.NewBuffer(make([]byte, 0, mp.ByteSize()))
	mp.Write(buf)
	return buf.Bytes(), nil
}

func ReadMultiPolygon(b []byte) ([]byte, MultiPolygon, error) {
	if len(b) < HeaderSize+CountSize {
		return nil, nil, ErrInvalidStorage
	}

	b, dec, err := header(b, GeomMultiPolygon)
	if err != nil {
		return nil, nil, err
	}

	b, n := readCount(b, dec)

	mp := make([]Polygon, n)
	for i := 0; i < n; i++ {
		b, mp[i], err = ReadPolygon(b)
		if err != nil {
			return nil, nil, err
		}
	}

	return b, mp, nil
}

func (mp MultiPolygon) ByteSize() int {
	size := HeaderSize + CountSize
	for _, p := range mp {
		size += p.ByteSize()
	}
	return size
}

func (mp MultiPolygon) Write(buf *bytes.Buffer) {
	writeHeader(buf, GeomMultiPolygon)
	writeCount(buf, len(mp))
	for _, p := range mp {
		p.Write(buf)
	}
}

func readLinearRing(b []byte, dec binary.ByteOrder) ([]byte, LinearRing, error) {
	b, pts, err := readPoints(b, dec)
	return b, LinearRing(pts), err
}

func (lr LinearRing) byteSize() int {
	return Points(lr).byteSize()
}

func (lr LinearRing) write(buf *bytes.Buffer) {
	Points(lr).write(buf)
}
