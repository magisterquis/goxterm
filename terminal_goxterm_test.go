package goxterm

/*
 * terminal_goxterm_test.go
 * Goxterm-specific tests for terminal.go
 * By J. Stuart McMurray
 * Created 20241128
 * Last Modified 20241128
 */

import (
	"errors"
	"fmt"
	"io"
	"net"
	"testing"

	"golang.org/x/sync/errgroup"
)

// Make sure ReadLine gives us a Ctrl+C-specific error on Ctrl+C.
func TestTerminalReadLine_CtrlC(t *testing.T) {
	/* Terminal wrapping a pretend pty. */
	near, far := net.Pipe()
	pt := NewTerminal(far, "")

	/* Send a raw mode Ctrl+c. */
	var eg errgroup.Group
	eg.Go(func() error {
		_, err := near.Write([]byte{keyCtrlC})
		if nil != err {
			err = fmt.Errorf("sending Ctrl+C: %w", err)
		}
		return err
	})
	/* Discard from the terminal. */
	eg.Go(func() error {
		_, err := io.Copy(io.Discard, near)
		if errors.Is(err, io.ErrClosedPipe) { /* Expected */
			err = nil
		} else if nil != err {
			err = fmt.Errorf("reading from terminal: %w", err)
		}
		return err
	})

	/* See if we get the right sort of error. */
	l, err := pt.ReadLine()
	if "" != l {
		t.Errorf("Unexpected line: %q", l)
	}
	if !errors.As(err, &CtrlC{}) {
		t.Errorf("Incorrect %T error: %s", err, err)
	}

	/* Make sure we actually sent Ctrl+C. */
	if err := near.Close(); nil != err {
		t.Errorf("Error closing pipe: %s", err)
	}
	if err := eg.Wait(); nil != err {
		t.Errorf("Additional terminal I/O error: %s", err)
	}
}
