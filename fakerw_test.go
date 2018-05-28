package fakerw

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFakeReadWriter_Active_IO(t *testing.T) {
	rw := NewFakeReadWriter(
		true,
		&IO{W: ShouldWrite([]byte{0x01}), R: Return([]byte{0x02})},
		&IO{W: ShouldWrite([]byte{0x03, 0x04}), R: Return([]byte{0x05, 0x06})},
	)

	n, err := rw.Write([]byte{0x01})
	assert.NoError(t, err)
	assert.Equal(t, 1, n)

	buf := make([]byte, 2)
	n, err = rw.Read(buf)
	assert.NoError(t, err)
	assert.Equal(t, 1, n)
	assert.Equal(t, []byte{0x02, 0x00}, buf)

	n, err = rw.Write([]byte{0x03, 0x04})
	assert.NoError(t, err)
	assert.Equal(t, 2, n)

	n, err = rw.Read(buf)
	assert.NoError(t, err)
	assert.Equal(t, 2, n)
	assert.Equal(t, []byte{0x05, 0x06}, buf)
}

func TestFakeReadWriter_Passive_IO(t *testing.T) {
	rw := NewFakeReadWriter(
		false,
		&IO{R: Return([]byte{0x01}), W: ShouldWrite([]byte{0x02})},
		&IO{R: Return([]byte{0x03, 0x04}), W: ShouldWrite([]byte{0x05, 0x06})},
	)

	buf := make([]byte, 2)
	n, err := rw.Read(buf)
	assert.NoError(t, err)
	assert.Equal(t, 1, n)
	assert.Equal(t, []byte{0x01, 0x00}, buf)

	n, err = rw.Write([]byte{0x02})
	assert.NoError(t, err)
	assert.Equal(t, 1, n)

	n, err = rw.Read(buf)
	assert.NoError(t, err)
	assert.Equal(t, 2, n)
	assert.Equal(t, []byte{0x03, 0x04}, buf)

	n, err = rw.Write([]byte{0x05, 0x06})
	assert.NoError(t, err)
	assert.Equal(t, 2, n)
}

func TestFakeReadWriter_SmallBufferReadError(t *testing.T) {
	rw := NewFakeReadWriter(
		false,
		&IO{R: Return([]byte{0x01, 0x02, 0x03}), W: ShouldWrite([]byte{0x00})},
	)

	buf := make([]byte, 1)
	n, err := rw.Read(buf)
	assert.EqualError(t, err, "read 1 bytes error @ 1:r1/w0 (buffer too small for data [01 02 03])")
	assert.Equal(t, 0, n)
}

func TestFakeReadWriter_InvalidWriteDataError(t *testing.T) {
	rw := NewFakeReadWriter(
		true,
		&IO{W: ShouldWrite([]byte{0x01, 0x02, 0x03}), R: Return([]byte{0x00})},
	)

	n, err := rw.Write([]byte{0x03})
	assert.EqualError(t, err, "write [03] error @ 1:r0/w1 (invalid write data while expected [01 02 03])")
	assert.Equal(t, 0, n)
}

func TestFakeReadWriter_ExpectedReadError(t *testing.T) {
	rw := NewFakeReadWriter(
		false,
		&IO{R: ReturnReadError(io.EOF, 3), W: ShouldWrite([]byte{0x00})},
	)

	buf := make([]byte, 8)
	n, err := rw.Read(buf)
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, 3, n)
}

func TestFakeReadWriter_ExpectedWriteError(t *testing.T) {
	rw := NewFakeReadWriter(
		true,
		&IO{W: ReturnWriteError(io.ErrShortWrite, 3), R: Return([]byte{0x00})},
	)

	n, err := rw.Write([]byte{0x03})
	assert.Equal(t, io.ErrShortWrite, err)
	assert.Equal(t, 3, n)
}

func TestFakeReadWriter_TruncatedWrite(t *testing.T) {
	rw := NewFakeReadWriter(
		true,
		&IO{W: TruncateWrite(1, ShouldWrite([]byte{0x01, 0x02, 0x03})), R: Return([]byte{0x00})},
	)

	n, err := rw.Write([]byte{0x01, 0x02, 0x03})
	assert.NoError(t, err)
	assert.Equal(t, 1, n)
}
