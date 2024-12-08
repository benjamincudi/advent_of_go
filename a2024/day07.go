package a2024

import (
	"io"
	"regexp"
)

func day7(r io.Reader) (int, int) {
	matchDigits := regexp.MustCompile(`[[:digit:]]+`)
	type problem struct {
		answer int
		inputs []int
	}
	var problems []problem
	mapFileLines(r, func(text string) {
		rowVals := mapValue(matchDigits.FindAllString(text, -1), mustInt)
		problems = append(problems, problem{rowVals[0], rowVals[1:]})
	})
	sum := 0
	for _, p := range problems {
		sum += aElseB(solveMathDepth(p.answer, p.inputs), p.answer, 0)
	}
	return sum, 0
}

func solveMathDepth(target int, values []int) bool {
	left, right, rest := values[0], values[1], values[2:]
	sum := left + right
	mult := left * right
	if len(rest) == 0 {
		return sum == target || mult == target
	}
	return solveMathDepth(target, append([]int{sum}, rest...)) || solveMathDepth(target, append([]int{mult}, rest...))
}
