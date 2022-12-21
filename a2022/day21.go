package a2022

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

func day21(in io.Reader) (int, int) {
	type monkey struct {
		name, operation string
		val, priority   int
		dependsOn       []string
		resolvedDepends map[string]bool
	}

	monkeys := map[string]*monkey{}
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ": ")
		name := parts[0]
		if val, err := strconv.Atoi(parts[1]); err == nil {
			monkeys[name] = &monkey{name, "", val, 0, nil, nil}
			continue
		}
		parts = strings.Split(parts[1], " ")
		monkeys[name] = &monkey{name, parts[1], 0, 10000, []string{parts[0], parts[2]}, nil}
	}

	var name string
	queue := []string{"root"}
	for len(queue) > 0 {
		name, queue = queue[0], queue[1:]
		m := monkeys[name]
		if m.dependsOn == nil {
			continue
		}
		m.resolvedDepends = map[string]bool{}
		for _, dName := range m.dependsOn {
			target, ok := monkeys[dName]
			if !ok {
				fmt.Printf("target monkey %s doesn't exist\n", dName)
				m.priority = -1
				continue
			}
			m.resolvedDepends[dName] = true
			target.priority = minInt(target.priority, m.priority-1)
			queue = append(queue, dName)
		}
	}

	sortedMonkeys := mapMapValues(monkeys, func(name string, m *monkey) *monkey {
		return m
	})
	sort.SliceStable(sortedMonkeys, func(i, j int) bool {
		return sortedMonkeys[i].priority < sortedMonkeys[j].priority
	})
	for _, m := range sortedMonkeys {
		if m.resolvedDepends == nil {
			continue
		}
		for _, maybeParent := range monkeys {
			if maybeParent.resolvedDepends[m.name] {
				for rd := range m.resolvedDepends {
					maybeParent.resolvedDepends[rd] = true
				}
			}
		}
	}

	opToFunc := map[string]func(string, string) int{
		"+": func(left, right string) int { return monkeys[left].val + monkeys[right].val },
		"-": func(left, right string) int { return monkeys[left].val - monkeys[right].val },
		"*": func(left, right string) int { return monkeys[left].val * monkeys[right].val },
		"/": func(left, right string) int { return monkeys[left].val / monkeys[right].val },
	}
	for _, m := range sortedMonkeys {
		if m.dependsOn == nil || m.priority == -1 {
			continue
		}
		m.val = opToFunc[m.operation](m.dependsOn[0], m.dependsOn[1])
	}

	part1 := monkeys["root"].val

	mustEqual := aElseB(
		!monkeys[monkeys["root"].dependsOn[0]].resolvedDepends["humn"],
		monkeys[monkeys["root"].dependsOn[0]].val,
		monkeys[monkeys["root"].dependsOn[1]].val)

	humnTree := aElseB(
		monkeys[monkeys["root"].dependsOn[0]].resolvedDepends["humn"],
		monkeys[monkeys["root"].dependsOn[0]],
		monkeys[monkeys["root"].dependsOn[1]])

	unwindMath := func(left, right *monkey, op string, result int) int {
		leftDependsOnHumn := left.resolvedDepends["humn"] || left.name == "humn"
		if op == "+" {
			return result - aElseB(leftDependsOnHumn, right.val, left.val)
		}
		if op == "*" {
			return result / aElseB(leftDependsOnHumn, right.val, left.val)
		}
		if leftDependsOnHumn {
			if op == "-" {
				return result + right.val
			} else {
				return result * right.val
			}
		} else {
			if op == "-" {
				return left.val - result
			} else {
				return left.val / result
			}
		}
	}

	part2 := mustEqual
	iter := 0
	for {
		left, right := humnTree.dependsOn[0], humnTree.dependsOn[1]
		part2 = unwindMath(monkeys[left], monkeys[right], humnTree.operation, part2)
		iter++
		if humnTree.dependsOn[0] == "humn" || humnTree.dependsOn[1] == "humn" {
			break
		}
		humnTree = aElseB(monkeys[left].resolvedDepends["humn"], monkeys[left], monkeys[right])
	}

	return part1, part2

}
