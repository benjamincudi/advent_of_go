package a2022

import (
	"bytes"
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
	if v := day2(b, true); v != 15 {
		t.Errorf("control case as plays should be 15, got %d", v)
	}

	if v := day2(mustOpen(t, "inputs-2022/day2.txt"), true); v != 13565 {
		t.Errorf("personal input as plays should be 13565, got %d", v)
	}

	b = bytes.NewReader([]byte(`A Y
B X
C Z
`))
	if v := day2(b, false); v != 12 {
		t.Errorf("control case as outcomes should be 12, got %d", v)
	}
	if v := day2(mustOpen(t, "inputs-2022/day2.txt"), false); v != 12424 {
		t.Errorf("personal input as outcomes should be 12424, got %d", v)
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

func Test_day5(t *testing.T) {
	b := bytes.NewReader([]byte(`    [D]    
[N] [C]    
[Z] [M] [P]
 1   2   3 

move 1 from 2 to 1
move 3 from 1 to 3
move 2 from 2 to 1
move 1 from 1 to 2`))
	part1, part2 := day5(b)
	if part1 != "CMZ" {
		t.Errorf("expected control part1 to be CMZ, got %s", part1)
	}
	if part2 != "MCD" {
		t.Errorf("expected control part2 to be MCD, got %s", part2)
	}

	part1, part2 = day5(mustOpen(t, "inputs-2022/day5.txt"))
	if part1 != "FZCMJCRHZ" {
		t.Errorf("expected personal input part1 to be FZCMJCRHZ, got %s", part1)
	}
	if part2 != "JSDHQMZGF" {
		t.Errorf("expected personal input part2 to be JSDHQMZGF, got %s", part2)
	}
}

func Test_day6(t *testing.T) {
	if r := day6(bytes.NewReader([]byte(`mjqjpqmgbljsphdztnvjfqwrcgsmlb`))); r.startOfPacket != 7 {
		t.Errorf("expected start-of-packet 7, got %d", r.startOfPacket)
	} else if r.startOfMessage != 19 {
		t.Errorf("expected start-of-message 19, got %d", r.startOfMessage)
	}
	if r := day6(bytes.NewReader([]byte(`bvwbjplbgvbhsrlpgdmjqwftvncz`))); r.startOfPacket != 5 {
		t.Errorf("expected start-of-packet 5, got %d", r.startOfPacket)
	} else if r.startOfMessage != 23 {
		t.Errorf("expected start-of-message 23, got %d", r.startOfMessage)
	}
	if r := day6(bytes.NewReader([]byte(`nppdvjthqldpwncqszvftbrmjlhg`))); r.startOfPacket != 6 {
		t.Errorf("expected start-of-packet 6, got %d", r.startOfPacket)
	} else if r.startOfMessage != 23 {
		t.Errorf("expected start-of-message 23, got %d", r.startOfMessage)
	}
	if r := day6(bytes.NewReader([]byte(`nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg`))); r.startOfPacket != 10 {
		t.Errorf("expected start-of-packet 10, got %d", r.startOfPacket)
	} else if r.startOfMessage != 29 {
		t.Errorf("expected start-of-message 29, got %d", r.startOfMessage)
	}
	if r := day6(bytes.NewReader([]byte(`zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw`))); r.startOfPacket != 11 {
		t.Errorf("expected start-of-packet 11, got %d", r.startOfPacket)
	} else if r.startOfMessage != 26 {
		t.Errorf("expected start-of-message 26, got %d", r.startOfMessage)
	}

	if r := day6(mustOpen(t, "inputs-2022/day6.txt")); r.startOfPacket != 1598 {
		t.Errorf("expected start-of-packet 1598, got %d", r.startOfPacket)
	} else if r.startOfMessage != 2414 {
		t.Errorf("expected start-of-message 2414, got %d", r.startOfMessage)
	}
}
