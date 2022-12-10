package a2022

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/gocarina/gocsv"
)

type crtProgram struct {
	x, cycle, sigSum int
	drawnPixels      []string
	latest40         int // debugging printing
}

func getPixel(x, pos int) string {
	if abs(x-pos) < 2 {
		return "#"
	}
	return "."
}

func (n *crtProgram) oneCycle() {
	// during cycle - draw a pixel
	// this has off-by-one shenanigans with x/screen being 0-index
	// but the cycle being effectively 1-indexed, which we can
	// handle most cleanly by printing before increment
	n.drawnPixels = append(n.drawnPixels, getPixel(n.x, n.cycle%40))
	n.printCurrentCRTRow()

	n.cycle++

	// during cycle - check strengths every 20mod40
	// before any addX modification can occur, since that happens "at the end"
	if n.cycle%40 == 20 {
		n.sigSum += n.signalStrength()
	}
}
func (n *crtProgram) addX(i int) {
	n.x += i
	n.printSprite()
}
func (n *crtProgram) printSprite() {
	if shouldLog {
		var spriteRow []string
		for j := 0; j < 40; j++ {
			spriteRow = append(spriteRow, getPixel(n.x, j))
		}
		fmt.Printf("sprite position: %s\n", strings.Join(spriteRow, ""))
		fmt.Printf("crtRows\n")
	}
}
func (n *crtProgram) printCurrentCRTRow() {
	if shouldLog {
		crtLoc := n.cycle % 40
		if crtLoc == 0 {
			n.latest40 = n.cycle
		}
		fmt.Printf("crt draws in position %d\n", n.cycle)
		fmt.Printf("%s\n", strings.Join(n.drawnPixels[n.latest40:], ""))
	}
}
func (n *crtProgram) signalStrength() int {
	return n.x * n.cycle
}

type crtProgramInstruction struct {
	cmd string
	val int
}

func (np *crtProgramInstruction) UnmarshalCSV(s string) error {
	parts := strings.Split(s, " ")
	np.cmd = parts[0]
	if parts[0] == "addx" {
		np.val = mustInt(parts[1])
	}
	return nil
}

func day10(in io.Reader) (int, string, error) {
	var ins []struct{ I crtProgramInstruction }
	if err := gocsv.UnmarshalCSVWithoutHeaders(csv.NewReader(in), &ins); err != nil {
		return 0, "", err
	}

	p := crtProgram{1, 0, 0, []string{}, 0}
	p.printSprite()
	for _, pi := range ins {
		if shouldLog {
			fmt.Printf("start cycle %d: begin %v", p.cycle, pi)
		}
		p.oneCycle()
		if pi.I.cmd == "addx" {
			p.oneCycle()
			p.addX(pi.I.val)
		}
	}
	screen := "\n"
	for i := 0; i < len(p.drawnPixels); i += 40 {
		screen += fmt.Sprintf("%s\n", strings.Join(p.drawnPixels[i:i+40], ""))
	}
	return p.sigSum, screen, nil
}
