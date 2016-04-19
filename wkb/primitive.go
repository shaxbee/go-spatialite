package wkb

import (
	"encoding/binary"
	"math"
)

func readUint32(b []byte, dec binary.ByteOrder) ([]byte, uint32) {
	return b[Uint32Size:], dec.Uint32(b)
}

func readCount(b []byte, dec binary.ByteOrder) ([]byte, int, error) {
	if len(b) < Uint32Size {
		return nil, 0, ErrInvalidStorage
	}
	b, n := readUint32(b, dec)
	return b, int(n), nil
}

func readFloat64(b []byte, dec binary.ByteOrder) ([]byte, float64) {
	return b[Float64Size:], math.Float64frombits(dec.Uint64(b))
}

func header(src interface{}, tpe Kind) ([]byte, binary.ByteOrder, error) {
	b, ok := src.([]byte)
	if !ok {
		return nil, nil, ErrInvalidStorage
	}

	if len(b) < HeaderSize {
		return nil, nil, ErrInvalidStorage
	}

	dec := byteOrder(b[0])
	if dec == nil {
		return nil, nil, ErrUnsupportedValue
	}

	b, kind := readUint32(b[ByteOrderSize:], dec)
	if tpe != Kind(kind) {
		return nil, nil, ErrUnsupportedValue
	}

	return b, dec, nil
}

func byteOrder(b byte) binary.ByteOrder {
	switch b {
	case BigEndian:
		return binary.BigEndian
	case LittleEndian:
		return binary.LittleEndian
	default:
		return nil
	}
}
