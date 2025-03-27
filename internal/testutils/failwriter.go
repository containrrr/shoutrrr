package testutils

import (
	"fmt"
	"io"
)

type failWriter struct {
	writeLimit int
	writeCount int
}

// Close is just a dummy function to implement io.Closer.
func (fw *failWriter) Close() error {
	return nil
}

// Write returns an error if the write limit has been reached.
func (fw *failWriter) Write(p []byte) (int, error) {
	fw.writeCount++
	if fw.writeCount > fw.writeLimit {
		return 0, fmt.Errorf("reached write limit %d", fw.writeLimit)
	}

	return len(p), nil
}

// CreateFailWriter returns a io.WriteCloser that returns an error after the amount of writes indicated by writeLimit.
func CreateFailWriter(writeLimit int) io.WriteCloser {
	return &failWriter{
		writeLimit: writeLimit,
	}
}
