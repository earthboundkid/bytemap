package bytemap_test

import (
	"bytes"
	"io"
	"os"
	"regexp"
	"testing"

	"github.com/carlmjohnson/bytemap"
)

var globalNaive map[byte]bool

func BenchmarkNaive(b *testing.B) {
	data, err := os.ReadFile("testdata/moby-dick.txt")
	if err != nil {
		b.Fatal(err)
	}
	s := string(data)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		globalNaive = naiveMap(s)
	}
}

var globalBool *bytemap.Bool

func BenchmarkBoolCopy(b *testing.B) {
	data, err := os.ReadFile("testdata/moby-dick.txt")
	if err != nil {
		b.Fatal(err)
	}
	var m bytemap.Bool
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = io.Copy(&m, bytes.NewReader(data))
		if err != nil {
			b.Fatal(err)
		}
	}
	globalBool = &m
}

func BenchmarkBoolWriteString(b *testing.B) {
	data, err := os.ReadFile("testdata/moby-dick.txt")
	if err != nil {
		b.Fatal(err)
	}
	s := string(data)
	var m bytemap.Bool
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.WriteString(s)
	}
	globalBool = &m
}

func BenchmarkBoolWrite(b *testing.B) {
	data, err := os.ReadFile("testdata/moby-dick.txt")
	if err != nil {
		b.Fatal(err)
	}
	var m bytemap.Bool
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Write(data)
	}
	globalBool = &m
}

var globalBitField *bytemap.BitField

func BenchmarkBitFieldCopy(b *testing.B) {
	data, err := os.ReadFile("testdata/moby-dick.txt")
	if err != nil {
		b.Fatal(err)
	}
	var m bytemap.BitField
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = io.Copy(&m, bytes.NewReader(data))
		if err != nil {
			b.Fatal(err)
		}
	}
	globalBitField = &m
}

func BenchmarkBitFieldWriteString(b *testing.B) {
	data, err := os.ReadFile("testdata/moby-dick.txt")
	if err != nil {
		b.Fatal(err)
	}
	s := string(data)
	var m bytemap.BitField
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.WriteString(s)
	}
	globalBitField = &m
}

var globalMatch bool

var testStrings = []string{
	"12356789",
	"12356789a",
	"1235a6789",
	"987654321",
}

func BenchmarkLoop(b *testing.B) {
	var match bool
	for i := 0; i < b.N; i++ {
		s := testStrings[i%len(testStrings)]
		match = true
		for _, c := range []byte(s) {
			if c < '0' || c > '9' {
				match = false
				break
			}
		}
	}
	globalMatch = match
}

func BenchmarkRegexp(b *testing.B) {
	r := regexp.MustCompile(`^[0-9]*$`)
	b.ResetTimer()
	var match bool
	for i := 0; i < b.N; i++ {
		s := testStrings[i%len(testStrings)]
		match = r.MatchString(s)
	}
	globalMatch = match
}

func BenchmarkNaiveContains(b *testing.B) {
	m := naiveMap("0123456789")
	b.ResetTimer()
	var match bool
	for i := 0; i < b.N; i++ {
		s := testStrings[i%len(testStrings)]
		match = naiveMapContains(m, s)
	}
	globalMatch = match
}

func BenchmarkBoolContains(b *testing.B) {
	m := bytemap.Make("0123456789")
	b.ResetTimer()
	var match bool
	for i := 0; i < b.N; i++ {
		s := testStrings[i%len(testStrings)]
		match = m.Contains(s)
	}
	globalMatch = match
}

func BenchmarkBitFieldContains(b *testing.B) {
	m := bytemap.MakeBitField("0123456789")
	b.ResetTimer()
	var match bool
	for i := 0; i < b.N; i++ {
		s := testStrings[i%len(testStrings)]
		match = m.Contains(s)
	}
	globalMatch = match
}
