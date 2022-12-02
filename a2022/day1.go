package a2022

import (
	"bufio"
	"io"
	"sort"
	"strconv"
	"strings"
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
	r := bufio.NewReader(inputReader)
	line := ""
	line, err = r.ReadString('\n')

	var elfTotals []int
	currentTotal := 0
	for err != io.EOF {
		if line != "\n" {
			if cal, err := strconv.Atoi(strings.TrimSuffix(line, "\n")); err != nil {
				panic(err)
			} else {
				currentTotal += cal
			}
		} else {
			elfTotals = append(elfTotals, currentTotal)
			currentTotal = 0
		}
		line, err = r.ReadString('\n')
	}
	elfTotals = append(elfTotals, currentTotal)
	sort.Sort(sort.Reverse(sort.IntSlice(elfTotals)))
	return elfTotals[0:3]
}
