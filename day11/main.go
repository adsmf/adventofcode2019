package main

import (
	"fmt"
	"sync"
)

func main() {
	fmt.Printf("Part 1: %d\n", part1())
	fmt.Printf("Part 2: %s\n", part2())
}

func part1() int {
	inputString := loadInputString()
	hull := runPainter(inputString, 0)
	return len(hull.visited) //hull.countVisited()
}

func part2() string {
	inputString := loadInputString()
	hull := runPainter(inputString, 1)
	return hull.decode()
}

type point struct {
	x, y int
}
type boolGrid map[point]bool

type shipHull struct {
	paintColour boolGrid
	visited     boolGrid
	min, max    point
}

func (h *shipHull) decode() string {
	return ""
}

func (h *shipHull) print() string {
	printout := ""
	for y := h.min.y; y < h.max.y+1; y++ {
		for x := h.min.x; x < h.max.x+1; x++ {
			if h.paintColour[point{x, y}] {
				printout += fmt.Sprint("█")
			} else {
				printout += fmt.Sprint(" ")
			}
		}
		printout += fmt.Sprintln()
	}
	return printout
}

func (h *shipHull) paint(pos point, colour bool) {
	h.visited[pos] = true
	h.paintColour[pos] = colour

	if pos.x < h.min.x {
		h.min.x = pos.x
	}
	if pos.y < h.min.y {
		h.min.y = pos.y
	}
	if pos.x > h.max.x {
		h.max.x = pos.x
	}
	if pos.y > h.max.y {
		h.max.y = pos.y
	}
}

type robot struct {
	pos    point
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
		r.pos.y--
	case facingDown:
		r.pos.y++
	case facingRight:
		r.pos.x++
	case facingLeft:
		r.pos.x--
	}
}

type facing int

const (
	facingUp facing = iota
	facingRight
	facingDown
	facingLeft
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
	hull := shipHull{
		paintColour: boolGrid{},
		visited:     boolGrid{},
	}
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
					hull.paint(robo.pos, true)
				} else {
					hull.paint(robo.pos, false)
				}
			} else {
				// We're moving
				if op == 1 {
					robo.turnRight()
				} else {
					robo.turnLeft()
				}
				robo.move()
				if hull.paintColour[robo.pos] {
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
