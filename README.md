# bytemap [![GoDoc](https://godoc.org/github.com/carlmjohnson/bytemap?status.svg)](https://godoc.org/github.com/carlmjohnson/bytemap) [![Go Report Card](https://goreportcard.com/badge/github.com/carlmjohnson/bytemap)](https://goreportcard.com/report/github.com/carlmjohnson/bytemap) [![Gocover.io](https://gocover.io/_badge/github.com/carlmjohnson/bytemap)](https://gocover.io/github.com/carlmjohnson/bytemap)

Bytemap contains types for making maps from bytes to bool, integer, or float using a backing array.

## Benchmarks

Micro-benchmarks are usually not a good way to evaluate systems. That said, using a bytemap array can be very fast while also providing readable code.

Let's say you want to test that string contains only digits. A very fast way is just to write a loop:

```go
match := true
for _, c := range []byte(s) {
    if c < '0' || c > '9' {
        match = false
        break
    }
}
```

This is very fast, but the code is somewhat tedious. One might decide to replace it with a simple regular expression.

```go
r := regexp.MustCompile(`^[0-9]*$`)
match := r.MatchString(s)
```

This is much shorter, but it's actually a little tricky to read if you're not very familiar with regular expressions, and it's much slower to execute.

Another idea might be to test against a `map[byte]bool`. This turns out to be almost as slow as the regular expression and about as verbose as the simple loop test.

A bytemap is **short, simple, and fast**:

```go
m := bytemap.Make("0123456789")
match := m.Contains(s)
```

Take these benchmarks with a grain of salt, but they show a bytemap can actually perform as well as a handwritten loop:

```
goos: darwin
goarch: amd64
pkg: github.com/carlmjohnson/bytemap
BenchmarkLoop-8                  184966533       6.314 ns/op
BenchmarkBoolContains-8          162605607       7.503 ns/op
BenchmarkBitFieldContains-8       80200012      16.85 ns/op
BenchmarkMapByteEmpty-8           53849732      23.19 ns/op
BenchmarkMapByteBool-8             6663114     165.9 ns/op
BenchmarkRegexp-8                  6080572     197.5 ns/op
BenchmarkRegexpSlow-8               330384    3252 ns/op
```

## How does it work?

There are only 256 different possible bit patterns in a byte, so `bytemap.Bool` just preallocates an array of 256 entries.

`bytemap.BitField` only allocates one bit per entry, which makes it 8 times smaller than `bytemap.Bool`, only 32 bytes long. In many cases however, it will be a bit slower than using a `bytemap.Bool`.
