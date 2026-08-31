// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

package view

import "io"

// stripCPR wraps r so that any ANSI Cursor Position Report replies
// (ESC '[' digits ';' digits 'R') are removed from the byte stream before
// they're passed on. Terminals emit these in response to a cursor-position
// query; if such a query/reply happens while a real TTY is being handed off
// to a remote exec session, the reply can otherwise be forwarded as literal
// (garbage) input to the remote shell.
func stripCPR(r io.Reader) io.Reader {
	return &cprFilter{r: r}
}

type cprFilter struct {
	r       io.Reader
	pending []byte // bytes held back while a possible CPR sequence is forming
}

func (f *cprFilter) Read(p []byte) (int, error) {
	for {
		raw := make([]byte, len(p))
		n, err := f.r.Read(raw)
		if n > 0 {
			f.pending = append(f.pending, raw[:n]...)
		}

		out, hold := extractCPR(f.pending)
		f.pending = hold

		if len(out) > 0 {
			return copy(p, out), nil
		}

		if err != nil {
			if len(f.pending) > 0 {
				// Whatever's held back wasn't a complete CPR sequence after
				// all (the underlying reader is done); flush it as-is.
				n := copy(p, f.pending)
				f.pending = f.pending[n:]
				return n, nil
			}
			return 0, err
		}
	}
}

// extractCPR scans buf for complete CPR sequences and removes them,
// returning the remaining bytes that are safe to emit now (out) and any
// trailing bytes that must be held back because they might still be the
// start of an in-progress CPR sequence (hold).
func extractCPR(buf []byte) (out, hold []byte) {
	for i := 0; i < len(buf); {
		if buf[i] != 0x1b {
			out = append(out, buf[i])
			i++
			continue
		}

		consumed, complete := matchCPR(buf[i:])
		switch {
		case complete:
			i += consumed
		case consumed == len(buf)-i:
			// The rest of the buffer is a possible in-progress match; hold
			// it all back until more data arrives.
			return out, buf[i:]
		default:
			// Not a CPR sequence: keep the ESC byte and move on.
			out = append(out, buf[i])
			i++
		}
	}
	return out, nil
}

// matchCPR checks whether buf starts with a CPR sequence: ESC '[' digits ';'
// digits 'R'.
//
// If complete is true, consumed is the length of the full matched sequence.
// If complete is false and consumed == len(buf), buf is a prefix of a
// sequence that might still complete once more data arrives, and the caller
// should hold it all back. Otherwise buf does not start with a CPR sequence
// at all, and consumed is meaningless (the caller only needs to know it's
// not a match).
func matchCPR(buf []byte) (consumed int, complete bool) {
	if len(buf) == 0 || buf[0] != 0x1b {
		return 0, false
	}
	pos := 1
	if pos == len(buf) {
		return pos, false
	}

	if buf[pos] != '[' {
		return 0, false
	}
	pos++

	digitsStart := pos
	for pos < len(buf) && buf[pos] >= '0' && buf[pos] <= '9' {
		pos++
	}
	if pos == len(buf) {
		return pos, false
	}
	if pos == digitsStart {
		return 0, false
	}

	if buf[pos] != ';' {
		return 0, false
	}
	pos++

	digitsStart = pos
	for pos < len(buf) && buf[pos] >= '0' && buf[pos] <= '9' {
		pos++
	}
	if pos == len(buf) {
		return pos, false
	}
	if pos == digitsStart {
		return 0, false
	}

	if buf[pos] != 'R' {
		return 0, false
	}
	return pos + 1, true
}
