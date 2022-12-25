package a2022

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strings"
)

var snafuToInt = map[string]int{
	"2": 2, "1": 1, "0": 0,
	"-": -1, "=": -2,
}

func fivePow(exp int) int {
	return int(math.Pow(5, float64(exp)))
}

func sumSnafu(s string) int {
	powerOrdered := reverse(strings.Split(s, ""))
	sum := 0
	for i, snafu := range powerOrdered {
		placeValue := snafuToInt[snafu] * fivePow(i)
		sum += placeValue
	}
	return sum
}

func base5toSnafu(explodedBase5 []int) string {
	convertSnafuInts := make([]int, len(explodedBase5)+1)
	intToSnafuInt := map[int]int{
		3: -2,
		4: -1,
	}
	snafuIntToString := map[int]string{
		-2: "=", -1: "-",
		0: "0", 1: "1", 2: "2",
	}
	for i := len(explodedBase5) - 1; i >= 0; i-- {
		b5 := explodedBase5[i]
		if b5 >= 3 {
			explodedBase5[i-1] += 1
			convertSnafuInts[i+1] = intToSnafuInt[b5]
		} else {
			convertSnafuInts[i+1] = b5
		}
		for j := i; j >= 0; j-- {
			if explodedBase5[j] > 4 {
				explodedBase5[j] -= 5
				explodedBase5[j-1] += 1
			}
		}
	}
	return strings.Join(
		mapValue(
			aElseB(convertSnafuInts[0] == 0, convertSnafuInts[1:], convertSnafuInts),
			func(snafuInt int) string { return fmt.Sprintf("%s", snafuIntToString[snafuInt]) },
		), "")
}

func day25(in io.Reader) string {
	scanner := bufio.NewScanner(in)
	sum := 0
	for scanner.Scan() {
		sum += sumSnafu(scanner.Text())
	}
	var maxPower int
	for i := 0; fivePow(i) < sum; i++ {
		maxPower = i
	}
	var base5 []int
	for i := maxPower; i >= 0; i-- {
		powerVal := sum / fivePow(i)
		base5 = append(base5, powerVal)
		sum -= powerVal * fivePow(i)
	}

	return base5toSnafu(base5)
}
