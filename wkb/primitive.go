package wkb

import (
	"bytes"
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

func writeCount(buf *bytes.Buffer, n int) {
	b := [Uint32Size]byte{}
	binary.LittleEndian.PutUint32(b[:], uint32(n))
	buf.Write(b[:])
}

func readFloat64(b []byte, dec binary.ByteOrder) ([]byte, float64) {
	return b[Float64Size:], math.Float64frombits(dec.Uint64(b))
}

func writeFloat64(buf *bytes.Buffer, f float64) {
	b := [Float64Size]byte{}
	binary.LittleEndian.PutUint64(b[:], math.Float64bits(f))
	buf.Write(b[:])
}

func header(src interface{}, tpe Kind) ([]byte, binary.ByteOrder, error) {
	b, ok := src.([]byte)
	if !ok {
		return nil, nil, ErrInvalidStorage
	}

	return byteHeader(b, tpe)
}

func writeHeader(buf *bytes.Buffer, tpe Kind) {
	b := [HeaderSize]byte{}
	b[0] = LittleEndian
	binary.LittleEndian.PutUint32(b[ByteOrderSize:], uint32(tpe))
	buf.Write(b[:])
}

func byteHeader(b []byte, tpe Kind) ([]byte, binary.ByteOrder, error) {
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
