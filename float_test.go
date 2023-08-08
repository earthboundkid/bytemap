package bytemap_test

import (
	"io"
	"maps"
	"strings"
	"testing"

	"github.com/carlmjohnson/bytemap"
)

func FuzzFloat(f *testing.F) {
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
		t.Run("WriteString", func(t *testing.T) {
			m := &bytemap.Float{}
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
			m := &bytemap.Float{}
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
			m := &bytemap.Int{}
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

func FuzzFloatSet(f *testing.F) {
	f.Add("", "", "")
	f.Add("a", "a", "a")
	f.Add("abc", "bcde", "b")
	f.Fuzz(func(t *testing.T, add, remove, restore string) {
		var mFloat bytemap.Float
		m := make(map[byte]float64)
		for _, c := range []byte(add) {
			mFloat.Set(c, 1)
			m[c] = 1
		}
		for _, c := range []byte(remove) {
			mFloat.Set(c, 0)
			m[c] = 0
		}
		for _, c := range []byte(restore) {
			mFloat.Set(c, 1)
			m[c] = 1
		}
		// Fill in blanks
		for i := 0; i < bytemap.Len; i++ {
			m[byte(i)] = m[byte(i)]
		}
		for i := 0; i < bytemap.Len; i++ {
			if mFloat.Get(byte(i)) != m[byte(i)] {
				t.Fatal(i, mFloat.Get(byte(i)), m[byte(i)])
			}
		}
		if !maps.Equal(mFloat.ToMap(), m) {
			t.Fatal(mFloat)
		}
	})
}
