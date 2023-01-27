// Package bytemap contains types for making maps
// from bytes to bool, integer, or float.
// The maps are backed by arrays of 256 entries.
package bytemap

// Len is the length of a byte map array, 256 items.
const Len = 1 << 8
