package a2024

import (
	"testing"
)

func Test_day1(t *testing.T) {
	f := mustOpen(t, "control-01.txt")
	p1, p2 := day1(f)
	if p1 != 11 || p2 != 31 {
		t.Errorf("control should be 11, 31, was %d, %d", p1, p2)
		t.FailNow()
	}
	f = mustOpen(t, "day-01.txt")
	p1, p2 = day1(f)
	if p1 != 2086478 || p2 != 24941624 {
		t.Errorf("day1 is %d, %d", p1, p2)
	}
}

func Test_day2(t *testing.T) {
	f := mustOpen(t, "control-02.txt")
	p1, p2 := day2(f)
	if p1 != 2 || p2 != 0 {
		t.Errorf("control should be 2, 0, was %d, %d", p1, p2)
		t.FailNow()
	}
	f = mustOpen(t, "day-02.txt")
	p1, p2 = day2(f)
	if p1 != 0 || p2 != 0 {
		t.Errorf("day2 is %d, %d", p1, p2)
	}
}
