package main

import (
	"fmt"
	"sort"

	"github.com/adsmf/adventofcode2019/utils"
	"github.com/adsmf/adventofcode2019/utils/pathfinding/astar"
)

func main() {
	fmt.Printf("Part 1: %d\n", part1())
	// fmt.Printf("Part 2: %d\n", part2())
	fmt.Printf("Part 2: 0\n")
}

func part1() int {
	crawl := loadMap("input.txt")
	// fmt.Printf("Crawler:\n%v\n", crawl)
	return crawl.collectKeys()
}

func part2() int {
	return 0
}

type crawler struct {
	vault vault
	keys  map[int]*key
	start point
}

func (c *crawler) collectKeys() int {
	keyIDs := []int{}
	// fmt.Print("Collecting keys: ")
	for id := range c.keys {
		// fmt.Printf("%c ", id)
		keyIDs = append(keyIDs, id)
	}
	// fmt.Println()
	type option struct {
		keys      []int
		lastPos   point
		lastSteps int
	}

	options := []option{}
	for _, id := range keyIDs {
		options = append(options, option{
			keys:      []int{id},
			lastPos:   c.start,
			lastSteps: 0,
			// remaining: remaining,
		})
	}

	bestSteps := utils.MaxInt
	lastLen := 0
	for len(options) > 0 {
		newOptions := []option{}

		iterOptions := append(options[0:0], options...)
		for _, opt := range iterOptions {
			optSteps, endPoint := c.tryOrder(opt.keys, opt.lastPos)
			if optSteps == utils.MaxInt {
				continue
			}

			if lastLen < len(opt.keys) {
				lastLen = len(opt.keys)
				fmt.Printf("Starting length %d\n", lastLen)
			}

			optSteps += opt.lastSteps
			if len(opt.keys) == len(c.keys) {
				if optSteps < bestSteps {
					bestSteps = optSteps
				}
				continue
			}

			remaining := without(c.keys, opt.keys)
			// remainingIter := append(remaining[0:0], remaining...)
			for _, next := range remaining {
				keysCopy := append(opt.keys[0:0], opt.keys...)
				for _, check := range keysCopy {
					if check == next {
						panic("Should not have next initial list")
					}
				}
				keysCopy = append(keysCopy, next)
				newWO := without(c.keys, keysCopy)
				for _, check := range newWO {
					if check == next {
						panic("Should not have next in remaining")
					}
				}
				if len(c.keys) != len(newWO)+len(keysCopy) {
					panic("Wrong number of keys")
				}
				newOpt := option{
					keys:      keysCopy,
					lastPos:   endPoint,
					lastSteps: optSteps,
				}
				newOptions = append(newOptions, newOpt)
			}
		}
		// options := []option{}
		// for _, option:=
		options = newOptions
	}
	return bestSteps
}

