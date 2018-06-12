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
	opnum    int
	rc       int
	wc       int
	err      error
}

func NewFakeReadWriter(a bool, s ...*IO) *FakeReadWriter {
	return &FakeReadWriter{active: a, sequence: s, opnum: 1}
}

func (rw *FakeReadWriter) Read(p []byte) (int, error) {
	if rw.err != nil {
		return 0, rw.err
	}

	if rw.opnum > len(rw.sequence) {
		rw.err = &IOError{message: fmt.Sprintf("unexpected read %d bytes @ %d:r%d/w%d", len(p), rw.opnum, rw.rc, rw.wc)}
		return 0, rw.err
	}
	op := rw.sequence[rw.opnum-1]

	rw.rc++

	if op.R == nil {
		rw.err = &IOError{message: fmt.Sprintf("unexpected read %d bytes @ %d:r%d/w%d", len(p), rw.opnum, rw.rc, rw.wc)}
		return 0, rw.err
	}

	if rw.active && rw.wc == 0 && op.W != nil {
		rw.err = &IOError{message: fmt.Sprintf("unexpected read %d bytes before write @ %d:r%d/w%d", len(p), rw.opnum, rw.rc, rw.wc)}
		return 0, rw.err
	}

	n, err := op.R.Read(p)

	if err != nil && err != ErrRepeat {
		if ioerr, ok := err.(*IOError); ok {
			rw.err = &IOError{message: fmt.Sprintf("read %d bytes error @ %d:r%d/w%d (%s)", len(p), rw.opnum, rw.rc, rw.wc, ioerr.message)}
			return 0, rw.err
		}
		return n, err
	}

	if err != ErrRepeat {
		if rw.active || op.W == nil {
			rw.opnum++
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

	if rw.opnum > len(rw.sequence) {
		rw.err = &IOError{message: fmt.Sprintf("unexpected write [% X] @ %d:r%d/w%d", p, rw.opnum, rw.rc, rw.wc)}
		return 0, rw.err
	}
	op := rw.sequence[rw.opnum-1]

	rw.wc++

	if op.W == nil {
		rw.err = &IOError{message: fmt.Sprintf("unexpected write [% X] @ %d:r%d/w%d", p, rw.opnum, rw.rc, rw.wc)}
		return 0, rw.err
	}

	if !rw.active && rw.rc == 0 && op.R != nil {
		rw.err = &IOError{message: fmt.Sprintf("unexpected write [% X] before read @ %d:r%d/w%d", p, rw.opnum, rw.rc, rw.wc)}
		return 0, rw.err
	}

	n, err := op.W.Write(p)

	if err != nil && err != ErrRepeat {
		if ioerr, ok := err.(*IOError); ok {
			rw.err = &IOError{message: fmt.Sprintf("write [% X] error @ %d:r%d/w%d (%s)", p, rw.opnum, rw.rc, rw.wc, ioerr.message)}
			return 0, rw.err
		}
		return n, err
	}

	if err != ErrRepeat {
		if !rw.active || op.R == nil {
			rw.opnum++
			rw.rc = 0
			rw.wc = 0
		}
	}

	return n, nil
}

func (rw *FakeReadWriter) LastError() error {
	return rw.err
}
