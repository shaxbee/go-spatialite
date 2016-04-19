package wkb

import "encoding/binary"

func (p *Polygon) Scan(src interface{}) error {
	b, dec, err := header(src, GeomPolygon)
	if err != nil {
		return err
	}

	_, *p, err = readPolygon(b, dec)
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
