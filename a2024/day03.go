package a2024

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

func day3(r io.Reader) (int, int) {
	matcher := regexp.MustCompile(`mul\(\d+,\d+\)`)
	instructionMatcher := regexp.MustCompile(`(mul\(\d+,\d+\))|(do\(\))|(don't\(\))`)
	scanner := bufio.NewScanner(r)
	var mults []string
	var instructions []string
	for scanner.Scan() {
		text := scanner.Text()
		mults = append(mults, matcher.FindAllString(text, -1)...)
		instructions = append(instructions, instructionMatcher.FindAllString(text, -1)...)
	}
	var multPairs [][]int
	for _, mult := range mults {
		tmp, _ := strings.CutPrefix(mult, "mul(")
		tmp, _ = strings.CutSuffix(tmp, ")")
		vals := mapValue(strings.Split(tmp, ","), func(e string) int {
			return mustInt(e)
		})
		multPairs = append(multPairs, vals)
	}
	var sum, conditionalSum int
	do := true
	multIndex := 0
	for _, val := range instructions {
		if val == "do()" {
			do = true
		} else if val == "don't()" {
			do = false
		} else {
			res := multPairs[multIndex][0] * multPairs[multIndex][1]
			sum += res
			if do {
				conditionalSum += res
			}
			multIndex += 1
		}
	}
	return sum, conditionalSum
}
