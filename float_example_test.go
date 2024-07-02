//go:build goexperiment.rangefunc || go1.23

package bytemap_test

import (
	"fmt"
	"io"
	"os"

	"github.com/earthboundkid/bytemap/v2"
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
	for c, freq := range freqmap.MostCommon() {
		if freq > 0.02 {
			fmt.Printf("%q: %04.1f%%\n", c, freq*100)
		}
	}
	// Output:
	// ' ': 15.5%
	// 'e': 09.3%
	// 't': 06.9%
	// 'a': 06.0%
	// 'o': 05.5%
	// 'n': 05.2%
	// 'i': 05.0%
	// 's': 05.0%
	// 'h': 04.9%
	// 'r': 04.1%
	// 'l': 03.3%
	// 'd': 03.0%
	// 'u': 02.1%
}
