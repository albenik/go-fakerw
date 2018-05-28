package fakerw

import (
	"fmt"
	"io"
)

type IOFunc func([]byte) (int, error)

type IO struct {
	R io.Reader
	W io.Writer
}

type FakeReadWriter struct {
	active   bool
	sequence []*IO
	r        io.Reader
	w        io.Writer
	lvl      int
	rc       int
	wc       int
	err      error
}

func NewFakeReadWriter(a bool, s ...*IO) *FakeReadWriter {
	op, s := s[0], s[1:]
	return &FakeReadWriter{active: a, sequence: s, r: op.R, w: op.W, lvl: 1}
}

func (rw *FakeReadWriter) Read(p []byte) (int, error) {
	if rw.err != nil {
		return 0, rw.err
	}

	rw.rc++

	if rw.r == nil {
		rw.err = &IOError{message: fmt.Sprintf("unexpected read %d bytes @ %d:r%d/w%d", len(p), rw.lvl, rw.rc, rw.wc)}
		return 0, rw.err
	}

	if rw.active && rw.wc == 0 {
		rw.err = &IOError{message: fmt.Sprintf("unexpected read %d bytes before write @ %d:r%d/w%d", len(p), rw.lvl, rw.rc, rw.wc)}
		return 0, rw.err
	}

	n, err := rw.r.Read(p)

	if err != nil && err != ErrRepeat {
		if ioerr, ok := err.(*IOError); ok {
			rw.err = &IOError{message: fmt.Sprintf("read %d bytes error @ %d:r%d/w%d (%s)", len(p), rw.lvl, rw.rc, rw.wc, ioerr.message)}
			return 0, rw.err
		}
		return n, err
	}

	if err != ErrRepeat {
		rw.r = nil
		if rw.active && len(rw.sequence) > 0 {
			var op *IO
			op, rw.sequence = rw.sequence[0], rw.sequence[1:]
			rw.r, rw.w = op.R, op.W
			rw.lvl++
			rw.rc = 0
			rw.wc = 0
		}
	}

	return n, nil
}

func (rw *FakeReadWriter) Write(p []byte) (int, error) {
	if rw.err != nil {
		return 0, rw.err
	}

	rw.wc++

	if rw.w == nil {
		rw.err = &IOError{message: fmt.Sprintf("unexpected write [% X] @ %d:r%d/w%d", p, rw.lvl, rw.rc, rw.wc)}
		return 0, rw.err
	}

	if !rw.active && rw.rc == 0 {
		rw.err = &IOError{message: fmt.Sprintf("unexpected write [% X] before read @ %d:r%d/w%d", p, rw.lvl, rw.rc, rw.wc)}
		return 0, rw.err
	}

	n, err := rw.w.Write(p)

	if err != nil && err != ErrRepeat {
		if ioerr, ok := err.(*IOError); ok {
			rw.err = &IOError{message: fmt.Sprintf("write [% X] error @ %d:r%d/w%d (%s)", p, rw.lvl, rw.rc, rw.wc, ioerr.message)}
			return 0, rw.err
		}
		return n, err
	}

	if err != ErrRepeat {
		rw.w = nil
		if !rw.active && len(rw.sequence) > 0 {
			var op *IO
			op, rw.sequence = rw.sequence[0], rw.sequence[1:]
			rw.r, rw.w = op.R, op.W
			rw.lvl++
			rw.rc = 0
			rw.wc = 0
		}
	}

	return n, nil
}

func (rw *FakeReadWriter) LastError() error {
	return rw.err
}
