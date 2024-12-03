package goxterm

/*
 * errors.go
 * Error types
 * By J. Stuart McMurray
 * Created 20241128
 * Last Modified 20241128
 */

import (
	"errors"
	"io"
)

// ErrCtrlC is wrapped by CtrlC.
var ErrCtrlC = errors.New("Ctrl+C")

// CtrlC indicates Terminal.ReadLine read a Ctrl+C.
// Both errors.Is(CtrlC{}, io.EOF) and errors.Is(CtrlC{}, ErrCtrlC) return
// true, as does errors.As(CtrlC{}, &CtrlC{}).
type CtrlC struct{}

// Error implements the error interface.
func (err CtrlC) Error() string { return ErrCtrlC.Error() }

// Unwrap returns []error{io.EOF, ErrCtrlC}.  This makes it easier to use
// errors.Is with both io.EOF and ErrCtrlC.
func (err CtrlC) Unwrap() []error { return []error{io.EOF, ErrCtrlC} }
