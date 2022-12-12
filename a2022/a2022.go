package a2022

import (
	"embed"
	"sort"
	"strconv"
)

const shouldLog bool = false

//go:embed inputs-2022
var inputs embed.FS

func mustInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

// Make a reversed copy of the given slice
//
// Does not mutate the input, assign the result to your input if that's what you want
func reverse[E any](in []E) []E {
	dest := append(make([]E, 0, len(in)), in...)
	sort.SliceStable(dest, func(i, j int) bool { return i > j })
	return dest
}

func abs[E ~int](x E) E {
	if x < 0 {
		return x * -1
	}
	return x
}

func sign[E ~int](x E) E {
	if x < 0 {
		return -1
	}
	if x > 0 {
		return 1
	}
	return 0
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
	ret := make([]F, len(in))
	for i, val := range in {
		ret[i] = getValue(val)
	}
	return ret
}

func mapValueWithIndex[E any, F any](in []E, getValue func(int, E) F) []F {
	ret := make([]F, len(in))
	for i, val := range in {
		ret[i] = getValue(i, val)
	}
	return ret
}
