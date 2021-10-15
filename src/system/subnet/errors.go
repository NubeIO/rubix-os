package checksubnet

import "errors"

var (
	// ErrMaskOutOfBound -
	ErrMaskOutOfBound = errors.New("ErrMaskOutOfBound")
	// ErrUnparsableIPV4 -
	ErrUnparsableIPV4 = errors.New("ErrUnparsableIPV4")
)
