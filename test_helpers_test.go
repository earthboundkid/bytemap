package bytemap_test

import (
	"io"
	"strings"
	"testing"

	"github.com/carlmjohnson/bytemap"
)

func naiveMap(charset string) map[byte]bool {
	m := make(map[byte]bool)
	for _, c := range []byte(charset) {
		m[c] = true
	}
	// fill out rest of the map
	// so there are entries for false as well as true
	for c := 0; c < bytemap.Len; c++ {
		m[byte(c)] = m[byte(c)]
	}
	return m
}

func naiveContains(s, charset string) bool {
	m := naiveMap(charset)
	return naiveMapContains(m, s)
}

func naiveMapContains(m map[byte]bool, s string) bool {
	for _, c := range []byte(s) {
		if !m[c] {
			return false
		}
	}
	return true
}

type ByteMapContainer interface {
	Contains(string) bool
	ContainsBytes([]byte) bool
	ContainsReader(io.Reader) (bool, error)
}

type ByteMap[T comparable] interface {
	ByteMapContainer
	Get(byte) T
}

func testContainment[M ByteMapContainer](t *testing.T, m M, s, charset string, want bool) {
	if want != m.Contains(s) {
		t.Fatalf("want: %v; s=%q charset=%q map=%v",
			want, s, charset, m)
	}
	if want != m.ContainsBytes([]byte(s)) {
		t.Fatal(want, s, charset, m)
	}
	got, err := m.ContainsReader(strings.NewReader(s))
	if err != nil {
		t.Fatal(err)
	}
	if want != got {
		t.Fatal(want, s, charset, m)
	}
}

func testGet[M ByteMap[T], T comparable](t *testing.T, m1 M, m2 map[byte]T) {
	for i := 0; i < bytemap.Len; i++ {
		if m1.Get(byte(i)) != m2[byte(i)] {
			t.Fatal(i, m1)
		}
	}
}
