package wkb

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
)

func (ls *LineString) Scan(src interface{}) error {
	b, dec, err := header(src, GeomLineString)
	if err != nil {
		return err
	}

	_, *ls, err = readLineString(b, dec)
	return err
}

func (ls LineString) Value() (driver.Value, error) {
	buf := bytes.NewBuffer(make([]byte, 0, ls.byteSize()))
	ls.write(buf)
	return buf.Bytes(), nil
}

func (mls *MultiLineString) Scan(src interface{}) error {
	b, dec, err := header(src, GeomMultiLineString)
	if err != nil {
		return err
	}

	_, *mls, err = readMultiLineString(b, dec)
	return err
}

func (mls MultiLineString) Value() (driver.Value, error) {
	buf := bytes.NewBuffer(make([]byte, 0, mls.byteSize()))
	mls.write(buf)
	return buf.Bytes(), nil
}

func readLineString(b []byte, dec binary.ByteOrder) ([]byte, LineString, error) {
	b, pts, err := readPoints(b, dec)
	return b, LineString(pts), err
}

func (ls LineString) byteSize() int {
	return HeaderSize + Points(ls).byteSize()
}

func (ls LineString) write(buf *bytes.Buffer) {
	writeHeader(buf, GeomLineString)
	Points(ls).write(buf)
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

func (mls MultiLineString) byteSize() int {
	size := HeaderSize + Uint32Size
	for _, ls := range mls {
		size += ls.byteSize()
	}
	return size
}

func (mls MultiLineString) write(buf *bytes.Buffer) {
	writeHeader(buf, GeomMultiLineString)
	writeCount(buf, len(mls))
	for _, ls := range mls {
		ls.write(buf)
	}
}