func without(keys map[int]*key, skipList []int) []int {
	filtered := []int{}
	skipMap := map[int]bool{}
	for _, skip := range skipList {
		skipMap[skip] = true
	}

	for item := range keys {
		if _, found := skipMap[item]; !found {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func (c *crawler) tryOrder(keys []int, startPoint point) (int, point) {
	pos := startPoint
	keyIndex := len(keys) - 1
	wantKey := keys[keyIndex]

	for id := range c.vault.openDoors {
		c.vault.openDoors[id] = false
	}
	for doorIdx := 0; doorIdx < keyIndex; doorIdx++ {
		key := keys[doorIdx]
		c.vault.openDoors[key] = true
	}
	keyPos := c.keys[wantKey].pos
	routeSteps := c.getKey(pos, keyPos)
	if routeSteps == utils.MaxInt {
		return utils.MaxInt, startPoint
	}
	routeSteps--
	pos = keyPos
	return routeSteps, pos
}

func (c *crawler) getKey(start, keyPos point) int {
	startNode := c.vault.vault[start]
	keyNode := c.vault.vault[keyPos]
	route, err := astar.Route(startNode, keyNode)
	if err != nil {
		return utils.MaxInt
	}
	return len(route)
}

func (c crawler) String() string {
	retString := "Keys:\n"
	keyIDs := []int{}
	for id := range c.keys {
		keyIDs = append(keyIDs, id)
	}
	sort.Ints(keyIDs)
	for _, idInt := range keyIDs {
		id := idInt
		var sym rune
		if c.keys[id].held {
			sym = '🔑'
		} else {
			sym = '❌'
		}
		retString += fmt.Sprintf("%c: %c ", id, sym)
	}
	retString += fmt.Sprintf("\nVault:\n%v\n", c.vault)
	return retString
}

// type routeNode struct {
// 	vault *vault
// 	pos   point
// }

func (t tile) Heuristic(from astar.Node) astar.Cost {
	fromNode := from.(tile)
	xDiff := t.pos.x - fromNode.pos.x
	if xDiff < 0 {
		xDiff *= -1
	}
	yDiff := t.pos.y - fromNode.pos.y
	if yDiff < 0 {
		yDiff *= -1
	}
	return astar.Cost(xDiff+yDiff) + 1
}

func (t tile) Paths() []astar.Edge {
	edges := []astar.Edge{}
	dirs := []point{
		point{t.pos.x, t.pos.y - 1},
		point{t.pos.x, t.pos.y + 1},
		point{t.pos.x - 1, t.pos.y},
		point{t.pos.x + 1, t.pos.y},
	}
	for _, dir := range dirs {
		if t.canTraverse(t.vault.vault[dir]) {
			edges = append(edges, astar.Edge{
				To:   t.vault.vault[dir],
				Cost: 1,
			})
		}
	}
	return edges
}

func (t tile) canTraverse(target tile) bool {
	switch target.tileType {
	case tileTypeEmpty:
		return true
	case tileTypeWall:
		return false
	case tileTypeKey:
		return true
	case tileTypeStart:
		return true
	case tileTypeDoor:
		open := t.vault.doorOpen(target.id)
		return open
	default:
		return false
	}
}

type key struct {
	pos  point
	held bool
}

type vault struct {
	vault      map[point]tile
	openDoors  map[int]bool
	minX, minY int
	maxX, maxY int
}

func (v vault) doorOpen(door int) bool {
	return v.openDoors[door]
}

func (v *vault) set(pos point, val tile) {
	val.pos = pos
	val.vault = v
	v.vault[pos] = val
	if v.minX > pos.x {
		v.minX = pos.x
	}
	if v.maxX < pos.x {
		v.maxX = pos.x
	}
	if v.minY > pos.y {
		v.minY = pos.y
	}
	if v.maxY < pos.y {
		v.maxY = pos.y
	}
}

func (v vault) String() string {
	newString := ""
	for y := v.minY; y <= v.maxY; y++ {
		for x := v.minX; x <= v.maxX; x++ {
			newString += fmt.Sprintf("%v", v.vault[point{x, y}])
		}
		newString += fmt.Sprintln()
	}
	return newString
}

func loadMap(filename string) crawler {
	newVault := vault{
		vault:     map[point]tile{},
		openDoors: map[int]bool{},
	}
	newCrawler := crawler{
		keys: map[int]*key{},
	}

	lines := utils.ReadInputLines(filename)
	for y, line := range lines {
		for x, char := range line {
			pos := point{x, y}
			switch {
			case char == '@':
				newVault.set(pos, tile{tileType: tileTypeEmpty})
				newCrawler.start = pos
			case char == '.':
				newVault.set(pos, tile{tileType: tileTypeEmpty})
			case char == '#':
				newVault.set(pos, tile{tileType: tileTypeWall})
			case 'a' <= char && char <= 'z':
				id := int(char)
				newKey := key{pos: pos}
				newVault.set(pos, tile{tileType: tileTypeKey, id: id})
				newCrawler.keys[id] = &newKey
			case 'A' <= char && char <= 'Z':
				id := int(char - 'A' + 'a')
				newVault.set(pos, tile{tileType: tileTypeDoor, id: id})
			}
		}
	}

	newCrawler.vault = newVault
	return newCrawler
}

type tile struct {
	tileType tileType
	id       int
	pos      point
	vault    *vault
}

func (t tile) String() string {
	switch t.tileType {
	case tileTypeEmpty:
		//return "⬛️"
		return " "
	case tileTypeWall:
		//return "⬜️"
		return "#"
	case tileTypeDoor:
		//return "🚪"
		return string(t.id - 'a' + 'A')
	case tileTypeKey:
		//return "🔑"
		return string(t.id)
	case tileTypeStart:
		//return "⭐️"
		return "s"
	case tileTypeUnkown:
		return "?"
	default:
		return "!"
		//return "❓"
	}
}

type tileType int

const (
	tileTypeUnkown tileType = iota
	tileTypeEmpty
	tileTypeWall
	tileTypeDoor
	tileTypeKey
	tileTypeStart
)

type point struct {
	x, y int
}
