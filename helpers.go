package fakerw

import (
	"bytes"
	"fmt"
	"io"
	"time"
)

type ReaderFunc func([]byte) (int, error)

func (fn ReaderFunc) Read(p []byte) (int, error) {
	return fn(p)
}

type WriterFunc func([]byte) (int, error)

func (fn WriterFunc) Write(p []byte) (int, error) {
	return fn(p)
}

func Return(data []byte) io.Reader {
	return ReaderFunc(func(p []byte) (int, error) {
		n := copy(p, data)
		if n < len(data) {
			return 0, &IOError{message: fmt.Sprintf("buffer too small for data [% X]", data)}
		}
		return n, nil
	})
}

func BufferSizeLimits(min, max int, r io.Reader) io.Reader {
	return ReaderFunc(func(p []byte) (int, error) {
		buflen := len(p)
		if buflen < min || buflen > max {
			return 0, &IOError{message: fmt.Sprintf("invalid read buffer length %d while expected between [%d, %d]", buflen, min, max)}
		}
		return r.Read(p)
	})
}

func DelayRead(d time.Duration, r io.Reader) io.Reader {
	return ReaderFunc(func(p []byte) (int, error) {
		time.Sleep(d)
		return r.Read(p)
	})
}

func ReturnReadError(err error, n int) io.Reader {
	return ReaderFunc(func(p []byte) (int, error) {
		return n, err
	})
}

func ShouldWrite(data []byte) io.Writer {
	return WriterFunc(func(p []byte) (int, error) {
		if !bytes.Equal(p, data) {
			return 0, &IOError{message: fmt.Sprintf("invalid write data while expected [% X]", data)}
		}
		return len(p), nil
	})
}

func ShouldWriteIn(t time.Duration, w io.Writer) io.Writer {
	return WriterFunc(func(p []byte) (int, error) {
		start := time.Now()
		n, err := w.Write(p)
		stop := time.Now()

		d := stop.Sub(start)
		if d > t {
			return 0, &IOError{message: fmt.Sprintf("write done in %s while expected in %s", d, t)}
		}
		return n, err
	})
}

func DelayWrite(d time.Duration, w io.Writer) io.Writer {
	return WriterFunc(func(p []byte) (int, error) {
		time.Sleep(d)
		return w.Write(p)
	})
}

func TruncateWrite(n int, w io.Writer) io.Writer {
	return WriterFunc(func(p []byte) (int, error) {
		_, err := w.Write(p)
		return n, err
	})
}

func ReturnWriteError(err error, n int) io.Writer {
	return WriterFunc(func(p []byte) (int, error) {
		return n, err
	})
}
