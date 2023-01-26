// Package bytemap contains types for making maps
// from bytes to bool, integer, or float.
// The maps are backed by arrays of 256 entries.
package bytemap

// Size is the size of a byte map array, 256.
const Size = 1 << 8
