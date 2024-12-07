package a2024

import (
	"bufio"
	"io"
	"slices"
	"strings"
)

func day2(r io.Reader) (int, int) {
	scanner := bufio.NewScanner(r)
	var reports [][]int
	for scanner.Scan() {
		text := scanner.Text()
		parts := filter(strings.Split(text, " "), func(s string) bool {
			return s != ""
		})
		reports = append(reports, mapValue(parts, mustInt))
	}
	var safeReportCount, dampenedSafeCount int
	for _, report := range reports {
		filtered, signs := processReport(report)
		if len(filtered) == 0 && len(signs) == 1 {
			safeReportCount += 1
			dampenedSafeCount += 1
			continue
		}
		for i := range report {
			filtered, signs = processReport(removeIndex(report, i))
			if len(filtered) == 0 && len(signs) == 1 {
				dampenedSafeCount += 1
				break
			}
		}
	}
	return safeReportCount, dampenedSafeCount
}

func processReport(report []int) ([]int, []int) {
	diffs := mapValueWithIndex(report, func(i int, e int) int {
		if i == len(report)-1 {
			return e - report[i-1] // repeat the previous calc
		}
		return report[i+1] - e
	})
	filtered := filter(diffs, func(i int) bool {
		return abs(i) > 3 || i == 0
	})
	signs := slices.Compact(mapValue(diffs, func(e int) int {
		return sign(e)
	}))
	return filtered, signs
}
