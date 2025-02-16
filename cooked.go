package goxterm

/*
 * cooked.go
 * Non-raw-specific Terminal things.
 * By J. Stuart McMurray
 * Created 20250215
 * Last Modified 20250215
 */

import (
	"bufio"
	"cmp"
	"io"
)

// cookedEscapeCodes are used when in cooked mode to prevent colors from
// gunking up simple stream output.
var cookedEscapeCodes = EscapeCodes{
	Black:   []byte{},
	Red:     []byte{},
	Green:   []byte{},
	Yellow:  []byte{},
	Blue:    []byte{},
	Magenta: []byte{},
	Cyan:    []byte{},
	White:   []byte{},

	Reset: []byte{},
}

// Cooked changes t to behave better when the underlying io.ReadWriter isn't a
// terminal in Raw mode.
//
// The specific changes are:
//   - No prompts are printed
//   - t.ReadLine performs no special treatment of control characters
//   - t.ReadPassword is equivalent to t.ReadLine
//   - Bracketed paste mode is ignored
//   - Changing the terminal's size has no effect
//   - t.Escape contains empty slices
//   - t.Write is a thin wrapper around the underlying [io.ReadWriter]'s Write
//
// Cooked should be called before calling any of t's other methods or accessing
// any of t's fields.
func (t *Terminal) Cooked() {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.Escape = &cookedEscapeCodes
	t.cooked = true
}

// cookedReadLine reads a line from t.scanner.
// The caller must hold t.lock.
func (t *Terminal) cookedReadLine() (string, error) {
	/* Make sure we have a scanner.  We only really need to lock here. */
	if nil == t.scanner {
		t.scanner = bufio.NewScanner(t.c)
	}

	/* Unlock here to unblock other methods on t, but relock before we
	return so the caller's unlock doesn't panic anything. */
	t.lock.Unlock()
	defer t.lock.Lock()

	/* If we've no more lines, let the user know.  Since scanner
	won't return EOF, we'll have to deal with it ourselves. */
	if !t.scanner.Scan() {
		return "", cmp.Or(t.scanner.Err(), io.EOF)
	}

	return t.scanner.Text(), nil
}
