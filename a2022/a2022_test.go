package a2022

import (
	"bytes"
	"fmt"
	"testing"
)

func Test_day1(t *testing.T) {
	inputReader, err := inputs.Open("inputs-2022/day1.txt")
	if err != nil {
		panic(err)
	}
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

	inputReader, err := inputs.Open("inputs-2022/day2.txt")
	if err != nil {
		fmt.Println(fmt.Errorf("input reading error: %v", err))
		t.FailNow()
	}
	fmt.Printf("score as plays: %d\n", day2(inputReader, true))

	b = bytes.NewReader([]byte(`A Y
B X
C Z
`))
	if day2(b, false) != 12 {
		fmt.Printf("control case should be 12")
		t.Fail()
	}
	inputReader, err = inputs.Open("inputs-2022/day2.txt")
	if err != nil {
		fmt.Println(fmt.Errorf("input reading error: %v", err))
		t.FailNow()
	}
	fmt.Printf("score as outcome: %d\n", day2(inputReader, false))
}
