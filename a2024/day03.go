package a2024

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

func day3(r io.Reader) (int, int) {
	matcher := regexp.MustCompile(`mul\(\d+,\d+\)`)
	scanner := bufio.NewScanner(r)
	var mults []string
	for scanner.Scan() {
		mults = append(mults, matcher.FindAllString(scanner.Text(), -1)...)
	}
	var sum int
	for _, mult := range mults {
		tmp, _ := strings.CutPrefix(mult, "mul(")
		tmp, _ = strings.CutSuffix(tmp, ")")
		vals := mapValue(strings.Split(tmp, ","), func(e string) int {
			return mustInt(e)
		})
		sum += vals[0] * vals[1]
	}
	return sum, 0
}
