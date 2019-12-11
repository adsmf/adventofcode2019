package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func main() {
	baseDir := "./letters"
	files, err := ioutil.ReadDir(baseDir)
	if err != nil {
		log.Fatal(err)
	}

	letters := letterMap{}
	for _, f := range files {
		// fmt.Print(f.Name(), "  ")
		// fmt.Println(loadPoints(baseDir + "/" + f.Name()))
		letters[f.Name()] = loadPoints(baseDir + "/" + f.Name())
	}
	findUnique(letters)
}

type letterMap map[string]pointMap
type pointUsage map[point]map[string]bool

type letterTest struct {
	sample []point
	letter string
}

func findUnique(letters letterMap) []letterTest {
	usage := pointUsage{}
	for letter, points := range letters {
		for point := range points {
			if usage[point] == nil {
				usage[point] = map[string]bool{}
			}
			usage[point][letter] = true
		}
	}
	samples := genSamples(usage)
	// samples := []letterTest{}
	// for len(usage) > 0 {
	// 	var found bool
	// 	var letter string

	// 	var usePoint point
	// 	found, usePoint, letter = findSingleUsePoint(usage)

	// 	if found {
	// 		samples = append(samples, letterTest{[]point{usePoint}, letter})
	// 	} else {
	// 		smallest := findSmallestPoint(usage)
	// 		fmt.Printf("Smallest: %v\n", smallest)
	// 		// pickLetter := ""
	// 		// for let := range usage[smallest] {
	// 		// 	pickLetter = let
	// 		// }
	// 		// withoutLetter := usageWithoutLetter(usage, pickLetter)
	// 		// fmt.Printf("Without %s: %v\n", pickLetter, withoutLetter)
	// 		// newTests := findUnique(withoutLetter)
	// 	}

	// 	if !found {
	// 		break
	// 	}
	// }

	for _, sample := range samples {
		fmt.Printf("If %v => %s\n", sample.sample, sample.letter)
	}
	if len(usage) > 0 {
		fmt.Printf("\nUnhandled points:\n")
		for point, inLetters := range usage {
			fmt.Println(point, inLetters)
		}
	}
	return samples
}

func genSamples(usage pointUsage) []letterTest {
	samples := []letterTest{}
	for len(usage) > 0 {
		var found bool
		var letter string

		var usePoint point
		found, usePoint, letter = findSingleUsePoint(usage)

		if found {
			samples = append(samples, letterTest{[]point{usePoint}, letter})
		} else {
			smallest := findSmallestPoint(usage)
			fmt.Printf("Smallest: %v\n", smallest)
			pickLetter := ""
			for let := range usage[smallest] {
				pickLetter = let
			}
			withoutLetter := usageWithoutLetter(usage, pickLetter)
			fmt.Printf("Without %s: %v\n", pickLetter, withoutLetter)
			newTests := genSamples(withoutLetter)
			fmt.Printf("newTests: %v\n", newTests)
		}

		if !found {
			break
		}
	}
	return samples
}

func usageWithoutLetter(usage pointUsage, letter string) pointUsage {
	newUsage := pointUsage{}
	for p, pLetters := range usage {
		newPLetters := map[string]bool{}
		for l := range pLetters {
			if l != letter {
				newPLetters[l] = true
			}
		}
		if len(newPLetters) > 0 {
			newUsage[p] = newPLetters
		}
	}
	return newUsage
}

func findSingleUsePoint(usage pointUsage) (bool, point, string) {
	for tryPoint, inLetters := range usage {
		if len(inLetters) == 1 {
			isLetter := ""
			for let := range inLetters {
				isLetter = let
			}

			delete(usage, tryPoint)
			for point := range usage {
				delete(usage[point], isLetter)
			}
			return true, tryPoint, isLetter
		}
	}
	return false, point{}, ""
}

func findSmallestPoint(usage pointUsage) point {
	var bestPoint point

	bestPointLen := 26

	for tryPoint, inLetters := range usage {
		if len(inLetters) < bestPointLen {
			bestPoint = tryPoint
			bestPointLen = len(inLetters)
		}
	}

	return bestPoint
}

type point struct {
	x, y int
}

type pointMap map[point]bool

func loadPoints(filename string) pointMap {
	points := pointMap{}
	bytes, _ := ioutil.ReadFile(filename)
	lines := strings.Split(string(bytes), "\n")
	for y, line := range lines {
		for x, char := range line {
			if char == '#' {
				points[point{x, y}] = true
			}
		}
	}
	return points
}
