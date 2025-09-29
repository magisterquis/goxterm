package goxterm

/*
 * cooked_test.go
 * Non-raw-specific Terminal things.
 * By J. Stuart McMurray
 * Created 20250215
 * Last Modified 20250929
 */

import (
	"bytes"
	"crypto/rand"
	"errors"
	"io"
	"net"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"
)

func newCookedTerminal() (*Terminal, net.Conn) {
	lc, tc := net.Pipe()
	ss := NewTerminal(tc, "dummy")
	ss.Cooked()
	return ss, lc
}

func TestTerminal_CookedRead(t *testing.T) {
	var (
		lines = []string{"line1", "line2"}
		have  = strings.Join(lines, "\n") + "\n"
		wants = make([]struct {
			line string
			err  error
		}, len(lines)+1)
	)
	for i, line := range lines {
		wants[i].line = line
	}
	wants[len(wants)-1].err = io.EOF

	/* At some point it'd be nice to synctest these, but maybe when it's
	no longer experimental. */
	testReadLine := func(
		t *testing.T,
		readLine func(t *Terminal) func() (string, error),
	) {
		var (
			ss, c = newCookedTerminal()
			werr  error
			cerr  error
			wg    sync.WaitGroup
		)
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, werr = io.WriteString(c, have)
			cerr = c.Close()
		}()
		for _, want := range wants {
			got, gerr := readLine(ss)()
			if got != want.line {
				t.Errorf(
					"Incorrect line:\n got: %s\nwant: %s",
					got,
					want.line,
				)
			}
			if !errors.Is(gerr, want.err) {
				t.Errorf(
					"Incorrect error:\n got: %s\nwant: %s",
					gerr,
					want.err,
				)
			}
		}
		wg.Wait() /* For just in case. */
		if nil != werr {
			t.Errorf("Write error: %s", werr)
		}
		if nil != cerr {
			t.Errorf("Error closing pipe: %s", cerr)
		}
	}
	t.Run("ReadLine", func(t *testing.T) {
		testReadLine(t, func(t *Terminal) func() (string, error) {
			return t.ReadLine
		})
	})
	t.Run("ReadPassword", func(t *testing.T) {
		testReadLine(t, func(t *Terminal) func() (string, error) {
			return func() (string, error) {
				return t.ReadPassword("dummy")
			}
		})
	})
}

func TestTerminalWrite_Cooked(t *testing.T) {
	sizes := []int{1, 10, 512, 1024, 1024 * 1024}
	gens := map[string]func(n int) []byte{
		"zeros": func(n int) []byte { return make([]byte, n) },
		"same": func(n int) []byte {
			return slices.Repeat([]byte{1}, n)
		},
		"identity": func(n int) []byte {
			return slices.Collect(func(yield func(byte) bool) {
				for i := range n {
					if !yield(byte(i & 0xFF)) {
						break
					}
				}
			})
		},
		"random": func(n int) []byte {
			b := make([]byte, n)
			rand.Read(b)
			return b
		},
	}
	for _, size := range sizes {
		for name, gen := range gens {
			t.Run(name+"/"+strconv.Itoa(size), func(t *testing.T) {
				var (
					ss, c = newCookedTerminal()
					have  = gen(size)
					werr  error
					wg    sync.WaitGroup
				)
				wg.Add(1)
				go func() {
					defer wg.Done()
					_, werr = ss.Write(have)
				}()
				got := make([]byte, len(have))
				_, err := io.ReadFull(c, got)
				if nil != err {
					t.Errorf("Read error: %s", err)
				}
				wg.Wait()
				if nil != werr {
					t.Errorf("Write error: %s", err)
				}
				if !bytes.Equal(have, got) {
					t.Errorf("Read incorrect")
				}
			})
		}
	}
}

func TestTerminalSetPrompt_Cooked(t *testing.T) {
	ss, _ := newCookedTerminal()
	ss.SetPrompt("> ") /* Shoud be kinda boring. */
}

func TestCookedEscapeCodes_AllEmpty(t *testing.T) {
	ec := reflect.ValueOf(CookedEscapeCodes()).Elem()
	for i := range ec.NumField() {
		f := ec.Field(i)
		if 0 != f.Len() {
			t.Errorf(
				"CookedEscapeCodes field %d/%d (%s) not empty",
				i+1,
				ec.NumField(),
				ec.Type().Field(i).Name,
			)
		}
	}
}
