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

func Test_day1(t *testing.T) {
	testForDay(t, "01", day1, 11, 31, 2086478, 24941624)
}

func Test_day2(t *testing.T) {
	testForDay(t, "02", day2, 2, 4, 680, 710)
}

func Test_day3(t *testing.T) {
	testForDay(t, "03", day3, 161, 48, 187833789, 94455185)
}
