package wkb

import (
	"bytes"
	"database/sql/driver"
)

func New(b []byte) (Geometry, error) {
	_, g, err := ReadGeometry(b)
	return g, err
}

func ReadGeometry(b []byte) ([]byte, Geometry, error) {
	if len(b) < HeaderSize {
		return nil, nil, ErrInvalidStorage
	}

	dec := byteOrder(b[0])
	if dec == nil {
		return nil, nil, ErrInvalidStorage
	}

	_, kind := readUint32(b[ByteOrderSize:], dec)

	var g Geometry
	var err error
	switch kind {
	case GeomPoint:
		b, g, err = ReadPoint(b)
	case GeomLineString:
		b, g, err = ReadLineString(b)
	case GeomPolygon:
		b, g, err = ReadPolygon(b)
	case GeomMultiPoint:
		b, g, err = ReadMultiPoint(b)
	case GeomMultiLineString:
		b, g, err = ReadMultiLineString(b)
	case GeomMultiPolygon:
		b, g, err = ReadMultiPolygon(b)
	case GeomCollection:
		b, g, err = ReadGeometryCollection(b)
	default:
		return nil, nil, ErrUnsupportedValue
	}

	return b, g, err
}

func (gc *GeometryCollection) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return ErrInvalidStorage
	}

	_, tmp, err := ReadGeometryCollection(b)
	if err != nil {
		return err
	}

	*gc = tmp
	return err
}

func (gc GeometryCollection) Value() (driver.Value, error) {
	buf := bytes.NewBuffer(make([]byte, 0, gc.ByteSize()))
	gc.Write(buf)
	return buf.Bytes(), nil
}

func ReadGeometryCollection(b []byte) ([]byte, GeometryCollection, error) {
	if len(b) < HeaderSize+CountSize {
		return nil, nil, ErrInvalidStorage
	}

	b, dec, err := header(b, GeomCollection)
	if err != nil {
		return nil, nil, err
	}

	b, n := readCount(b, dec)

	gc := make([]Geometry, n)
	for i := 0; i < n; i++ {
		b, gc[i], err = ReadGeometry(b)
		if err != nil {
			return nil, nil, err
		}
	}

	return b, gc, nil
}

func (gc GeometryCollection) ByteSize() int {
	size := HeaderSize + CountSize
	for _, g := range gc {
		size += g.ByteSize()
	}
	return size
}

func (gc GeometryCollection) Write(buf *bytes.Buffer) {
	writeHeader(buf, GeomCollection)
	writeCount(buf, len(gc))
	for _, g := range gc {
		g.Write(buf)
	}
}
