package a2022

import (
	"bytes"
	"io"
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
	b := bytes.NewReader([]byte(`1000
2000
3000

4000

5000
6000

7000
8000
9000

10000

`))
	cases := []struct {
		name       string
		in         io.Reader
		top1, top3 int
	}{
		{"control", b, 24000, 45000},
		{"personal", mustOpen(t, "inputs-2022/day1.txt"), 71924, 210406},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			top3 := day1(c.in)
			sum := 0
			for _, cal := range top3 {
				sum += cal
			}
			if top3[0] != c.top1 {
				t.Errorf("expected max %d, got %d", c.top1, top3[0])
			}
			if sum != c.top3 {
				t.Errorf("expected top 3 sum of %d, got %d", c.top3, sum)
			}
		})
	}
}

func Test_day2(t *testing.T) {
	controlInput := []byte(`A Y
B X
C Z
`)
	cases := []struct {
		name                string
		in                  io.Reader
		secondColumnAsPlays bool
		outcome             int
	}{
		{"control case - as plays", bytes.NewReader(controlInput), true, 15},
		{"control case - as outcomes", bytes.NewReader(controlInput), false, 12},
		{"personal input - as plays", mustOpen(t, "inputs-2022/day2.txt"), true, 13565},
		{"personal input - as outcomes", mustOpen(t, "inputs-2022/day2.txt"), false, 12424},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if v := day2(c.in, c.secondColumnAsPlays); v != c.outcome {
				t.Errorf("expected %d, got %d", c.outcome, v)
			}
		})
	}
}

func Test_day3(t *testing.T) {
	b := bytes.NewReader([]byte(`vJrwpWtwJgWrhcsFMMfFFhFp
jqHRNqRjqzjGDLGLrsFMfFZSrLrFZsSL
PmmdzqPrVvPwwTWBwg
wMqvLMZHhHMvwLHjbvcjnnSBnvTQFn
ttgJtRGJQctTZtZT
CrZsJsPPZsGzwwsLwLmpwMDw
`))
	cases := []struct {
		name         string
		in           io.Reader
		sacks, group int
	}{
		{"control case", b, 157, 70},
		{"personal input", mustOpen(t, "inputs-2022/day3.txt"), 7967, 2716},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			sacks, group, err := day3(c.in)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if sacks != c.sacks {
				t.Errorf("expected %d, got %d", c.sacks, sacks)
			}
			if group != c.group {
				t.Errorf("expected group score %d, got %d", c.group, group)
			}
		})
	}
}

func Test_day4(t *testing.T) {
	b := bytes.NewReader([]byte(`2-4,6-8
2-3,4-5
5-7,7-9
2-8,3-7
6-6,4-6
2-6,4-8`))
	cases := []struct {
		name                              string
		in                                io.Reader
		fullOverlapCount, anyOverlapCount int
	}{
		{"control case", b, 2, 4},
		{"personal input", mustOpen(t, "inputs-2022/day4.txt"), 441, 861},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if fullOverlapCount, anyOverlapCount, err := day4(c.in); err != nil {
				t.Errorf("unexpected error: %v", err)
			} else if fullOverlapCount != c.fullOverlapCount {
				t.Errorf("expected %d overlapping fullOverlap, got %d", c.fullOverlapCount, fullOverlapCount)
			} else if anyOverlapCount != c.anyOverlapCount {
				t.Errorf("expected %d with any overlap, got %d", c.anyOverlapCount, anyOverlapCount)
			}
		})
	}
}

func Test_day5(t *testing.T) {
	b := bytes.NewReader([]byte(`    [D]    
[N] [C]    
[Z] [M] [P]
 1   2   3 

move 1 from 2 to 1
move 3 from 1 to 3
move 2 from 2 to 1
move 1 from 1 to 2`))
	cases := []struct {
		name         string
		in           io.Reader
		part1, part2 string
	}{
		{"control case", b, "CMZ", "MCD"},
		{"personal input", mustOpen(t, "inputs-2022/day5.txt"), "FZCMJCRHZ", "JSDHQMZGF"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			part1, part2 := day5(c.in)
			if part1 != c.part1 {
				t.Errorf("expected part1 to be %s, got %s", c.part1, part1)
			}
			if part2 != c.part2 {
				t.Errorf("expected part2 to be %s, got %s", c.part2, part2)
			}
		})
	}
}

func Test_day6(t *testing.T) {
	cases := []struct {
		name                          string
		in                            io.Reader
		startOfPacket, startOfMessage int
	}{
		{"control case 1", bytes.NewReader([]byte(`mjqjpqmgbljsphdztnvjfqwrcgsmlb`)), 7, 19},
		{"control case 2", bytes.NewReader([]byte(`bvwbjplbgvbhsrlpgdmjqwftvncz`)), 5, 23},
		{"control case 3", bytes.NewReader([]byte(`nppdvjthqldpwncqszvftbrmjlhg`)), 6, 23},
		{"control case 4", bytes.NewReader([]byte(`nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg`)), 10, 29},
		{"control case 5", bytes.NewReader([]byte(`zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw`)), 11, 26},
		{"personal input", mustOpen(t, "inputs-2022/day6.txt"), 1598, 2414},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if r := day6(c.in); r.startOfPacket != c.startOfPacket {
				t.Errorf("expected start-of-packet %d, got %d", c.startOfPacket, r.startOfPacket)
			} else if r.startOfMessage != c.startOfMessage {
				t.Errorf("expected start-of-message %d, got %d", c.startOfMessage, r.startOfMessage)
			}
		})
	}
}
