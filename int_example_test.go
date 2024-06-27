//go:build goexperiment.rangefunc || go1.23

package bytemap_test

import (
	"fmt"
	"io"
	"strings"

	"github.com/earthboundkid/bytemap/v2"
)

func ExampleInt_MostCommon() {
	var freqmap bytemap.Int
	r := strings.NewReader(`The quick brown fox jumps over the lazy dog.`)
	_, _ = io.Copy(&freqmap, r)
	for c, freq := range freqmap.MostCommon() {
		if freq > 0 {
			fmt.Printf("%q: %d\n", c, freq)
		}
	}
	// Output:
	// ' ': 8
	// 'o': 4
	// 'e': 3
	// 'h': 2
	// 'r': 2
	// 'u': 2
	// '.': 1
	// 'T': 1
	// 'a': 1
	// 'b': 1
	// 'c': 1
	// 'd': 1
	// 'f': 1
	// 'g': 1
	// 'i': 1
	// 'j': 1
	// 'k': 1
	// 'l': 1
	// 'm': 1
	// 'n': 1
	// 'p': 1
	// 'q': 1
	// 's': 1
	// 't': 1
	// 'v': 1
	// 'w': 1
	// 'x': 1
	// 'y': 1
	// 'z': 1
}
