package bytemap

import "io"

// BitField is a map from byte to bool backed by a bit field.
// It is not as fast as Bool,
// but if memory is an important consideration,
// it is 8 times smaller.
type BitField [Size / 8]byte

func MakeBitField[byteseq []byte | string](s byteseq) *BitField {
	var m BitField
	for _, c := range []byte(s) {
		i := c / 8
		mask := byte(1 << (c % 8))
		m[i] |= mask
	}
	return &m
}

// Write satisfies io.Writer.
func (m *BitField) Write(p []byte) (int, error) {
	for _, c := range p {
		i := c / 8
		mask := byte(1 << (c % 8))
		m[i] |= mask
	}
	return len(p), nil
}

// WriteString satisfies io.StringWriter.
func (m *BitField) WriteString(s string) (n int, err error) {
	for _, c := range []byte(s) {
		i := c / 8
		mask := byte(1 << (c % 8))
		m[i] |= mask
	}
	return len(s), nil
}

// Contains reports whether all bytes in s are already in m.
func (m *BitField) Contains(s string) bool {
	for _, c := range []byte(s) {
		i := c / 8
		mask := byte(1 << (c % 8))
		if m[i]&mask == 0 {
			return false
		}
	}
	return true
}

// ContainsBytes reports whether all bytes in b are already in m.
func (m *BitField) ContainsBytes(b []byte) bool {
	for _, c := range b {
		mask := byte(1 << (c % 8))
		masked := m[c/8] & mask
		if contains := masked != 0; !contains {
			return false
		}
	}
	return true
}

// ContainsReader reports whether all bytes in r are already in m.
// If the reader fails, it returns false, error.
// If it reads to io.EOF, it returns true, nil.
func (m *BitField) ContainsReader(r io.Reader) (bool, error) {
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
func (m *BitField) ToMap() map[byte]bool {
	m2 := make(map[byte]bool)
	for c := 0; c < Size; c++ {
		mask := byte(1 << (c % 8))
		masked := m[c/8] & mask
		m2[byte(c)] = masked != 0
	}
	return m2
}

// Equals reports if two BitFields are equal.
func (m *BitField) Equals(other *BitField) bool {
	return *m == *other
}

// Set sets one byte in the Bitfield byte map
func (m *BitField) Set(key byte, value bool) {
	i := key / 8
	if value {
		mask := byte(1 << (key % 8))
		m[i] |= mask
	} else {
		mask := ^byte(1 << (key % 8))
		m[i] &= mask
	}
}

// Get looks up one byte in the BitField byte map.
func (m *BitField) Get(key byte) bool {
	mask := byte(1 << (key % 8))
	masked := m[key/8] & mask
	return masked != 0
}
