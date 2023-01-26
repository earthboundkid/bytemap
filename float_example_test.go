package bytemap_test

import (
	"fmt"
	"io"
	"strings"

	"github.com/carlmjohnson/bytemap"
)

func ExampleFloat_SetFrequencies() {
	var freqmap bytemap.Float
	r := strings.NewReader(`The quick brown fox jumps over the lazy dog.`)
	_, _ = io.Copy(&freqmap, r)
	freqmap.SetFrequencies()
	for _, freq := range freqmap.MostCommon() {
		if freq.N > 0 {
			fmt.Printf("%q: %04.1f%%\n", []byte{freq.Byte}, freq.N*100)
		}
	}
	// Output:
	// " ": 18.2%
	// "o": 09.1%
	// "e": 06.8%
	// "h": 04.5%
	// "r": 04.5%
	// "u": 04.5%
	// ".": 02.3%
	// "T": 02.3%
	// "a": 02.3%
	// "b": 02.3%
	// "c": 02.3%
	// "d": 02.3%
	// "f": 02.3%
	// "g": 02.3%
	// "i": 02.3%
	// "j": 02.3%
	// "k": 02.3%
	// "l": 02.3%
	// "m": 02.3%
	// "n": 02.3%
	// "p": 02.3%
	// "q": 02.3%
	// "s": 02.3%
	// "t": 02.3%
	// "v": 02.3%
	// "w": 02.3%
	// "x": 02.3%
	// "y": 02.3%
	// "z": 02.3%
}
