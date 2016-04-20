package wkb

import "encoding/binary"

func (g *Geometry) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return ErrInvalidStorage
	}

	err := error(nil)
	_, *g, err = readGeometry(b)
	return err
}

func (gc *GeometryCollection) Scan(src interface{}) error {
	b, dec, err := header(src, GeomCollection)
	if err != nil {
		return err
	}

	_, *gc, err = readGeometryCollection(b, dec)
	return err
}

func readGeometry(b []byte) ([]byte, Geometry, error) {
	g := Geometry{}
	if len(b) < HeaderSize {
		return nil, g, ErrInvalidStorage
	}

	dec := byteOrder(b[0])
	if dec == nil {
		return nil, g, ErrInvalidStorage
	}

	err := error(nil)

	b, kind := readUint32(b[ByteOrderSize:], dec)

	switch kind {
	case GeomPoint:
		if len(b) < PointSize {
			return nil, g, ErrInvalidStorage
		}
		b, g.Value = readPoint(b, dec)
	case GeomLineString:
		b, g.Value, err = readLineString(b, dec)
	case GeomPolygon:
		b, g.Value, err = readPolygon(b, dec)
	case GeomMultiPoint:
		b, g.Value, err = readMultiPoint(b, dec)
	case GeomMultiLineString:
		b, g.Value, err = readMultiLineString(b, dec)
	case GeomMultiPolygon:
		b, g.Value, err = readMultiPolygon(b, dec)
	case GeomCollection:
		b, g.Value, err = readGeometryCollection(b, dec)
	default:
		err = ErrUnsupportedValue
	}

	if err != nil {
		return nil, g, err
	}

	return b, g, nil
}

func readGeometryCollection(b []byte, dec binary.ByteOrder) ([]byte, GeometryCollection, error) {
	b, n, err := readCount(b, dec)
	if err != nil {
		return nil, nil, err
	}

	gc := make([]Geometry, n)
	for i := 0; i < n; i++ {
		b, gc[i], err = readGeometry(b)
		if err != nil {
			return nil, nil, err
		}
	}

	return b, gc, nil
}
