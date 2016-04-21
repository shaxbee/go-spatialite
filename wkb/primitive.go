package wkb

import (
	"bytes"
	"encoding/binary"
	"math"
)

func readUint32(b []byte, dec binary.ByteOrder) ([]byte, uint32) {
	return b[CountSize:], dec.Uint32(b)
}

func readCount(b []byte, dec binary.ByteOrder) ([]byte, int) {
	b, n := readUint32(b, dec)
	return b, int(n)
}

func writeCount(buf *bytes.Buffer, n int) {
	b := [CountSize]byte{}
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

func header(b []byte, tpe Kind) ([]byte, binary.ByteOrder, error) {
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

func writeHeader(buf *bytes.Buffer, tpe Kind) {
	b := [HeaderSize]byte{}
	b[0] = LittleEndian
	binary.LittleEndian.PutUint32(b[ByteOrderSize:], uint32(tpe))
	buf.Write(b[:])
}
