package a2022

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strings"
)

type packetCompareResult int

const (
	inOrder packetCompareResult = iota + 1
	outOfOrder
	indeterminate
)

func (c packetCompareResult) string() string {
	if c == inOrder {
		return "in order"
	}
	if c == outOfOrder {
		return "out of order"
	}
	return "indeterminate"
}

type intOrSlice struct {
	v    *int
	list []intOrSlice
}

func (is intOrSlice) equals(right intOrSlice) bool {
	if is.v != nil && right.v != nil {
		return *is.v == *right.v
	}
	if is.v != nil || right.v != nil {
		return false
	}
	if len(is.list) != len(right.list) {
		return false
	}
	for i, left := range is.list {
		if !left.equals(right.list[i]) {
			return false
		}
	}
	return true
}

func (is intOrSlice) string() string {
	if is.v != nil {
		return fmt.Sprintf("%d", *is.v)
	}
	members := mapValue(is.list, func(listVal intOrSlice) string { return listVal.string() })
	return fmt.Sprintf("[%s]", strings.Join(members, ","))
}

func (is intOrSlice) getList() []intOrSlice {
	if is.v != nil {
		return []intOrSlice{is}
	}
	return is.list
}

func parseListFromString(s string) intOrSlice {
	var parsedList []intOrSlice
	for i := 0; i < len(s); {
		switch s[i] {
		case ',':
			i++
		case '[':
			openBrackets := 1
			endIndex := -1
			for j := i + 1; openBrackets > 0; j++ {
				switch s[j] {
				case '[':
					openBrackets++
				case ']':
					openBrackets--
				}
				if openBrackets == 0 {
					endIndex = j
				}
			}
			sCopy := strings.Join([]string{"", s}, "")
			innerList := parseListFromString(sCopy[i+1 : endIndex])
			if innerList.list != nil {
				parsedList = append(parsedList, intOrSlice{list: innerList.list})
			} else {
				parsedList = append(parsedList, intOrSlice{list: []intOrSlice{innerList}})
			}
			i = endIndex + 1
		case ']':
			i++
		default:
			// must be a number
			sCopy := strings.Join([]string{"", s}, "")
			endIndex := strings.Index(sCopy[i:], ",")
			var v int
			if endIndex == -1 {
				v = mustInt(sCopy[i:])
				i = len(s)
			} else {
				v = mustInt(sCopy[i : i+endIndex])
				i += endIndex + 1
			}
			parsedList = append(parsedList, intOrSlice{v: &v})
		}
	}
	return intOrSlice{list: parsedList}
}

func compare(left, right intOrSlice) packetCompareResult {
	if shouldLog {
		fmt.Printf("comparing %s against %s\n", left.string(), right.string())
	}
	// if both are integers
	if left.v != nil && right.v != nil {
		if *left.v < *right.v {
			if shouldLog {
				fmt.Println("integers, left less than right - in order")
			}
			return inOrder
		}
		if *left.v == *right.v {
			if shouldLog {
				fmt.Println("integers, left less equals right - indeterminate")
			}
			return indeterminate
		}
		if shouldLog {
			fmt.Println("integers, left greater than right - out of order")
		}
		return outOfOrder
	}
	leftList := left.getList()
	rightList := right.getList()
	maxRight := len(rightList) - 1
	for i := range leftList {
		if i > maxRight {
			if shouldLog {
				fmt.Println("right ran out of items - out of order")
			}
			return outOfOrder
		}
		l, r := leftList[i], rightList[i]
		if res := compare(l, r); res != indeterminate {
			if shouldLog {
				fmt.Printf("list order determined as %s\n", res.string())
			}
			return res
		}
	}
	if len(leftList) < len(rightList) {
		if shouldLog {
			fmt.Println("left ran out of items - in order")
		}
		return inOrder
	}
	return indeterminate
}

func day13(in io.Reader) (int, int) {
	var packets []intOrSlice

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		packets = append(packets, parseListFromString(line[1:len(line)-1]))
	}

	numOrderedPackets := 0
	for i := 0; i < len(packets); i += 2 {
		pairIndex := (i / 2) + 1
		if shouldLog {
			fmt.Printf("\npair index %d\n", pairIndex)
		}
		if res := compare(packets[i], packets[i+1]); res == inOrder {
			numOrderedPackets += pairIndex
		}
	}

	two, six := 2, 6
	firstDivider := intOrSlice{list: []intOrSlice{{list: []intOrSlice{{v: &two}}}}}
	secondDivider := intOrSlice{list: []intOrSlice{{list: []intOrSlice{{v: &six}}}}}
	packets = append(packets, firstDivider, secondDivider)

	sort.SliceStable(packets, func(left, right int) bool {
		res := compare(packets[left], packets[right])
		return res == inOrder
	})

	start, end := 0, 0
	for i, packet := range packets {
		if packet.equals(firstDivider) {
			start = i + 1
		}
		if packet.equals(secondDivider) {
			end = i + 1
		}
	}

	return numOrderedPackets, start * end
}
