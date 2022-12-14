package a2022

import (
	"encoding/csv"
	"errors"
	"io"
	"strings"

	"github.com/gocarina/gocsv"
)

type coordinates struct{ X, Y int }

func (c *coordinates) UnmarshalString(s string) error {
	parts := strings.Split(s, ",")
	c.X, c.Y = mustInt(parts[0]), mustInt(parts[1])
	return nil
}

// vector represents the amount a ropeKnot will move
// it aliases coordinates for convenience as I am lazy
type vector coordinates

// a basicDirection corresponds to one of four preset vectors
// this makes moving the head knot super simple
var letterToVector = map[string]vector{
	"U": {0, 1}, "L": {-1, 0},
	"D": {0, -1}, "R": {1, 0},
}

// basicDirection adds very minimal safety to ropeInstruction parsing
//
// input is very constraints so this doesn't really guard against much,
// but it does let us translate to a vector on the way in
type basicDirection struct{ vec vector }

func (d *basicDirection) UnmarshalCSV(s string) error {
	if vec, ok := letterToVector[s]; !ok {
		return errors.New("unknown direction")
	} else {
		d.vec = vec
	}
	return nil
}

type ropeInstruction struct {
	D     basicDirection
	Count int
}

type ropeKnot struct {
	position coordinates
	history  map[coordinates]bool
	tail     *ropeKnot
}

func makeRopeKnot(i, count int) *ropeKnot {
	if i == count {
		return nil
	}
	return &ropeKnot{coordinates{0, 0}, map[coordinates]bool{coordinates{0, 0}: true}, makeRopeKnot(i+1, count)}
}

func (r *ropeKnot) move(d basicDirection) {
	r.position = coordinates{r.position.X + d.vec.X, r.position.Y + d.vec.Y}
	r.history[r.position] = true
	r.tail.follow(r.position)
}

// knotPair only used to disambiguate head and tail inputs for getMoveVector
type knotPair struct{ head, tail coordinates }

func getMoveVector(kp knotPair) vector {
	deltaX, deltaY := kp.head.X-kp.tail.X, kp.head.Y-kp.tail.Y
	if abs(deltaX) == 2 || abs(deltaY) == 2 {
		return vector{sign(deltaX), sign(deltaY)}
	}
	return vector{0, 0}
}

func (r *ropeKnot) follow(head coordinates) {
	moveVec := getMoveVector(knotPair{head, r.position})
	r.position = coordinates{r.position.X + moveVec.X, r.position.Y + moveVec.Y}
	r.history[r.position] = true
	if r.tail != nil {
		r.tail.follow(r.position)
	}
}

func (r *ropeKnot) getTailHistoryLength() int {
	if r.tail != nil {
		return r.tail.getTailHistoryLength()
	}
	return len(r.history)
}

func day9(in io.Reader) (int, int) {
	r := csv.NewReader(in)
	r.Comma = ' '
	var ropeSteps []ropeInstruction
	if err := gocsv.UnmarshalCSVWithoutHeaders(r, &ropeSteps); err != nil {
		panic(err)
	}

	head := makeRopeKnot(0, 2)
	h := makeRopeKnot(0, 10)
	for _, s := range ropeSteps {
		for c := 0; c < s.Count; c++ {
			head.move(s.D)
			h.move(s.D)
		}
	}

	return head.getTailHistoryLength(), h.getTailHistoryLength()
}
