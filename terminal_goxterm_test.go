package goxterm

/*
 * terminal_goxterm_test.go
 * Goxterm-specific tests for terminal.go
 * By J. Stuart McMurray
 * Created 20241128
 * Last Modified 20250926
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

// Make sure we don't store dupes in the history.
func TestTerminal_DuplicateHistory(t *testing.T) {
	near, far := net.Pipe()
	pt := NewTerminal(far, "")
	defer near.Close()
	defer far.Close()
	go io.Copy(io.Discard, near)

	const (
		lineOne = "one"
		lineTwo = "two"
	)
	/* Send several lines, of which all but the first are the same. */
	lines := []string{lineOne, lineTwo, lineTwo, lineTwo}
	for i, line := range lines {
		var (
			eg  errgroup.Group
			got string
		)
		eg.Go(func() error {
			if _, err := fmt.Fprintf(
				near,
				"%s\r\n",
				line,
			); nil != err {
				return fmt.Errorf(
					"sending line %q: %w",
					line,
					err,
				)
			}
			return nil
		})
		eg.Go(func() error {
			var err error
			if got, err = pt.ReadLine(); nil != err {
				return fmt.Errorf("reading line: %w", err)
			}
			return nil
		})
		if err := eg.Wait(); nil != err {
			t.Fatalf(
				"Error sending line %d/%d: %s",
				i+1,
				len(lines),
				err,
			)
		} else if want := line; got != want {
			t.Fatalf(
				"Line %d/%d incorrect\n got: %q\nwant: %q",
				i+1,
				len(lines),
				got,
				want,
			)
		}
	}

	/* Two up arrows and an enter should send the first line. */
	const up = "\x1b[A" /* up, from terminal_test.go. */
	var (
		got string
		eg  errgroup.Group
	)
	eg.Go(func() error {
		if _, err := fmt.Fprintf(near, "%s%s\r", up, up); nil != err {
			return fmt.Errorf(
				"sending up arrows and an enter: %w",
				err,
			)
		}
		return nil
	})
	eg.Go(func() error {
		var err error
		if got, err = pt.ReadLine(); nil != err {
			return fmt.Errorf(
				"reading a line after up arrows: %s",
				err,
			)
		}
		return nil
	})
	if err := eg.Wait(); nil != err {
		t.Fatalf("Error trying to up arrow and enter: %s", err)
	} else if want := lineOne; got != want {
		t.Errorf(
			"Incorrect line after two up arros and an enter\n"+
				" got: %q\n"+
				"want: %q",
			got,
			want,
		)
	}
}
