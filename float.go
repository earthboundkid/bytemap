package bytemap

import (
	"cmp"
	"io"
	"slices"
)

// Float is an array backed map from byte to float64.
type Float [Len]float64

var _ io.Writer = (*Float)(nil)

// Write satisfies io.Writer.
func (m *Float) Write(p []byte) (int, error) {
	for _, c := range p {
		m[c]++
	}
	return len(p), nil
}

var _ io.StringWriter = (*Float)(nil)

// WriteString satisfies io.StringWriter.
func (m *Float) WriteString(s string) (n int, err error) {
	for _, c := range []byte(s) {
		m[c]++
	}
	return len(s), nil
}

// Contains reports whether all bytes in s are already in m.
func (m *Float) Contains(s string) bool {
	for _, b := range []byte(s) {
		if m[b] < 1 {
			return false
		}
	}
	return true
}

// ContainsBytes reports whether all bytes in b are already in m.
func (m *Float) ContainsBytes(b []byte) bool {
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
func (m *Float) ContainsReader(r io.Reader) (bool, error) {
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

// ToMap makes a map[byte]float64 from the bytemap.
func (m *Float) ToMap() map[byte]float64 {
	m2 := make(map[byte]float64)
	for i := range m {
		m2[byte(i)] = m[i]
	}
	return m2
}

// ToBool makes a Bool from the bytemap.
func (m *Float) ToBool() *Bool {
	var m2 Bool
	for i := range m {
		m2[byte(i)] = m[i] > 0
	}
	return &m2
}

// Equals reports if two Floats are equal.
func (m *Float) Equals(other *Float) bool {
	return *m == *other
}

// Set sets one byte in the Float byte map.
func (m *Float) Set(key byte, value float64) {
	m[key] = value
}

// Get looks up one byte in the Float byte map.
func (m *Float) Get(key byte) float64 {
	return m[key]
}

type FloatN struct {
	Byte byte
	N    float64
}

// SetFrequencies sets each value in m to its overall frequency from 0 to 1.
func (m *Float) SetFrequencies() {
	sum := float64(0)
	for i := range m {
		sum += m[i]
	}
	for i, n := range m {
		m[i] = n / sum
	}
}

// MostCommon returns a slice of character counts for m
// from highest count to lowest.
func (m *Float) MostCommon() []FloatN {
	freqs := make([]FloatN, len(m))
	for i, n := range m {
		freqs[i] = FloatN{byte(i), n}
	}
	slices.SortStableFunc(freqs, func(a, b FloatN) int {
		return cmp.Compare(b.N, a.N)
	})
	return freqs
}

// Clone copies m.
func (m *Float) Clone() *Float {
	m2 := *m
	return &m2
}
