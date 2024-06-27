//go:build goexperiment.rangefunc || go1.23

package bytemap

import (
	"cmp"
	"iter"
	"slices"
)

// MostCommon returns a sequence of characters and counts for m
// from highest count to lowest.
func (m *Float) MostCommon() iter.Seq2[byte, float64] {
	return func(yield func(byte, float64) bool) {
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
