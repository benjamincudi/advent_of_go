package a2022

import (
	"github.com/gocarina/gocsv"
	"sort"
	"strconv"
)

func day1() []int {
	type Calories struct {
		Cal      string `csv:"lmao"`
		NotThere string `cav:"wrong"`
	}

	inputReader, err := inputs.Open("inputs-2022/day1.txt")
	if err != nil {
		panic(err)
	}
	var elfSnacks []Calories
	if err := gocsv.UnmarshalWithoutHeaders(inputReader, &elfSnacks); err != nil {
		panic(err)
	}
	var elfTotals []int
	currentTotal := 0
	for _, snack := range elfSnacks {
		if snack.Cal != "" {
			if cal, err := strconv.Atoi(snack.Cal); err != nil {
				panic(err)
			} else {
				currentTotal += cal
			}
			continue
		}
		elfTotals = append(elfTotals, currentTotal)
		currentTotal = 0
	}
	elfTotals = append(elfTotals, currentTotal)
	sort.Sort(sort.Reverse(sort.IntSlice(elfTotals)))
	return elfTotals[0:3]
}
