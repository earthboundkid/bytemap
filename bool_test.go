package bytemap_test

import (
	"io"
	"strings"
	"testing"

	"github.com/carlmjohnson/bytemap"
	"golang.org/x/exp/maps"
)

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
		for j := 0; j < 3; j++ {
			s := strings.Repeat("a", i)
			charset := strings.Repeat("a", j)
			f.Add(s, charset)
			f.Add(s+"b", charset)
		}
	}
	f.Fuzz(func(t *testing.T, s, charset string) {
		want := naiveContains(s, charset)
		t.Run("Make", func(t *testing.T) {
			m := bytemap.Make(charset)
			testContainment(t, m, s, charset, want)
		})
		t.Run("WriteString", func(t *testing.T) {
			m := &bytemap.Bool{}
			n, err := m.WriteString(charset)
			if err != nil {
				t.Fatal(err)
			}
			if n != len(charset) {
				t.Fatal(len(charset))
			}
			testContainment(t, m, s, charset, want)
		})
		t.Run("Write", func(t *testing.T) {
			m := &bytemap.Bool{}
			n, err := m.Write([]byte(charset))
			if err != nil {
				t.Fatal(err)
			}
			if n != len(charset) {
				t.Fatal(len(charset))
			}
			testContainment(t, m, s, charset, want)
		})
		// Test io.Copy
		t.Run("Copy", func(t *testing.T) {
			m := &bytemap.Bool{}
			n64, err := io.Copy(m, strings.NewReader(charset))
			if err != nil {
				t.Fatal(err)
			}
			if n64 != int64(len(charset)) {
				t.Fatal(len(charset))
			}
			testContainment(t, m, s, charset, want)
		})
	})
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
		testGet(t, aMap, aNaive)
		bNaive := naiveMap(b)
		bMap := bytemap.Make(b)
		if !maps.Equal(bNaive, bMap.ToMap()) {
			t.Fatal(b, bMap)
		}
		testGet(t, bMap, bNaive)
		if maps.Equal(aNaive, bNaive) != aMap.Equals(bMap) {
			t.Fatal(aMap, bMap)
		}
	})
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

func FuzzRoundTrip(f *testing.F) {
	f.Add("")
	f.Add("a")
	f.Add("a")
	f.Add("ab")
	f.Add("abc")
	f.Fuzz(func(t *testing.T, charset string) {
		m1 := bytemap.Make(charset)
		bf := m1.ToBitField()
		m2 := bf.ToBool()
		if !m1.Equals(m2) {
			t.Fatal(m1, m2)
		}
	})
}
