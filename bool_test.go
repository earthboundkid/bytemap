package bytemap_test

import (
	"io"
	"strings"
	"testing"

	"github.com/carlmjohnson/bytemap"
	"golang.org/x/exp/maps"
)

func naiveMap(charset string) map[byte]bool {
	m := make(map[byte]bool)
	for _, c := range []byte(charset) {
		m[c] = true
	}
	// fill out rest of the map
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

func TestMakeBool(t *testing.T) {
	for _, tc := range []struct {
		s, charset string
		want       bool
	}{
		{"", "", true},
		{"a", "a", true},
		{"ab", "a", false},
		{"the quick brown fox jumps over the lazy dog",
			"abcdefghijklmnopqrstuvwxyz ", true},
		{"the quick brown fox jumps over the lazy dog.",
			"abcdefghijklmnopqrstuvwxyz ", false},
	} {
		m := bytemap.Make(tc.charset)
		testContainment(t, m, tc.s, tc.charset, tc.want)
	}
}

func FuzzMakeBool(f *testing.F) {
	f.Add("", "")
	f.Add("a", "a")
	f.Add("a", "b")
	f.Add("ab", "ab")
	f.Add("ab", "abc")
	for i := 0; i < 1_000_000; i = (i + 1) * 2 {
		for j := 0; j < 1_000_000; j = (j + 1) * 2 {
			s := strings.Repeat("a", i)
			charset := strings.Repeat("a", j)
			f.Add(s, charset)
			f.Add(s+"b", charset)
		}
	}
	f.Fuzz(func(t *testing.T, s, charset string) {
		want := naiveContains(s, charset)
		// Test Make
		{
			m := bytemap.Make(charset)
			testContainment(t, m, s, charset, want)
		}
		// Test WriteString
		{
			m := &bytemap.Bool{}
			n, err := m.WriteString(charset)
			if err != nil {
				t.Fatal(err)
			}
			if n != len(charset) {
				t.Fatal(len(charset))
			}
			testContainment(t, m, s, charset, want)
		}
		// Test Write
		{
			m := &bytemap.Bool{}
			n, err := m.Write([]byte(charset))
			if err != nil {
				t.Fatal(err)
			}
			if n != len(charset) {
				t.Fatal(len(charset))
			}
			testContainment(t, m, s, charset, want)
		}
		// Test io.Copy
		{
			m := &bytemap.Bool{}
			n64, err := io.Copy(m, strings.NewReader(charset))
			if err != nil {
				t.Fatal(err)
			}
			if n64 != int64(len(charset)) {
				t.Fatal(len(charset))
			}
			testContainment(t, m, s, charset, want)
		}
	})
}

func testContainment(t *testing.T, m *bytemap.Bool, s, charset string, want bool) {
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

func FuzzBoolToMap(f *testing.F) {
	f.Add("", "")
	f.Add("a", "b")
	f.Add(
		"the quick brown fox jumps over a lazy dog.",
		"abcdefghijklmnopqrstuvwxyz. ",
	)
	f.Fuzz(func(t *testing.T, a, b string) {
		aNaive := naiveMap(a)
		aMap := bytemap.Make(a)
		if !maps.Equal(aNaive, aMap.ToMap()) {
			t.Fatal(a, aMap)
		}
		testBoolGet(t, aMap, aNaive)
		bNaive := naiveMap(b)
		bMap := bytemap.Make(b)
		if !maps.Equal(bNaive, bMap.ToMap()) {
			t.Fatal(b, bMap)
		}
		testBoolGet(t, bMap, bNaive)
		if maps.Equal(aNaive, bNaive) != aMap.Equals(bMap) {
			t.Fatal(aMap, bMap)
		}
	})
}

func testBoolGet(t *testing.T, m1 *bytemap.Bool, m2 map[byte]bool) {
	for i := 0; i < bytemap.Len; i++ {
		if m1.Get(byte(i)) != m2[byte(i)] {
			t.Fatal(i, m1)
		}
	}
}
func FuzzBoolSet(f *testing.F) {
	f.Add("", "", "")
	f.Add("a", "a", "a")
	f.Add("abc", "bcde", "b")
	f.Fuzz(func(t *testing.T, add, remove, restore string) {
		var bf bytemap.Bool
		m := make(map[byte]bool)
		for _, c := range []byte(add) {
			bf.Set(c, true)
			m[c] = true
		}
		for _, c := range []byte(remove) {
			bf.Set(c, false)
			m[c] = false
		}
		for _, c := range []byte(restore) {
			bf.Set(c, true)
			m[c] = true
		}
		// Fill in blanks
		for i := 0; i < bytemap.Len; i++ {
			m[byte(i)] = m[byte(i)]
		}
		if !maps.Equal(bf.ToMap(), m) {
			t.Fatal(bf)
		}
	})
}
