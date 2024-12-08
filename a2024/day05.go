package a2024

import (
	"bufio"
	"io"
	"slices"
	"strings"
)

func day5(r io.Reader) (int, int) {
	scanner := bufio.NewScanner(r)
	rules := map[int][]int{}
	var updates [][]int
	for scanner.Scan() {
		text := strings.Split(scanner.Text(), "|")
		if len(text) == 2 {
			vals := mapValue(text, mustInt)
			rules[vals[0]] = append(rules[vals[0]], vals[1])
		} else {
			if text[0] == "" {
				continue
			}
			updates = append(updates, mapValue(strings.Split(text[0], ","), mustInt))
		}
	}

	validUpdates := 0
	middleSums := 0

	var invalidUpdates [][]int

updateLoop:
	for _, update := range updates {
		seen := map[int]struct{}{}
		for _, page := range update {
			if pageRules, mustBeBefore := rules[page]; mustBeBefore {
				for _, mustBeAfter := range pageRules {
					if _, invalidUpdate := seen[mustBeAfter]; invalidUpdate {
						invalidUpdates = append(invalidUpdates, update)
						continue updateLoop
					}
				}
			}
			seen[page] = struct{}{}
		}
		validUpdates++
		middleSums += update[(len(update)-1)/2]
	}

	fixedMiddleSums := 0
	for _, update := range invalidUpdates {
		slices.SortFunc(update, func(a, b int) int {
			if pageRules, ok := rules[a]; ok {
				if slices.Contains(pageRules, b) {
					return -1
				}
			}
			if pageRules, ok := rules[b]; ok {
				if slices.Contains(pageRules, a) {
					return 1
				}
			}
			return 0
		})
		fixedMiddleSums += update[(len(update)-1)/2]
	}

	return middleSums, fixedMiddleSums
}
