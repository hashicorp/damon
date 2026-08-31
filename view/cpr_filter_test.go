// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func readAll(t *testing.T, r io.Reader, chunkSize int) []byte {
	t.Helper()

	var out []byte
	buf := make([]byte, chunkSize)
	for {
		n, err := r.Read(buf)
		out = append(out, buf[:n]...)
		if err != nil {
			return out
		}
	}
}

func TestStripCPR(t *testing.T) {
	r := require.New(t)

	t.Run("passes plain text through unchanged", func(t *testing.T) {
		in := "ls -la\n"
		out := readAll(t, stripCPR(bytes.NewBufferString(in)), 64)
		r.Equal(in, string(out))
	})

	t.Run("strips a complete CPR sequence", func(t *testing.T) {
		in := "/ $ \x1b[40;5R\x1b[41;5Rls\n"
		out := readAll(t, stripCPR(bytes.NewBufferString(in)), 64)
		r.Equal("/ $ ls\n", string(out))
	})

	t.Run("strips a CPR sequence split across multiple reads", func(t *testing.T) {
		pr, pw := io.Pipe()
		go func() {
			pw.Write([]byte("/ $ \x1b[4"))
			pw.Write([]byte("0;5"))
			pw.Write([]byte("R"))
			pw.Write([]byte("ls\n"))
			pw.Close()
		}()
		out := readAll(t, stripCPR(pr), 1)
		r.Equal("/ $ ls\n", string(out))
	})

	t.Run("does not strip unrelated CSI sequences like arrow keys", func(t *testing.T) {
		in := "\x1b[A\x1b[B"
		out := readAll(t, stripCPR(bytes.NewBufferString(in)), 64)
		r.Equal(in, string(out))
	})

	t.Run("flushes a lone trailing ESC that never completes a sequence", func(t *testing.T) {
		in := "ls\x1b"
		out := readAll(t, stripCPR(bytes.NewBufferString(in)), 64)
		r.Equal(in, string(out))
	})

	t.Run("flushes a near-miss that turns out not to be a CPR", func(t *testing.T) {
		in := "\x1b[40;5X" // ends in X, not R
		out := readAll(t, stripCPR(bytes.NewBufferString(in)), 64)
		r.Equal(in, string(out))
	})
}
