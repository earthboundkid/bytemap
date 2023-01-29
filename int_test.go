package bytemap_test

import (
	"io"
	"strings"
	"testing"

	"github.com/carlmjohnson/bytemap"
	"golang.org/x/exp/maps"
)

func FuzzInt(f *testing.F) {
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
			m := &bytemap.Int{}
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
			m := &bytemap.Int{}
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

func FuzzIntSet(f *testing.F) {
	f.Add("", "", "")
	f.Add("a", "a", "a")
	f.Add("abc", "bcde", "b")
	f.Fuzz(func(t *testing.T, add, remove, restore string) {
		var mInt bytemap.Int
		m := make(map[byte]int)
		for _, c := range []byte(add) {
			mInt.Set(c, 1)
			m[c] = 1
		}
		for _, c := range []byte(remove) {
			mInt.Set(c, 0)
			m[c] = 0
		}
		for _, c := range []byte(restore) {
			mInt.Set(c, 1)
			m[c] = 1
		}
		// Fill in blanks
		for i := 0; i < bytemap.Len; i++ {
			m[byte(i)] = m[byte(i)]
		}
		for i := 0; i < bytemap.Len; i++ {
			if mInt.Get(byte(i)) != m[byte(i)] {
				t.Fatal(i, mInt.Get(byte(i)), m[byte(i)])
			}
		}
		if !maps.Equal(mInt.ToMap(), m) {
			t.Fatal(mInt)
		}
	})
}
