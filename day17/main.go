package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/adsmf/adventofcode2019/utils/intcode"
)

func main() {
	fmt.Printf("Part 1: %d\n", part1())
	fmt.Printf("Part 2: %d\n", part2())
}

func part1() int {
	s := readCamera()
	return s.sumAlignment()
}

func part2() int {
	s := readCamera()
	route := s.generateRoute()
	commands := s.generateCommands(route)
	fmt.Printf("Route: %s\n", route)
	fmt.Printf("Commands: %#v\n", commands)
	return testDustCount(commands)
}

func part2manualCommands() []string {
	return []string{
		"A,B,A,B,C,C,B,A,B,C\n",
		"L,8,R,12,R,12,R,10\n",
		"R,10,R,12,R,10\n",
		"L,10,R,10,L,6\n",
		"n\n",
	}
}

func part2manual() int {
	cmds := part2manualCommands()
	return testDustCount(cmds)
}

func testDustCount(commands []string) int {
	s := scaffold{
		commands: commands,
	}
	input := loadInputString()
	m := intcode.NewMachine(intcode.M19(s.driver, s.outputHandler))
	m.LoadProgram(input)
	m.WriteRAM(0, 2)
	m.Run(false)
	return s.collected

}

func readCamera() scaffold {
	s := scaffold{}
	input := loadInputString()
	m := intcode.NewMachine(intcode.M19(nil, s.outputHandler))
	m.LoadProgram(input)
	m.Run(false)
	s.processCamera()
	return s
}

type facing int

const (
	facingUp facing = iota
	facingRight
	facingDown
	facingLeft
)

func (f facing) forward() vector {
	switch f {
	case facingUp:
		return vector{0, -1}
	case facingDown:
		return vector{0, 1}
	case facingLeft:
		return vector{-1, 0}
	case facingRight:
		return vector{1, 0}
	}
	return vector{}
}
func (f facing) left() facing {
	newFacing := f - 1
	if newFacing < 0 {
		newFacing += 4
	}
	return newFacing
}
func (f facing) right() facing {
	newFacing := (f + 1) % 4
	return newFacing
}

// func (f facing) right() vector {}

type scaffold struct {
	cameraViewRaw  string
	grid           map[point]scaffoldTile
	intersections  map[point]intersection
	commands       []string
	currentCommand string
	collected      int
	start          point
	facing         facing
}

func (s *scaffold) generateCommands(route string) []string {
	cmds := []string{}

	type pair struct {
		dir string
		len string
	}

	pairs := []pair{}
	parts := strings.Split(route, ",")
	for _ := range pairs {

	}
	mainloop := ""
	functions := map[string]string{}

	cmds = append(cmds, "n\n")
	return cmds
}

func (s *scaffold) generateRoute() string {
	route := ""
	count := 0
	pos := s.start
	dir := s.facing

	for {
		// fmt.Printf("Pos: %v\n", pos)
		fwdPos := pos.add(dir.forward())
		if _, found := s.grid[fwdPos]; found {
			count++
			pos = fwdPos
		} else {
			if count > 0 {
				route += fmt.Sprintf("%d,", count)
				count = 0
			}
			altFound := false
			leftPos := pos.add(dir.left().forward())
			if _, found := s.grid[leftPos]; found {
				route += "L,"
				dir = dir.left()
				altFound = true
			}
			rightPos := pos.add(dir.right().forward())
			if _, found := s.grid[rightPos]; found {
				route += "R,"
				dir = dir.right()
				altFound = true
			}
			if !altFound {
				break
			}
		}
	}
	return strings.Trim(route, ",")
}

func (s *scaffold) driver() (int, bool) {
	// fmt.Printf("Input called!\n")
	if s.currentCommand != "" {
		ch := s.currentCommand[0]
		s.currentCommand = s.currentCommand[1:]
		return int(ch), false
	}
	if len(s.commands) > 0 {
		s.currentCommand = s.commands[0]
		if len(s.commands) > 1 {
			s.commands = s.commands[1:]
		} else {
			s.commands = s.commands[0:0]
		}
		ch := s.currentCommand[0]
		s.currentCommand = s.currentCommand[1:]
		return int(ch), false
	}
	return 0, true
}

func (s *scaffold) outputHandler(out int) {
	if out < 255 {
		s.cameraViewRaw += fmt.Sprintf("%c", out)
	} else {
		s.collected = out
	}
}

func (s *scaffold) sumAlignment() int {
	sum := 0
	for p := range s.intersections {
		sum += p.x * p.y
	}
	return sum
}

func (s *scaffold) processCamera() {
	s.intersections = map[point]intersection{}
	s.grid = map[point]scaffoldTile{}

	lines := strings.Split(s.cameraViewRaw, "\n")
	for y, line := range lines {
		for x, inp := range line {
			pos := point{x, y}
			switch inp {
			case '#':
				s.grid[pos] = scaffoldTile{}
			case '^':
				s.grid[pos] = scaffoldTile{}
				s.start = pos
				s.facing = facingUp
			}
		}
	}
	for p := range s.grid {
		_, up := s.grid[point{p.x, p.y - 1}]
		_, down := s.grid[point{p.x, p.y + 1}]
		_, left := s.grid[point{p.x - 1, p.y}]
		_, right := s.grid[point{p.x + 1, p.y}]
		if up && down && left && right {
			s.intersections[p] = intersection{}
		}
	}
}

type point struct {
	x, y int
}

func (p point) add(v vector) point {
	return point{
		x: p.x + v.x,
		y: p.y + v.y,
	}
}

type vector struct {
	x, y int
}

type intersection struct{}
type scaffoldTile struct{}

func loadInputString() string {
	inputRaw, err := ioutil.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	return string(inputRaw)

}
