package wkb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLineString(t *testing.T) {
	invalid := []struct {
		err error
		b   []byte
	}{
		{
			// invalid type
			ErrUnsupportedValue,
			[]byte{
				0x01, 0x42, 0x00, 0x00, 0x00,
			},
		},
		{
			// no payload
			ErrInvalidStorage,
			[]byte{
				0x01, 0x02, 0x00, 0x00, 0x00, // header
			},
		},
		{
			// no points
			ErrInvalidStorage,
			[]byte{
				0x01, 0x02, 0x00, 0x00, 0x00, // header
				0x01, 0x00, 0x00, 0x00, // numpoints - 1
			},
		},
	}

	for _, e := range invalid {
		ls := LineString{}
		if err := ls.Scan(e.b); assert.Error(t, err) {
			assert.Exactly(t, e.err, err)
		}
	}

	valid := []byte{
		0x01, 0x02, 0x00, 0x00, 0x00, // header
		0x03, 0x00, 0x00, 0x00, // numpoints - 3
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40, // point 1
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // point 2
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40, // point 3
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40,
	}

	ls := LineString{}
	if err := ls.Scan(valid); assert.NoError(t, err) {
		assert.Equal(t, LineString{{30, 10}, {10, 30}, {40, 40}}, ls)
	}
}

func TestMultiLineString(t *testing.T) {
	invalid := []struct {
		err error
		b   []byte
	}{
		{
			// invalid type
			ErrUnsupportedValue,
			[]byte{
				0x01, 0x42, 0x00, 0x00, 0x00,
			},
		},
		{
			// no payload
			ErrInvalidStorage,
			[]byte{
				0x01, 0x05, 0x00, 0x00, 0x00, // header
			},
		},
		{
			// no elements
			ErrInvalidStorage,
			[]byte{
				0x01, 0x05, 0x00, 0x00, 0x00, // header
				0x01, 0x00, 0x00, 0x00, // numlinestring - 1
			},
		},
		{
			//invalid element type
			ErrUnsupportedValue,
			[]byte{
				0x01, 0x05, 0x00, 0x00, 0x00, // header
				0x01, 0x00, 0x00, 0x00, // numlinestring - 2
				0x01, 0x42, 0x00, 0x00, 0x00, // header - invalid type
				0x00, 0x00, 0x00, 0x00, // numpoints - 0
			},
		},
		{
			// no payload in element
			ErrInvalidStorage,
			[]byte{
				0x01, 0x05, 0x00, 0x00, 0x00, // header
				0x01, 0x00, 0x00, 0x00, // numlinestring - 2
				0x01, 0x02, 0x00, 0x00, 0x00, // header - invalid type
			},
		},
	}

	for _, e := range invalid {
		mls := MultiLineString{}
		if err := mls.Scan(e.b); assert.Error(t, err) {
			assert.Exactly(t, e.err, err)
		}
	}

	valid := []byte{
		0x01, 0x05, 0x00, 0x00, 0x00, // header
		0x02, 0x00, 0x00, 0x00, // numlinestring - 2
		0x01, 0x02, 0x00, 0x00, 0x00, // linestring - 1
		0x03, 0x00, 0x00, 0x00, // numpoints - 3
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40,
		0x01, 0x02, 0x00, 0x00, 0x00, // linestring - 2
		0x04, 0x00, 0x00, 0x00, // numpoints - 4
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40,
	}

	mls := MultiLineString{}
	if err := mls.Scan(valid); assert.NoError(t, err) {
		assert.Equal(t, MultiLineString{
			LineString{{10, 10}, {20, 20}, {10, 40}},
			LineString{{40, 40}, {30, 30}, {40, 20}, {30, 10}},
		}, mls)
	}
}
