package bytemap

import (
	"io"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Bool is an array backed map from byte to bool.
type Bool [Len]bool

// Make initializes a bytemap.Bool with a byte sequence.
func Make[byteseq []byte | string](seq byteseq) *Bool {
	var m Bool
	for _, c := range []byte(seq) {
		m[c] = true
	}
	return &m
}

var _ io.Writer = (*Bool)(nil)

// Write satisfies io.Writer.
func (m *Bool) Write(p []byte) (int, error) {
	for _, c := range p {
		m[c] = true
	}
	return len(p), nil
}

var _ io.StringWriter = (*Bool)(nil)

// WriteString satisfies io.StringWriter.
func (m *Bool) WriteString(s string) (n int, err error) {
	for _, c := range []byte(s) {
		m[c] = true
	}
	return len(s), nil
}

// Contains reports whether all bytes in s are already in m.
func (m *Bool) Contains(s string) bool {
	for _, b := range []byte(s) {
		if !m[b] {
			return false
		}
	}
	return true
}

// ContainsBytes reports whether all bytes in b are already in m.
func (m *Bool) ContainsBytes(b []byte) bool {
	for _, c := range b {
		if !m[c] {
			return false
		}
	}
	return true
}

// ContainsReader reports whether all bytes in r are already in m.
// If the reader fails, it returns false, error.
// If it reads to io.EOF, it returns true, nil.
func (m *Bool) ContainsReader(r io.Reader) (bool, error) {
	var buf [4096]byte
	for {
		n, err := r.Read(buf[:])
		if err != nil && err != io.EOF {
			return false, err
		}
		if !m.ContainsBytes(buf[:n]) {
			return false, nil
		}
		if err == io.EOF {
			return true, nil
		}
	}
}

// ToMap makes a map[byte]bool from the bytemap.
func (m *Bool) ToMap() map[byte]bool {
	m2 := make(map[byte]bool)
	for i := range m {
		m2[byte(i)] = m[i]
	}
	return m2
}

// Equals reports if two Bools are equal.
func (m *Bool) Equals(other *Bool) bool {
	return *m == *other
}

// Set sets one byte in the Bool byte map.
func (m *Bool) Set(key byte, value bool) {
	m[key] = value
}

// Get looks up one byte in the Bool byte map.
func (m *Bool) Get(key byte) bool {
	return m[key]
}

// Clone copies m.
func (m *Bool) Clone() *Bool {
	m2 := *m
	return &m2
}

func (m *Bool) String() string {
	var buf strings.Builder
	buf.Grow(Len + len("Bool()"))
	buf.WriteString("Bool(")
	for i := 0; i < Len; i++ {
		if !m[i] {
			buf.WriteString("_")
			continue
		}
		c := byte(i)
		b := []byte{c}
		if utf8.Valid(b) && unicode.IsPrint(rune(c)) && c != '_' {
			buf.Write(b)
		} else {
			buf.WriteString(".")
		}
	}
	buf.WriteString(")")
	return buf.String()
}
