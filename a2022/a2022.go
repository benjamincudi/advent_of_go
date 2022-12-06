package a2022

import (
	"embed"
	"sort"
	"strconv"
)

//go:embed inputs-2022
var inputs embed.FS

func mustInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

// Reverse the given slice
//
// Mutates the input, but also returns it so it can be inlined to function calls
func reverse[E any](in []E) []E {
	sort.SliceStable(in, func(i, j int) bool { return i > j })
	return in
}
