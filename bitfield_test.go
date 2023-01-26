package bytemap_test

import (
	"io"
	"strings"
	"testing"

	"github.com/carlmjohnson/bytemap"
	"golang.org/x/exp/maps"
)

func FuzzMakeBitField(f *testing.F) {
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
			m := bytemap.MakeBitField(charset)
			testContainmentBitfield(t, m, s, charset, want)
		}
		// Test WriteString
		{
			m := &bytemap.BitField{}
			n, err := m.WriteString(charset)
			if err != nil {
				t.Fatal(err)
			}
			if n != len(charset) {
				t.Fatal(len(charset))
			}
			testContainmentBitfield(t, m, s, charset, want)
		}
		// Test Write
		{
			m := &bytemap.BitField{}
			n, err := m.Write([]byte(charset))
			if err != nil {
				t.Fatal(err)
			}
			if n != len(charset) {
				t.Fatal(len(charset))
			}
			testContainmentBitfield(t, m, s, charset, want)
		}
		// Test io.Copy
		{
			m := &bytemap.BitField{}
			n64, err := io.Copy(m, strings.NewReader(charset))
			if err != nil {
				t.Fatal(err)
			}
			if n64 != int64(len(charset)) {
				t.Fatal(len(charset))
			}
			testContainmentBitfield(t, m, s, charset, want)
		}
	})
}

func testContainmentBitfield(t *testing.T, m *bytemap.BitField, s, charset string, want bool) {
	if want != m.Contains(s) {
		t.Fatal(want, s, charset, m)
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

func FuzzBitFieldToMap(f *testing.F) {
	f.Add("", "")
	f.Add("a", "b")
	f.Add(
		"the quick brown fox jumps over a lazy dog.",
		"abcdefghijklmnopqrstuvwxyz. ",
	)
	f.Fuzz(func(t *testing.T, a, b string) {
		aNaive := naiveMap(a)
		aMap := bytemap.MakeBitField(a)
		if !maps.Equal(aNaive, aMap.ToMap()) {
			t.Fatalf("input=%q want=%v got=%v",
				a, aNaive, aMap.ToMap())
		}
		bNaive := naiveMap(b)
		bMap := bytemap.MakeBitField(b)
		if !maps.Equal(bNaive, bMap.ToMap()) {
			t.Fatal(b, bMap)
		}
		if maps.Equal(aNaive, bNaive) != aMap.Equals(bMap) {
			t.Fatal(aMap, bMap)
		}
	})
}
