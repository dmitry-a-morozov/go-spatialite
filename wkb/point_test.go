package wkb

import (
	"encoding/binary"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPoint(t *testing.T) {
	invalid := map[error][]byte{
		ErrInvalidStorage: {
			0x01,
		}, // header too short
		ErrInvalidStorage: {
			0x01, 0x01, 0x00, 0x00, 0x00, 0x00,
		}, // no payload
		ErrInvalidStorage: {
			0x02, 0x01, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40,
		}, // invalid endianness
		ErrInvalidStorage: {
			0x01, 0x01, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40,
		}, // single coordinate only
		ErrUnsupportedValue: {
			0x01, 0x02, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40,
		}, // invalid type
	}

	for expected, b := range invalid {
		p := Point{}
		if err := p.Scan(b); assert.Error(t, err) {
			assert.Exactly(t, expected, err, "Expected point <%s> to fail", hex.EncodeToString(b))
		}
	}

	valid := []byte{
		0x01, 0x01, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40,
	}
	p := Point{}
	if assert.NoError(t, p.Scan(valid)) {
		assert.Equal(t, Point{30, 10}, p)
	}
}

func TestMultiPoint(t *testing.T) {
	invalid := []struct {
		err error
		b   []byte
	}{
		{
			// invalid byte order
			ErrUnsupportedValue,
			[]byte{
				0x42, 0x04, 0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			// invalid type
			ErrUnsupportedValue, []byte{
				0x01, 0x42, 0x00, 0x00, 0x00,
			},
		},
		{
			// no payload
			ErrInvalidStorage, []byte{
				0x01, 0x04, 0x00, 0x00, 0x00,
			},
		},
		{
			// no points
			ErrInvalidStorage, []byte{
				0x01, 0x04, 0x00, 0x00, 0x00,
				0x01, 0x00, 0x00, 0x00, // numpoints - 1
			},
		},
		{
			// incomplete point
			ErrInvalidStorage, []byte{
				0x01, 0x04, 0x00, 0x00, 0x00,
				0x01, 0x00, 0x00, 0x00, // numpoints - 1
				0x01, 0x01, 0x00, 0x00, 0x00, // point without payload
			},
		},
		{
			// element not a point
			ErrUnsupportedValue, []byte{
				0x01, 0x04, 0x00, 0x00, 0x00,
				0x01, 0x00, 0x00, 0x00, // numpoints - 1
				0x01, 0x02, 0x00, 0x00, 0x00, // invalid element
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40,
			},
		},
	}

	for _, e := range invalid {
		mp := MultiPoint{}
		if err := mp.Scan(e.b); assert.Error(t, err) {
			assert.Exactly(t, e.err, err, "Expected multipoint <%s> to fail", hex.EncodeToString(e.b))
		}
	}

	valid := []byte{
		0x01, 0x04, 0x00, 0x00, 0x00, // header
		0x04, 0x00, 0x00, 0x00, // numpoints - 4
		0x01, 0x01, 0x00, 0x00, 0x00, // point 1
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40,
		0x01, 0x01, 0x00, 0x00, 0x00, // point 2
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40,
		0x01, 0x01, 0x00, 0x00, 0x00, // point 3
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40,
		0x01, 0x01, 0x00, 0x00, 0x00, // point 4
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40,
	}

	mp := MultiPoint{}
	if assert.NoError(t, mp.Scan(valid)) {
		assert.Equal(t, MultiPoint{{10, 40}, {40, 30}, {20, 20}, {30, 10}}, mp)
	}
}

func TestReadPoints(t *testing.T) {
	invalid := [][]byte{
		{0x01, 0x00, 0x00},       // numpoints too short
		{0x01, 0x00, 0x00, 0x00}, // no payload
	}

	for _, b := range invalid {
		_, _, err := readPoints(b, binary.LittleEndian)
		if assert.Error(t, err) {
			assert.Exactly(t, ErrInvalidStorage, err)
		}
	}
}

func TestEqual(t *testing.T) {
	assert.True(t, Point{10, 10}.Equal(Point{10, 10}))
	assert.False(t, Point{10, 20}.Equal(Point{20, 10}))
}
