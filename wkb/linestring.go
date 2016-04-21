package wkb

import (
	"bytes"
	"database/sql/driver"
)

func (ls *LineString) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return ErrInvalidStorage
	}

	_, tmp, err := ReadLineString(b)
	if err != nil {
		return err
	}

	*ls = tmp
	return nil
}

func (ls LineString) Value() (driver.Value, error) {
	buf := bytes.NewBuffer(make([]byte, 0, ls.ByteSize()))
	ls.Write(buf)
	return buf.Bytes(), nil
}

func ReadLineString(b []byte) ([]byte, LineString, error) {
	if len(b) < HeaderSize+CountSize {
		return nil, nil, ErrInvalidStorage
	}

	b, dec, err := header(b, GeomLineString)
	if err != nil {
		return nil, nil, err
	}

	b, pts, err := readPoints(b, dec)
	if err != nil {
		return nil, nil, err
	}
	return b, LineString(pts), err
}

func (ls LineString) ByteSize() int {
	return HeaderSize + Points(ls).byteSize()
}

func (ls LineString) Write(buf *bytes.Buffer) {
	writeHeader(buf, GeomLineString)
	Points(ls).write(buf)
}

func (mls *MultiLineString) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return ErrInvalidStorage
	}

	_, tmp, err := ReadMultiLineString(b)
	if err != nil {
		return err
	}

	*mls = tmp
	return nil
}

func (mls MultiLineString) Value() (driver.Value, error) {
	buf := bytes.NewBuffer(make([]byte, 0, mls.ByteSize()))
	mls.Write(buf)
	return buf.Bytes(), nil
}

func ReadMultiLineString(b []byte) ([]byte, MultiLineString, error) {
	if len(b) < HeaderSize+CountSize {
		return nil, nil, ErrInvalidStorage
	}

	b, dec, err := header(b, GeomMultiLineString)
	if err != nil {
		return nil, nil, err
	}

	b, n := readCount(b, dec)

	mls := make([]LineString, n)
	for i := 0; i < n; i++ {
		b, mls[i], err = ReadLineString(b)
		if err != nil {
			return nil, nil, err
		}
	}
	return b, mls, err
}

func (mls MultiLineString) ByteSize() int {
	size := HeaderSize + CountSize
	for _, ls := range mls {
		size += ls.ByteSize()
	}
	return size
}

func (mls MultiLineString) Write(buf *bytes.Buffer) {
	writeHeader(buf, GeomMultiLineString)
	writeCount(buf, len(mls))
	for _, ls := range mls {
		ls.Write(buf)
	}
}
