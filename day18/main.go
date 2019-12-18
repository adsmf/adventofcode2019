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
	for id := range c.keys {
		keyIDs = append(keyIDs, int(id))
	}

	perms := utils.PermuteInts(keyIDs)
	iterSteps := make([]int, len(perms))
	bestSteps := utils.MaxInt
	fmt.Println("Permutations: ", len(perms))
	for iter, keys := range perms {
		fmt.Print(".")
		steps := c.tryOrder(keys)
		iterSteps[iter] = steps
		if steps < bestSteps {
			bestSteps = steps
		}
	}
	return bestSteps
}

func (c *crawler) tryOrder(keys []int) int {
	// fmt.Printf("Trying permutation: %v\n", keys)
	steps := 0
	pos := c.start
	for keyIndex := 0; keyIndex < len(keys); keyIndex++ {
		// fmt.Printf("Key index: %d\n", keyIndex)
		wantKey := keys[keyIndex]

		for doorIdx := 0; doorIdx < len(keys); doorIdx++ {
			if doorIdx < keyIndex {
				c.vault.openDoors[keys[doorIdx]] = true
			} else {
				c.vault.openDoors[keys[doorIdx]] = false
			}
		}
		keyPos := c.keys[wantKey].pos
		// fmt.Printf("Trying to get key %d (%v) from %v with keys %v\n", wantKey, keyPos, pos, c.vault.openDoors)
		// fmt.Printf("%v", c.vault)
		routeSteps := c.getKey(c.vault, pos, keyPos)
		if routeSteps == utils.MaxInt {
			// fmt.Printf("Couldn't get key\n")
			return utils.MaxInt
		}
		routeSteps--
		// fmt.Printf("Got key: %v\n", routeSteps)
		pos = keyPos
		steps += routeSteps
		// fmt.Println()
	}
	fmt.Printf("Total steps: %d\n\n", steps)
	return steps
}

func (c *crawler) getKey(vaultCopy vault, start, keyPos point) int {
	startNode := vaultCopy.vault[start]
	keyNode := vaultCopy.vault[keyPos]
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
		keyIDs = append(keyIDs, int(id))
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
	return astar.Cost(xDiff + yDiff)
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
	default:
		return "?"
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
