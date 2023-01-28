package bytemap_test

import (
	"fmt"
	"io"
	"os"

	"github.com/carlmjohnson/bytemap"
)

func ExampleFloat_SetFrequencies() {
	f, err := os.Open("testdata/moby-dick.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var freqmap bytemap.Float
	_, _ = io.Copy(&freqmap, f)
	freqmap.SetFrequencies()
	for _, freq := range freqmap.MostCommon() {
		if freq.N > 0.02 {
			fmt.Printf("%q: %04.1f%%\n", []byte{freq.Byte}, freq.N*100)
		}
	}
	// Output:
	// " ": 15.5%
	// "e": 09.3%
	// "t": 06.9%
	// "a": 06.0%
	// "o": 05.5%
	// "n": 05.2%
	// "i": 05.0%
	// "s": 05.0%
	// "h": 04.9%
	// "r": 04.1%
	// "l": 03.3%
	// "d": 03.0%
	// "u": 02.1%
}
