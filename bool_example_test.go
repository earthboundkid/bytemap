package bytemap_test

import (
	"fmt"

	"github.com/earthboundkid/bytemap/v2"
)

func ExampleRange() {
	ascii := bytemap.Range(0, 127)
	fmt.Println(ascii.Contains("Hello, world"))
	fmt.Println(ascii.Contains("Hello, 🌎"))

	upper := bytemap.Range('A', 'Z')
	nonupper := upper.Invert()
	fmt.Println(nonupper.Contains("hello, world!"))

	// Output:
	// true
	// false
	// true
}

func ExampleUnion() {
	upper := bytemap.Range('A', 'Z')
	lower := bytemap.Range('a', 'z')
	alpha := bytemap.Union(upper, lower)
	word := bytemap.Union(
		upper,
		lower,
		bytemap.Range('0', '9'),
		bytemap.Make("_"),
	)
	fmt.Println(alpha.Contains("CamelCase"))
	fmt.Println(alpha.Contains("snake_case"))
	fmt.Println(word.Contains("snake_case"))
	// Output:
	// true
	// false
	// true
}
