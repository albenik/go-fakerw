package fakerw

import (
	"errors"
)

type IOError struct {
	message string
}

func (e *IOError) Error() string {
	return e.message
}

var ErrRepeat = errors.New("wait next IO")
