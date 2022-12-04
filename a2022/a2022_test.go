package a2022

import (
	"bytes"
	"fmt"
	"io/fs"
	"testing"
)

func mustOpen(t *testing.T, name string) fs.File {
	inputReader, err := inputs.Open(name)
	if err != nil {
		t.Fatalf("%v", err)
	}
	return inputReader
}

func Test_day1(t *testing.T) {
	inputReader := mustOpen(t, "inputs-2022/day1.txt")
	top3 := day1(inputReader)
	sum := 0
	for _, c := range top3 {
		sum += c
	}
	if top3[0] != 71924 {
		t.Errorf("expected max 71924, got %d", top3[0])
	}
	if sum != 210406 {
		t.Errorf("expected top 3 sum of 210406, got %d", sum)
	}
}

func Test_day2(t *testing.T) {
	b := bytes.NewReader([]byte(`A Y
B X
C Z
`))
	if day2(b, true) != 15 {
		fmt.Printf("control case should be 15")
		t.Fail()
	}

	inputReader := mustOpen(t, "inputs-2022/day2.txt")
	fmt.Printf("score as plays: %d\n", day2(inputReader, true))

	b = bytes.NewReader([]byte(`A Y
B X
C Z
`))
	if day2(b, false) != 12 {
		fmt.Printf("control case should be 12")
		t.Fail()
	}
	inputReader = mustOpen(t, "inputs-2022/day2.txt")
	fmt.Printf("score as outcome: %d\n", day2(inputReader, false))
}

func Test_day3(t *testing.T) {
	b := bytes.NewReader([]byte(`vJrwpWtwJgWrhcsFMMfFFhFp
jqHRNqRjqzjGDLGLrsFMfFZSrLrFZsSL
PmmdzqPrVvPwwTWBwg
wMqvLMZHhHMvwLHjbvcjnnSBnvTQFn
ttgJtRGJQctTZtZT
CrZsJsPPZsGzwwsLwLmpwMDw
`))
	sacks, group, err := day3(b)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if sacks != 157 {
		t.Errorf("expected 157, got %d", sacks)
	}
	if group != 70 {
		t.Errorf("expected group score 70, got %d", group)
	}

	input := mustOpen(t, "inputs-2022/day3.txt")
	sacks, group, err = day3(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if sacks != 7967 {
		t.Errorf("expected 157, got %d", sacks)
	}
	if group != 2716 {
		t.Errorf("expected group score 70, got %d", group)
	}
}

func Test_day4(t *testing.T) {
	b := bytes.NewReader([]byte(`2-4,6-8
2-3,4-5
5-7,7-9
2-8,3-7
6-6,4-6
2-6,4-8`))
	if fullOverlap, anyOverlap, err := day4(b); err != nil {
		t.Errorf("unexpected error: %v", err)
	} else if fullOverlap != 2 {
		t.Errorf("expected 2 overlapping fullOverlap, got %d", fullOverlap)
	} else if anyOverlap != 4 {
		t.Errorf("expected 4 with any overlap, got %d", anyOverlap)
	}
	input := mustOpen(t, "inputs-2022/day4.txt")
	if fullOverlap, anyOverlap, err := day4(input); err != nil {
		t.Errorf("unexpected error: %v", err)
	} else if fullOverlap != 441 {
		t.Errorf("expected 441 overlapping fullOverlap, got %d", fullOverlap)
	} else if anyOverlap != 861 {
		t.Errorf("expected 861 with any overlap, got %d", anyOverlap)
	}
}
