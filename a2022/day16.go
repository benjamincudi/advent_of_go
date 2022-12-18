package a2022

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

var lineParser = regexp.MustCompile(`^Valve (?P<ValveName>\w+) has flow rate=(?P<flowRate>\d+); tunnels? leads? to valves? (?P<tunnels>.+)$`)

type pressureValve struct {
	name       string
	flowRate   int // pressure released per remaining minute
	toValves   []string
	isOpen     bool
	valueEdges map[string]int
}

func copyMap[K comparable, V any](in map[K]V) map[K]V {
	ret := map[K]V{}
	for k, v := range in {
		ret[k] = v
	}
	return ret
}

func recursePaths(valveMap map[string]*pressureValve, walkerLocations []string, valvesOpenedFor map[string]int) (int, int) {
	best, checked := 0, 0
	remainingTime := maxInt(mapValue(walkerLocations, func(name string) int { return valvesOpenedFor[name] })...)
	currentWalkerIndex := aElseB(valvesOpenedFor[walkerLocations[0]] == remainingTime, 0, 1)
	for next, time := range valveMap[walkerLocations[currentWalkerIndex]].valueEdges {
		// if the valve was already opened, don't go back
		if _, opened := valvesOpenedFor[next]; opened {
			continue
		}
		// if we can't get there in time, don't consider it
		if time > remainingTime {
			continue
		}
		voat := copyMap(valvesOpenedFor)
		// we can make it, and it is currently closed, so go there and open it
		voat[next] = remainingTime - time
		nextWalkerLocations := append([]string{}, walkerLocations...)
		nextWalkerLocations[currentWalkerIndex] = next
		total, permutations := recursePaths(valveMap, nextWalkerLocations, voat)
		checked += permutations
		best = aElseB(total > best, total, best)
	}
	// if we ran out of time to go anywhere else, calculate our total and return it
	if best == 0 {
		totals := mapMapValues(valvesOpenedFor, func(name string, time int) int {
			return time * valveMap[name].flowRate
		})
		sum := 0
		for _, t := range totals {
			sum += t
		}
		return sum, 1
	}
	// return the best subtree we found
	return best, checked
}

func day16(in io.Reader) (int, int) {
	scanner := bufio.NewScanner(in)

	valveMap := map[string]*pressureValve{}
	for scanner.Scan() {
		t := scanner.Text()
		if lineParser.MatchString(t) {
			subgroups := lineParser.FindStringSubmatch(t)
			name := subgroups[1]
			valveMap[name] = &pressureValve{
				name,
				mustInt(subgroups[2]),
				strings.Split(subgroups[3], ", "),
				false,
				nil}
			//fmt.Printf("%s: %v\n", name, *valveMap[name])
		}
	}

	// Build "shortest path to each valuable valve" so we can skip path construction later
	for _, pv := range valveMap {
		// We'll never intentionally visit 0-value nodes, so don't build a map for them
		// AA is special because it is everyone's starting point
		if pv.flowRate == 0 && pv.name != "AA" {
			continue
		}
		valueEdges := map[string]int{}
		type pathQueue struct {
			name string
			path []string
		}
		seen := map[string]int{pv.name: len(pv.toValves)} // queue might not actually visit in the fastest order, allow N visits
		queue := mapValue(pv.toValves, func(name string) pathQueue { return pathQueue{name, []string{name}} })
		var current pathQueue
		for len(queue) > 0 {
			current, queue = queue[0], queue[1:]
			seen[current.name] += 1
			if valveMap[current.name].flowRate > 0 {
				distance := len(current.path) + 1 // include the valve-opening turn
				existingDistance, ok := valueEdges[current.name]
				valueEdges[current.name] = aElseB(ok, minInt(distance, existingDistance), distance)
			}
			for _, nextName := range valveMap[current.name].toValves {
				if seen[nextName] < len(valveMap[nextName].toValves) {
					queue = append(queue, pathQueue{nextName, append(append([]string{}, current.path...), nextName)})
				}
			}
		}

		pv.valueEdges = valueEdges

		fmt.Printf("%s has edges: %v\n", pv.name, valueEdges)
	}

	bestRelease, permutationsChecked := recursePaths(valveMap, []string{"AA"}, map[string]int{"AA": 30})
	fmt.Printf("checked %d options\n", permutationsChecked)

	withElephant, permutationsChecked := recursePaths(valveMap, []string{"AA", "AA"}, map[string]int{"AA": 26})
	fmt.Printf("checked %d options\n", permutationsChecked)

	// fully naive, no heuristics
	// checked 215953 options (bestRelease) ~0.6s
	// checked 127119358 options
	//    --- PASS: Test_day16/personal_input (559.26s)

	return bestRelease, withElephant
}
