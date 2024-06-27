//go:build go1.23 || goexperiment.rangefunc

package bytemap

import (
	"cmp"
	"iter"
	"slices"
)

// MostCommon returns a sequence of characters and counts for m
// from highest count to lowest.
func (m *Int) MostCommon() iter.Seq2[byte, int] {
	return func(yield func(byte, int) bool) {
		freqs := make([]byte, Len)
		for i := range freqs {
			freqs[i] = byte(i)
		}
		slices.SortStableFunc(freqs[:], func(a, b byte) int {
			return cmp.Compare(m[b], m[a])
		})
		for _, c := range freqs {
			if !yield(c, m[c]) {
				return
			}
		}
	}
}
