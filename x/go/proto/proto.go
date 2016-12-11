package proto

import (
	"errors"
)

const Unused = 0

var (
	ErrNegativeLength = errors.New("negative length")
)

type Message interface {
	MessageName() string
	Encode(Writer) error
	Decode([]byte) (int, error)
}
