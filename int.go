package bytemap

import (
	"io"

	"golang.org/x/exp/slices"
)

// Int is an array backed map from byte to integer.
type Int [Size]int

var _ io.Writer = (*Int)(nil)

// Write satisfies io.Writer.
func (m *Int) Write(p []byte) (int, error) {
	for _, c := range p {
		m[c]++
	}
	return len(p), nil
}

var _ io.StringWriter = (*Int)(nil)

// WriteString satisfies io.StringWriter.
func (m *Int) WriteString(s string) (n int, err error) {
	for _, c := range []byte(s) {
		m[c]++
	}
	return len(s), nil
}

// Contains reports whether all bytes in s are already in m.
func (m *Int) Contains(s string) bool {
	for _, b := range []byte(s) {
		if m[b] < 1 {
			return false
		}
	}
	return true
}

// ContainsBytes reports whether all bytes in b are already in m.
func (m *Int) ContainsBytes(b []byte) bool {
	for _, c := range b {
		if m[c] < 1 {
			return false
		}
	}
	return true
}

// ContainsReader reports whether all bytes in r are already in m.
// If the reader fails, it returns false, error.
// If it reads to io.EOF, it returns true, nil.
func (m *Int) ContainsReader(r io.Reader) (bool, error) {
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

// ToMap makes a map[byte]int from the bytemap.
func (m *Int) ToMap() map[byte]int {
	m2 := make(map[byte]int)
	for i := range m {
		m2[byte(i)] = m[i]
	}
	return m2
}

// ToBool makes a map[byte]bool from the bytemap.
func (m *Int) ToBool() map[byte]bool {
	m2 := make(map[byte]bool)
	for i := range m {
		m2[byte(i)] = m[i] > 0
	}
	return m2
}

// Equals reports if two Ints are equal.
func (m *Int) Equals(other *Int) bool {
	return *m == *other
}

// Set sets one byte in the Int byte map.
func (m *Int) Set(key byte, value int) {
	m[key] = value
}

// Get looks up one byte in the Int byte map.
func (m *Int) Get(key byte) int {
	return m[key]
}

type IntN struct {
	Byte byte
	N    int
}

// MostCommon returns a slice of character counts for m
// from highest count to lowest.
func (m *Int) MostCommon() []IntN {
	freqs := make([]IntN, len(m))
	for i, n := range m {
		freqs[i] = IntN{byte(i), n}
	}
	slices.SortStableFunc(freqs, func(a, b IntN) bool {
		return a.N > b.N
	})
	return freqs
}

// Clone copies m.
func (m *Int) Clone() *Int {
	m2 := *m
	return &m2
}
