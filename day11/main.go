package main

import (
	"fmt"
	"sync"
)

func main() {
	fmt.Printf("Part 1: %d\n", part1())
	fmt.Printf("Part 2:\n%s\n", part2())
}

func part1() int {
	inputString := loadInputString()
	hull := runPainter(inputString, 0)
	return hull.countVisited()
}

func part2() string {
	inputString := loadInputString()
	hull := runPainter(inputString, 1)
	registration := decodePaint(hull)
	fmt.Print(hull.print())
	return registration
}

func decodePaint(hull shipHull) string {
	registration := ""
	for letterNum := 0; letterNum < 6; letterNum++ {
		letter := getLetter(hull, letterNum)
		registration += letter.decode()
	}
	return registration
}
func getLetter(hull shipHull, num int) letter {
	l := letter{}
	offset := (5 * num) + 1
	for x := 0; x < 4; x++ {
		for y := 0; y < 6; y++ {
			if hull.painted(x+offset, y) {
				l.set(x, y, true)
			}
		}
	}
	return l
}

type letter boolGrid

func (l letter) set(x, y int, val bool) { boolGrid(l).set(x, y, val) }
func (l letter) value(x, y int) bool    { return boolGrid(l).value(x, y) }

func (l letter) decode() string {
	debug := "> "
	// Handles:
	//  CEFHKPUZ

	type position struct {
		x, y int
	}
	checkChars := []position{
		position{0, 0},
		position{1, 0},
		position{2, 0},
		position{3, 0},
	}
	for _, pos := range checkChars {
		if l.value(pos.x, pos.y) {
			debug += "T"
		} else {
			debug += "F"
		}
	}

	// if !l.value(0, 0) {
	// 	// C
	// 	if l.value(1, 0) {
	// 		return "C"
	// 	}
	// 	return "3"
	// } else {
	// 	// P F K H E Z U
	// 	if l.value(1, 0) {
	// 		// P F E Z
	// 		return "1"
	// 	} else {
	// 		// K H U
	// 		if l.value(1, 2) {
	// 			// K H
	// 			return "2"
	// 		} else {
	// 			return "U"
	// 		}
	// 	}
	// }
	return debug + "\n"
}

type boolGrid map[int]map[int]bool

func (bg boolGrid) set(x, y int, value bool) {
	if bg[x] == nil {
		bg[x] = map[int]bool{}
	}
	bg[x][y] = value
}

func (bg boolGrid) value(x, y int) bool {
	if bg == nil {
		return false
	}
	if bg[x] == nil {
		return false
	}
	return bg[x][y]
}

type shipHull struct {
	paintColour boolGrid
	visited     boolGrid
}

func (h *shipHull) print() string {
	printout := ""
	var minX, maxX int
	var minY, maxY int
	for x, cols := range h.paintColour {
		for y, tile := range cols {
			if tile {
				if x < minX {
					minX = x
				}
				if x > maxX {
					maxX = x
				}
				if y < minY {
					minY = y
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}
	for y := minY; y < maxY+1; y++ {
		for x := minX; x < maxX+1; x++ {
			if h.painted(x, y) {
				printout += fmt.Sprint("#")
			} else {
				printout += fmt.Sprint(".")
			}
		}
		printout += fmt.Sprintln()
	}
	return printout
}

func (h *shipHull) visit(x, y int) {
	if h.visited == nil {
		h.visited = boolGrid{}
	}
	h.visited.set(x, y, true)
}

func (h *shipHull) paint(x, y int, colour bool) {
	h.visit(x, y)
	if h.paintColour == nil {
		h.paintColour = boolGrid{}
	}
	h.paintColour.set(x, y, colour)
}

func (h *shipHull) painted(x, y int) bool {
	return h.paintColour.value(x, y)
}

func (h *shipHull) countVisited() int {
	count := 0
	for _, cols := range h.visited {
		for _, tile := range cols {
			if tile {
				count++
			}
		}
	}
	return count
}

type robot struct {
	x, y   int
	facing facing
}

func (r *robot) turnRight() {
	r.facing = r.facing.right()
}
func (r *robot) turnLeft() {
	r.facing = r.facing.left()
}
func (r *robot) move() {
	switch r.facing {
	case facingUp:
		r.y--
	case facingDown:
		r.y++
	case facingRight:
		r.x++
	case facingLeft:
		r.x--
	}
}

type facing int

const (
	facingUp    facing = 0
	facingRight facing = 1
	facingDown  facing = 2
	facingLeft  facing = 3
)

func (f facing) right() facing {
	return facing((f + 1) % 4)
}

func (f facing) left() facing {
	if f == facingUp {
		return facingLeft
	}
	return facing(f - 1)
}

func runPainter(program string, startingPanel int64) shipHull {
	hull := shipHull{}
	robo := robot{}
	output := make(chan int64)
	input := make(chan int64, 1)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		nextInstructionIsPaint := true
		for op := range output {
			if nextInstructionIsPaint {
				// We're painting
				if op == 1 {
					hull.paint(robo.x, robo.y, true)
				} else {
					hull.paint(robo.x, robo.y, false)
				}
			} else {
				// We're moving
				if op == 1 {
					robo.turnRight()
				} else {
					robo.turnLeft()
				}
				robo.move()
				if hull.painted(robo.x, robo.y) {
					input <- 1
				} else {
					input <- 0
				}
			}
			nextInstructionIsPaint = !nextInstructionIsPaint
		}
		wg.Done()
	}()

	input <- startingPanel
	tape := newMachine(program, input, output)
	tape.run()

	wg.Wait()
	return hull
}
