package wkb

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeader(t *testing.T) {
	invalid := map[error][]byte{
		ErrUnsupportedValue: {0x02, 0x01, 0x00, 0x00, 0x00},
		ErrUnsupportedValue: {0x01, 0x02, 0x00, 0x00, 0x00},
	}
	for expected, b := range invalid {
		if _, _, err := header(b, GeomPoint); assert.Error(t, err) {
			assert.Exactly(t, expected, err)
		}
	}

	valid := []byte{0x01, 0x01, 0x00, 0x00, 0x00}
	if b, bo, err := header(valid, GeomPoint); assert.NoError(t, err) {
		assert.Len(t, b, 0)
		assert.Exactly(t, binary.LittleEndian, bo)
	}
}

func TestByteOrder(t *testing.T) {
	assert.Exactly(t, binary.BigEndian, byteOrder(0x00))
	assert.Exactly(t, binary.LittleEndian, byteOrder(0x01))
	assert.Nil(t, byteOrder(0x42))
}
