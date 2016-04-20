package wkb

import "encoding/binary"

func (ls *LineString) Scan(src interface{}) error {
	b, dec, err := header(src, GeomLineString)
	if err != nil {
		return err
	}

	_, *ls, err = readPoints(b, dec)
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

func readLineString(b []byte, dec binary.ByteOrder) ([]byte, LineString, error) {
	return readPoints(b, dec)
}

func readMultiLineString(b []byte, dec binary.ByteOrder) ([]byte, MultiLineString, error) {
	b, n, err := readCount(b, dec)
	if err != nil {
		return nil, nil, err
	}

	mls := make([]LineString, n)
	for i := 0; i < n; i++ {
		b, dec, err = byteHeader(b, GeomLineString)
		if err != nil {
			return nil, nil, err
		}

		b, mls[i], err = readLineString(b, dec)
		if err != nil {
			return nil, nil, err
		}
	}
	return b, mls, err
}
