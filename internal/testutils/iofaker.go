package testutils

import (
	"io"
)

type ioFaker struct {
	io.ReadWriter
}

// Close is just a dummy function to implement the io.Closer interface.
func (iof ioFaker) Close() error {
	return nil
}
