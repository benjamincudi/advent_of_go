package a2022

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
)

var allNumbers = regexp.MustCompile(`(\d+)`)

type compositeCost struct {
	ore, other int
}
type blueprint struct {
	id, oreRobotCost, clayRobotCost   int
	obsidianRobotCost, geodeRobotCost compositeCost
}

func day19(in io.Reader) int {
	scanner := bufio.NewScanner(in)
	var blueprints []blueprint
	for scanner.Scan() {
		numbers := mapValue(allNumbers.FindAllString(scanner.Text(), -1), mustInt)
		blueprints = append(blueprints, blueprint{
			numbers[0],
			numbers[1],
			numbers[2],
			compositeCost{numbers[3], numbers[4]},
			compositeCost{numbers[5], numbers[6]}})
	}

	totalQuality := 0
	for _, bp := range blueprints {
		v := mostGeodes(bp, harvestState{0, 0, 0, 0, 1, 0, 0, 0}, 0)
		if shouldLog {
			fmt.Printf("blueprint id %d has max geodes %d, score of %d\n", bp.id, v, bp.id*v)
		}
		totalQuality += bp.id * v
	}

	return totalQuality
}

type harvestState struct {
	ore, clay, obsidian, geodes                        int
	oreRobots, clayRobots, obsidianRobots, geodeRobots int
}

func (hs harvestState) canBuildAfter(bp blueprint) map[buildTarget]int {
	oreForOre := maxInt(bp.oreRobotCost-hs.ore, 0)
	oreForClay := maxInt(bp.clayRobotCost-hs.ore, 0)
	oreForObsidian := maxInt(bp.obsidianRobotCost.ore-hs.ore, 0)
	clayForObsidian := maxInt(bp.obsidianRobotCost.other-hs.clay, 0)
	oreForGeode := maxInt(bp.geodeRobotCost.ore-hs.ore, 0)
	obsidianForGeode := maxInt(bp.geodeRobotCost.other-hs.obsidian, 0)

	// If we already make enough ore each turn to build anything, stop making them
	maxOre := maxInt(bp.oreRobotCost, bp.clayRobotCost, bp.obsidianRobotCost.ore, bp.geodeRobotCost.ore)

	return map[buildTarget]int{
		oreRobot: aElseB(hs.oreRobots < maxOre,
			(oreForOre/hs.oreRobots)+aElseB(oreForOre%hs.oreRobots == 0, 1, 2),
			0,
		),
		clayRobot: aElseB(hs.clayRobots < bp.obsidianRobotCost.other,
			(oreForClay/hs.oreRobots)+aElseB(oreForClay%hs.oreRobots == 0, 1, 2),
			0,
		),
		obsidianRobot: aElseB(hs.clayRobots > 0 && hs.obsidianRobots < bp.geodeRobotCost.other,
			maxInt(
				(oreForObsidian/hs.oreRobots)+aElseB(oreForObsidian%hs.oreRobots == 0, 1, 2),
				(clayForObsidian/maxInt(hs.clayRobots, 1))+aElseB(clayForObsidian%maxInt(hs.clayRobots, 1) == 0, 1, 2)),
			0),
		geodeRobot: aElseB(hs.obsidianRobots > 0,
			maxInt(
				(oreForGeode/hs.oreRobots)+1,
				(obsidianForGeode/maxInt(hs.obsidianRobots, 1))+aElseB(obsidianForGeode%maxInt(hs.obsidianRobots, 1) == 0, 1, 2)),
			0),
	}
}

type buildTarget int

const (
	nothing buildTarget = iota
	oreRobot
	clayRobot
	obsidianRobot
	geodeRobot
)

func (hs harvestState) build(bp blueprint, t buildTarget, minutes int) harvestState {
	nextOre := (hs.oreRobots * minutes) + hs.ore
	nextClay := (hs.clayRobots * minutes) + hs.clay
	nextObsidian := (hs.obsidianRobots * minutes) + hs.obsidian
	nextGeodes := (hs.geodeRobots * minutes) + hs.geodes
	switch t {
	case oreRobot:
		return harvestState{
			nextOre - bp.oreRobotCost,
			nextClay,
			nextObsidian,
			nextGeodes,
			hs.oreRobots + 1,
			hs.clayRobots,
			hs.obsidianRobots,
			hs.geodeRobots,
		}
	case clayRobot:
		return harvestState{
			nextOre - bp.clayRobotCost,
			nextClay,
			nextObsidian,
			nextGeodes,
			hs.oreRobots,
			hs.clayRobots + 1,
			hs.obsidianRobots,
			hs.geodeRobots,
		}
	case obsidianRobot:
		return harvestState{
			nextOre - bp.obsidianRobotCost.ore,
			nextClay - bp.obsidianRobotCost.other,
			nextObsidian,
			nextGeodes,
			hs.oreRobots,
			hs.clayRobots,
			hs.obsidianRobots + 1,
			hs.geodeRobots,
		}
	case geodeRobot:
		return harvestState{
			nextOre - bp.geodeRobotCost.ore,
			nextClay,
			nextObsidian - bp.geodeRobotCost.other,
			nextGeodes,
			hs.oreRobots,
			hs.clayRobots,
			hs.obsidianRobots,
			hs.geodeRobots + 1,
		}
	case nothing:
		return harvestState{
			nextOre,
			nextClay,
			nextObsidian,
			nextGeodes,
			hs.oreRobots,
			hs.clayRobots,
			hs.obsidianRobots,
			hs.geodeRobots,
		}
	default:
		panic("unknown build target")
	}
}

func mostGeodes(bp blueprint, hs harvestState, minute int) int {
	if minute == 24 {
		return hs.geodes
	}
	remainingTime := 24 - minute
	buildOptions := hs.canBuildAfter(bp)
	best := 0
	allowDoNothing := minInt(mapMapValues(buildOptions, func(t buildTarget, time int) int { return time })...) > remainingTime
	for t, canBuildAfter := range buildOptions {
		// We don't have a necessary gathering robot yet
		// Or it would complete after time is up
		if canBuildAfter < 1 || minute+canBuildAfter > 24 {
			continue
		}
		best = maxInt(best, mostGeodes(bp, hs.build(bp, t, canBuildAfter), minute+canBuildAfter))
	}
	if allowDoNothing {
		best = maxInt(best, mostGeodes(bp, hs.build(bp, nothing, remainingTime), minute+remainingTime))
	}
	return best
}
