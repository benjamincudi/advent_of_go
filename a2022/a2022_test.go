package a2022

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"strings"
	"testing"
)

func mustOpen(t *testing.T, name string) fs.File {
	t.Helper()
	inputReader, err := inputs.Open(fmt.Sprintf("inputs-2022/%s", name))
	if err != nil {
		t.Fatalf("%v", err)
	}
	return inputReader
}

func Test_day1(t *testing.T) {
	testCases := []struct {
		name       string
		in         io.Reader
		top1, top3 int
	}{
		{"control", mustOpen(t, "control1.txt"), 24000, 45000},
		{"personal", mustOpen(t, "day1.txt"), 71924, 210406},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			top3 := day1(tc.in)
			sum := 0
			for _, cal := range top3 {
				sum += cal
			}
			if top3[0] != tc.top1 {
				t.Errorf("expected max %d, got %d", tc.top1, top3[0])
			}
			if sum != tc.top3 {
				t.Errorf("expected top 3 sum of %d, got %d", tc.top3, sum)
			}
		})
	}
}

func Test_day2(t *testing.T) {
	testCases := []struct {
		name                string
		in                  io.Reader
		secondColumnAsPlays bool
		outcome             int
	}{
		{"control case - as plays", mustOpen(t, "control2.txt"), true, 15},
		{"control case - as outcomes", mustOpen(t, "control2.txt"), false, 12},
		{"personal input - as plays", mustOpen(t, "day2.txt"), true, 13565},
		{"personal input - as outcomes", mustOpen(t, "day2.txt"), false, 12424},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if v := day2(tc.in, tc.secondColumnAsPlays); v != tc.outcome {
				t.Errorf("expected %d, got %d", tc.outcome, v)
			}
		})
	}
}

func Test_day3(t *testing.T) {
	testCases := []struct {
		name         string
		in           io.Reader
		sacks, group int
	}{
		{"control case", mustOpen(t, "control3.txt"), 157, 70},
		{"personal input", mustOpen(t, "day3.txt"), 7967, 2716},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sacks, group, err := day3(tc.in)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if sacks != tc.sacks {
				t.Errorf("expected %d, got %d", tc.sacks, sacks)
			}
			if group != tc.group {
				t.Errorf("expected group score %d, got %d", tc.group, group)
			}
		})
	}
}

func Test_day4(t *testing.T) {
	testCases := []struct {
		name                              string
		in                                io.Reader
		fullOverlapCount, anyOverlapCount int
	}{
		{"control case", mustOpen(t, "control4.txt"), 2, 4},
		{"personal input", mustOpen(t, "day4.txt"), 441, 861},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if fullOverlapCount, anyOverlapCount, err := day4(tc.in); err != nil {
				t.Errorf("unexpected error: %v", err)
			} else if fullOverlapCount != tc.fullOverlapCount {
				t.Errorf("expected %d overlapping fullOverlap, got %d", tc.fullOverlapCount, fullOverlapCount)
			} else if anyOverlapCount != tc.anyOverlapCount {
				t.Errorf("expected %d with any overlap, got %d", tc.anyOverlapCount, anyOverlapCount)
			}
		})
	}
}

func Test_day5(t *testing.T) {
	testCases := []struct {
		name         string
		in           io.Reader
		part1, part2 string
	}{
		{"control case", mustOpen(t, "control5.txt"), "CMZ", "MCD"},
		{"personal input", mustOpen(t, "day5.txt"), "FZCMJCRHZ", "JSDHQMZGF"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			part1, part2 := day5(tc.in)
			if part1 != tc.part1 {
				t.Errorf("expected part1 to be %s, got %s", tc.part1, part1)
			}
			if part2 != tc.part2 {
				t.Errorf("expected part2 to be %s, got %s", tc.part2, part2)
			}
		})
	}
}

func Test_day6(t *testing.T) {
	r := bufio.NewScanner(mustOpen(t, "control6.txt"))
	var control []string
	for r.Scan() {
		control = append(control, r.Text())
	}
	testCases := []struct {
		name                          string
		in                            io.Reader
		startOfPacket, startOfMessage int
	}{
		{"control case 1", strings.NewReader(control[0]), 7, 19},
		{"control case 2", strings.NewReader(control[1]), 5, 23},
		{"control case 3", strings.NewReader(control[2]), 6, 23},
		{"control case 4", strings.NewReader(control[3]), 10, 29},
		{"control case 5", strings.NewReader(control[4]), 11, 26},
		{"personal input", mustOpen(t, "day6.txt"), 1598, 2414},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if r := day6(tc.in); r.startOfPacket != tc.startOfPacket {
				t.Errorf("expected start-of-packet %d, got %d", tc.startOfPacket, r.startOfPacket)
			} else if r.startOfMessage != tc.startOfMessage {
				t.Errorf("expected start-of-message %d, got %d", tc.startOfMessage, r.startOfMessage)
			}
		})
	}
}
