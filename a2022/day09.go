package a2022

import (
	"encoding/csv"
	"errors"
	"image"
	"io"

	"github.com/gocarina/gocsv"
)

// a basicDirection corresponds to one of four preset vectors
// this makes moving the head knot super simple
var letterToVector = map[string]image.Point{
	"U": image.Pt(0, 1), "L": image.Pt(-1, 0),
	"D": image.Pt(0, -1), "R": image.Pt(1, 0),
}

// basicDirection adds very minimal safety to ropeInstruction parsing
//
// input is very constraints so this doesn't really guard against much,
// but it does let us translate to a vector on the way in
type basicDirection struct{ vec image.Point }

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
	position image.Point
	history  map[image.Point]bool
	tail     *ropeKnot
}

func makeRopeKnot(i, count int) *ropeKnot {
	if i == count {
		return nil
	}
	return &ropeKnot{image.Pt(0, 0), map[image.Point]bool{image.Pt(0, 0): true}, makeRopeKnot(i+1, count)}
}

func (r *ropeKnot) move(d basicDirection) {
	r.position = image.Pt(r.position.X+d.vec.X, r.position.Y+d.vec.Y)
	r.history[r.position] = true
	r.tail.follow(r.position)
}

// knotPair only used to disambiguate head and tail inputs for getMoveVector
type knotPair struct{ head, tail image.Point }

func getMoveVector(kp knotPair) image.Point {
	deltaX, deltaY := kp.head.X-kp.tail.X, kp.head.Y-kp.tail.Y
	if abs(deltaX) == 2 || abs(deltaY) == 2 {
		return image.Pt(sign(deltaX), sign(deltaY))
	}
	return image.Pt(0, 0)
}

func (r *ropeKnot) follow(head image.Point) {
	moveVec := getMoveVector(knotPair{head, r.position})
	r.position = image.Pt(r.position.X+moveVec.X, r.position.Y+moveVec.Y)
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
