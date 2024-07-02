# bytemap [![GoDoc](https://godoc.org/github.com/earthboundkid/bytemap?status.svg)](https://godoc.org/github.com/earthboundkid/bytemap/v2) [![Go Report Card](https://goreportcard.com/badge/github.com/earthboundkid/bytemap)](https://goreportcard.com/report/github.com/earthboundkid/bytemap) [![Coverage Status](https://coveralls.io/repos/github/earthboundkid/bytemap/badge.svg)](https://coveralls.io/github/earthboundkid/bytemap)

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

Take these benchmarks with a grain of salt, but they show a bytemap can actually perform as well as a handwritten loop or better:

```
goos: darwin
goarch: amd64
pkg: github.com/earthboundkid/bytemap
BenchmarkBoolContains-8         318648963        3.762 ns/op
BenchmarkBitFieldContains-8     216729614        5.526 ns/op
BenchmarkLoop-8                 170478852        6.954 ns/op
BenchmarkMapByteEmpty-8         143119087        8.991 ns/op
BenchmarkMapByteBool-8           19294774       61.44 ns/op
BenchmarkRegexp-8                11613380      102.7 ns/op
BenchmarkRegexpSlow-8              822931     1289 ns/op
```

## How does it work?

There are only 256 different possible bit patterns in a byte, so `bytemap.Bool` just preallocates an array of 256 entries.

`bytemap.BitField` only allocates one bit per entry, which makes it 8 times smaller than `bytemap.Bool`, only 32 bytes long. In many cases however, it will be a bit slower than using a `bytemap.Bool`.
