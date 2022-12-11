package a2022

import (
	"bufio"
	"io"
	"sort"
	"strings"
)

type monkey struct {
	items               []int
	inspect             func(int) int
	testMod             int
	testTrue, testFalse int
	inspectedCount      int
}

func (m *monkey) addItem(worry int) {
	m.items = append(m.items, worry)
}

func getMonkeyBusinessScore(in []*monkey) int {
	counts := mapValue(in, func(m *monkey) int { return m.inspectedCount })
	sort.Sort(sort.Reverse(sort.IntSlice(counts)))
	return counts[0] * counts[1]
}

func day11(in io.Reader) (int, int) {
	reader := bufio.NewReader(in)

	var monkeys []*monkey
	currentMonkey := &monkey{}
	safeMod := 1
	for s, err := reader.ReadString('\n'); err != io.EOF; s, err = reader.ReadString('\n') {
		if len(strings.TrimSpace(s)) == 0 {
			monkeys = append(monkeys, currentMonkey)
			currentMonkey = &monkey{}
		}
		parts := strings.Split(strings.TrimSpace(s), " ")
		switch parts[0] {
		case "Monkey":
			// do nothing - only append on separators so 0-index is right
		case "Starting":
			items := parts[2:]
			currentMonkey.items = mapValue(items, func(v string) int { return mustInt(strings.TrimRight(v, ",")) })
		case "Operation:":
			operation, value := parts[4], parts[5]
			switch {
			case value == "old" && operation == "*":
				currentMonkey.inspect = func(v int) int { return v * v }
			case value == "old" && operation == "+":
				// this case didn't show up in either input, but it felt weird not to allow it
				currentMonkey.inspect = func(v int) int { return v + v }
			case operation == "*":
				vInt := mustInt(value)
				currentMonkey.inspect = func(v int) int { return v * vInt }
			case operation == "+":
				vInt := mustInt(value)
				currentMonkey.inspect = func(v int) int { return v + vInt }
			}
		case "Test:":
			mod := parts[len(parts)-1]
			currentMonkey.testMod = mustInt(mod)
			safeMod *= currentMonkey.testMod
		case "If":
			testBool, monkeyIndex := parts[1], mustInt(parts[len(parts)-1])
			if testBool == "true:" {
				currentMonkey.testTrue = monkeyIndex
			} else {
				currentMonkey.testFalse = monkeyIndex
			}
		}
	}

	// items and inspectedCount are the only things that need to roll back between problem parts
	backupItems := make([][]int, 0, len(monkeys))
	for _, m := range monkeys {
		backupItems = append(backupItems, append([]int{}, m.items...))
	}

	for i := 0; i < 20; i++ {
		for _, m := range monkeys {
			m.inspectedCount += len(m.items)
			for _, item := range m.items {
				newWorry := m.inspect(item) / 3
				if newWorry%m.testMod == 0 {
					monkeys[m.testTrue].addItem(newWorry)
				} else {
					monkeys[m.testFalse].addItem(newWorry)
				}
			}
			m.items = []int{}
		}
	}

	part1 := getMonkeyBusinessScore(monkeys)

	for i, m := range monkeys {
		m.items, m.inspectedCount = backupItems[i], 0
	}

	for i := 0; i < 10000; i++ {
		for _, m := range monkeys {
			m.inspectedCount += len(m.items)
			for _, item := range m.items {
				// it'd be nice to have a proof of why this is the right choice
				// in lieu of that, it was the only idea that led to the control case passing
				newWorry := m.inspect(item) % safeMod
				if newWorry%m.testMod == 0 {
					monkeys[m.testTrue].addItem(newWorry)
				} else {
					monkeys[m.testFalse].addItem(newWorry)
				}
			}
			m.items = []int{}
		}
	}

	return part1, getMonkeyBusinessScore(monkeys)
}
