package a2022

import (
	"embed"
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
