package a2024

import (
	"fmt"
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
		sum += aElseB(solveMathDepth(p.answer, p.inputs, false), p.answer, 0)
	}
	withConcat := 0
	for _, p := range problems {
		withConcat += aElseB(solveMathDepth(p.answer, p.inputs, true), p.answer, 0)
	}
	return sum, withConcat
}

func solveMathDepth(target int, values []int, allowConcat bool) bool {
	left, right, rest := values[0], values[1], values[2:]
	sum := left + right
	mult := left * right
	concat := mustInt(fmt.Sprintf("%d%d", left, right))
	if len(rest) == 0 {
		return sum == target || mult == target || (allowConcat && concat == target)
	}
	return solveMathDepth(target, append([]int{sum}, rest...), allowConcat) ||
		solveMathDepth(target, append([]int{mult}, rest...), allowConcat) ||
		(allowConcat && solveMathDepth(target, append([]int{concat}, rest...), allowConcat))
}
