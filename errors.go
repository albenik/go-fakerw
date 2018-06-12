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

func NewError(m string) *IOError {
	return &IOError{message: m}
}

var ErrRepeat = errors.New("wait next IO")
