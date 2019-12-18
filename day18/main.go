package main

import (
	"fmt"
	"sort"

	"github.com/adsmf/adventofcode2019/utils"
	"github.com/adsmf/adventofcode2019/utils/pathfinding/astar"
)

func main() {
	fmt.Printf("Part 1: %d\n", part1())
	fmt.Printf("Part 2: %d\n", part2())
}

func part1() int {
	crawl := loadMap("input.txt")
	fmt.Printf("Crawler:\n%v\n", crawl)
	return crawl.collectKeys()
}

func part2() int {
	return 0
}

type crawler struct {
	vault vault
	keys  map[int]key
	start point
}

func (c *crawler) collectKeys() int {
	steps := 0
	keyIDs := []int{}
	for id := range c.keys {
		keyIDs = append(keyIDs, int(id))
	}

	perms := utils.PermuteInts(keyIDs)
	iterSteps := make([]int, len(perms))
	for iter, keys := range perms {
		iterSteps[iter] = c.tryOrder(keys)
	}
	return steps
}

func (c *crawler) tryOrder(keys []int) int {
	steps := 0
	pos := c.start
	for keyIndex := 0; keyIndex < len(keys); keyIndex++ {
		wantKey := keys[keyIndex]
		var haveKeys []int
		if keyIndex > 0 {
			haveKeys = keys[:keyIndex-1]
		}
		// notKeys := keys[keyIndex:]
		vaultCopy := vault{vault: map[point]tile{}}
		var keyPos point
		for pos, t := range c.vault.vault {
			switch t.tileType {
			case tileTypeDoor:
				haveKey := false
				for _, key := range haveKeys {
					if key == t.id {
						haveKey = true
						break
					}
				}
				if haveKey {
					t.tileType = tileTypeEmpty
				} else {
					t.tileType = tileTypeWall
				}
			case tileTypeKey:
				if t.id == wantKey {
					keyPos = pos
				} else {
					t.tileType = tileTypeEmpty
				}
			}
			if pos == c.start {
				t.tileType = tileTypeStart
			}
			vaultCopy.set(pos, t)
		}
		fmt.Printf("Trying to get key %d (%v)\n", wantKey, keyPos)
		routeSteps := c.getKey(vaultCopy, pos, keyPos)
		if routeSteps == utils.MaxInt {
			fmt.Printf("Couldn't get key\n")
			return utils.MaxInt
		}
		pos = keyPos
		steps += routeSteps
	}
	return steps
}

func (c *crawler) getKey(vaultCopy vault, start, keyPos point) int {
	startNode := routeNode{vault: vaultCopy, pos: start}
	keyNode := routeNode{vault: vaultCopy, pos: keyPos}
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

type routeNode struct {
	vault vault
	pos   point
}

func (r routeNode) Heuristic(from astar.Node) astar.Cost {
	fromNode := from.(routeNode)
	xDiff := r.pos.x - fromNode.pos.x
	if xDiff < 0 {
		xDiff *= -1
	}
	yDiff := r.pos.y - fromNode.pos.y
	if yDiff < 0 {
		yDiff *= -1
	}
	return astar.Cost(xDiff + yDiff)
}

func (r routeNode) Paths() []astar.Edge {
	return nil
}

type key struct {
	pos  point
	held bool
}

type vault struct {
	vault      map[point]tile
	minX, minY int
	maxX, maxY int
}

func (v *vault) set(pos point, val tile) {
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
		vault: map[point]tile{},
	}
	newCrawler := crawler{
		keys: map[int]key{},
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
				newCrawler.keys[id] = newKey
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
}

func (t tile) String() string {
	switch t.tileType {
	case tileTypeEmpty:
		return "⬛️"
	case tileTypeWall:
		return "⬜️"
	case tileTypeDoor:
		return "🚪"
	case tileTypeKey:
		return "🔑"
	case tileTypeStart:
		return "⭐️"
	default:
		return "❓"
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
