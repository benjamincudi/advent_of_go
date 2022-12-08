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
// Mutates the input, but also returns it, so it can be inlined to function calls
func reverse[E any](in []E) []E {
	sort.SliceStable(in, func(i, j int) bool { return i > j })
	return in
}

func maxInt[E ~int](values ...E) E {
	max := values[0]
	for _, v := range values {
		if v > max {
			max = v
		}
	}
	return max
}

func minInt[E ~int](values ...E) E {
	max := values[0]
	for _, v := range values {
		if v < max {
			max = v
		}
	}
	return max
}

func mapValue[E any, F any](in []E, getValue func(E) F) []F {
	var ret []F
	for _, val := range in {
		ret = append(ret, getValue(val))
	}
	return ret
}
