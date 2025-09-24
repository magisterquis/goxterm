package goxterm

/*
 * errors_test.go
 * Tests for errors.go
 * By J. Stuart McMurray
 * Created 20241128
 * Last Modified 20241128
 */

import (
	"errors"
	"io"
	"testing"
)

// Make sure we can check if a interfacey error is a CtrlC.
func TestCtrlC(t *testing.T) {
	var err error = CtrlC{}
	if !errors.Is(err, io.EOF) {
		t.Errorf("Is io.EOF returned false")
	}
	if !errors.Is(err, ErrCtrlC) {
		t.Errorf("Is ErrCtrlC returned false")
	}
	if !errors.As(err, &CtrlC{}) {
		t.Errorf("As CtrlC returned false")
	}
}
