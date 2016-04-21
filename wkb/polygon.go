package wkb

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
)

func (p *Polygon) Scan(src interface{}) error {
	b, dec, err := header(src, GeomPolygon)
	if err != nil {
		return err
	}

	_, *p, err = readPolygon(b, dec)
	return err
}

func (p Polygon) Value() (driver.Value, error) {
	buf := bytes.NewBuffer(make([]byte, 0, p.byteSize()))
	p.write(buf)
	return buf.Bytes(), nil
}

func (mp *MultiPolygon) Scan(src interface{}) error {
	b, dec, err := header(src, GeomMultiPolygon)
	if err != nil {
		return err
	}

	_, *mp, err = readMultiPolygon(b, dec)
	return err
}

func (mp MultiPolygon) Value() (driver.Value, error) {
	buf := bytes.NewBuffer(make([]byte, 0, mp.byteSize()))
	mp.write(buf)
	return buf.Bytes(), nil
}

func readPolygon(b []byte, dec binary.ByteOrder) ([]byte, Polygon, error) {
	b, n, err := readCount(b, dec)
	if err != nil {
		return nil, nil, err
	}

	lr := make([]LinearRing, n)
	for i := 0; i < n; i++ {
		b, lr[i], err = readLinearRing(b, dec)
		if err != nil {
			return nil, nil, err
		}
	}
	return b, lr, nil
}

func (p Polygon) byteSize() int {
	size := HeaderSize + Uint32Size
	for _, lr := range p {
		size += lr.byteSize()
	}
	return size
}

func (p Polygon) write(buf *bytes.Buffer) {
	writeHeader(buf, GeomPolygon)
	writeCount(buf, len(p))
	for _, lr := range p {
		lr.write(buf)
	}
}

func readMultiPolygon(b []byte, dec binary.ByteOrder) ([]byte, MultiPolygon, error) {
	b, n, err := readCount(b, dec)
	if err != nil {
		return nil, nil, err
	}

	mp := make([]Polygon, n)
	for i := 0; i < n; i++ {
		b, dec, err = byteHeader(b, GeomPolygon)
		if err != nil {
			return nil, nil, err
		}

		b, mp[i], err = readPolygon(b, dec)
		if err != nil {
			return nil, nil, err
		}
	}

	return b, mp, nil
}

func (mp MultiPolygon) byteSize() int {
	size := HeaderSize + Uint32Size
	for _, p := range mp {
		size += p.byteSize()
	}
	return size
}

func (mp MultiPolygon) write(buf *bytes.Buffer) {
	writeHeader(buf, GeomMultiPolygon)
	writeCount(buf, len(mp))
	for _, p := range mp {
		p.write(buf)
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
