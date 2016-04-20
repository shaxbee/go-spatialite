package wkb

import "encoding/binary"

type Point struct {
	X, Y float64
}

func (p *Point) Scan(src interface{}) error {
	b, dec, err := header(src, GeomPoint)
	if err != nil {
		return err
	}

	_, *p, err = readPoint(b, dec)
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

func readPoint(b []byte, dec binary.ByteOrder) ([]byte, Point, error) {
	p := Point{}
	if len(b) < PointSize {
		return nil, p, ErrInvalidStorage
	}

	b, p.X = readFloat64(b, dec)
	b, p.Y = readFloat64(b, dec)
	return b, p, nil
}

func readMultiPoint(b []byte, dec binary.ByteOrder) ([]byte, MultiPoint, error) {
	return readPoints(b, dec)
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
