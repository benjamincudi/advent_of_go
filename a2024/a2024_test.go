package a2024

import (
	"fmt"
	"io"
	"testing"
)

func testForDay(t *testing.T, day string, solver func(r io.Reader) (int, int), control1, control2, result1, result2 int) {
	t.Helper()
	f := mustOpen(t, fmt.Sprintf("control-%s.txt", day))
	p1, p2 := solver(f)
	if p1 != control1 || p2 != control2 {
		t.Errorf("control should be %d, %d, was %d, %d", control1, control2, p1, p2)
		t.FailNow()
	}
	f = mustOpen(t, fmt.Sprintf("day-%s.txt", day))
	p1, p2 = solver(f)
	if p1 != result1 || p2 != result2 {
		t.Errorf("day%s is %d, %d", day, p1, p2)
	}
}

func Test_a2024(t *testing.T) {
	for i, tc := range []struct {
		solver func(r io.Reader) (int, int)
		control1, control2,
		result1, result2 int
	}{
		{day1, 11, 31, 2086478, 24941624},
		{day2, 2, 4, 680, 710},
		{day3, 161, 48, 187833789, 94455185},
		{day4, 18, 9, 2532, 1941},
		{day5, 143, 123, 6612, 4944},
		{day6, 41, 6, 4778, 1618},
		{day7, 3749, 11387, 3312271365652, 509463489296712},
		{day8, 14, 34, 327, 1233},
		{day9, 1928, 2858, 6200294120911, 6227018762750},
	} {
		prefix := aElseB(i < 9, "0", "")
		day := fmt.Sprintf("%s%d", prefix, i+1)
		t.Run(fmt.Sprintf("day %s", day), func(t *testing.T) {
			testForDay(t, day, tc.solver, tc.control1, tc.control2, tc.result1, tc.result2)
		})
	}
}
