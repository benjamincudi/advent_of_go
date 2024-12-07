package a2024

import (
	"bufio"
	"io"
	"sort"
	"strings"
)

func day1(r io.Reader) (int, int) {
	scanner := bufio.NewScanner(r)
	var left, right []int
	for scanner.Scan() {
		text := scanner.Text()
		parts := filter(strings.Split(text, " "), func(s string) bool {
			return s != ""
		})
		left = append(left, mustInt(parts[0]))
		right = append(right, mustInt(parts[1]))
	}
	sort.Ints(left)
	sort.Ints(right)
	var distance int
	leftCount := map[int]int{}
	var similarity int
	for i := range left {
		distance += abs(left[i] - right[i])
		if _, ok := leftCount[left[i]]; !ok {
			leftCount[left[i]] = len(filter(right, func(j int) bool {
				return j == left[i]
			}))
		}
		similarity += left[i] * leftCount[left[i]]
	}
	return distance, similarity
}
