package goxterm

/*
 * rw.go
 * Turn an io.Reader and io.Writer into an io.ReadWriter
 * By J. Stuart McMurray
 * Created 20240323
 * Last Modified 20241203
 */

import (
	"io"
	"os"
)

// StdioRW combines os.Stdin and os.Stdout into a single io.ReadWriter.
var StdioRW = ReadWriter{
	Reader: os.Stdin,
	Writer: os.Stdout,
}

// ReadWriter is an io.ReadWriter which combines a separate io.Reader and
// io.Writer.
type ReadWriter struct {
	Reader io.Reader
	Writer io.Writer
}

// Read wraps os.Stdin.Read.
func (rw ReadWriter) Read(p []byte) (int, error) { return rw.Reader.Read(p) }

// Write wraps os.Stdout.Write.
func (rw ReadWriter) Write(p []byte) (int, error) { return rw.Writer.Write(p) }
