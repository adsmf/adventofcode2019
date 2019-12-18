package main

import (
	"fmt"
	"testing"

	"github.com/adsmf/adventofcode2019/utils"
	"github.com/stretchr/testify/assert"
)

func TestPart1Examples(t *testing.T) {
	type testDef struct {
		file  string
		steps int
	}
	tests := []testDef{
		testDef{
			"amf1.txt",
			2,
		},
		testDef{
			"amf2.txt",
			8,
		},
		testDef{
			file:  "p1ex0.txt",
			steps: 8,
		},
		testDef{
			file:  "p1ex1.txt",
			steps: 86,
		},
		testDef{
			file:  "p1ex2.txt",
			steps: 132,
		},
		// testDef{
		// 	file:  "p1ex3.txt",
		// 	steps: 136,
		// },
		// testDef{
		// 	file:  "p1ex4.txt",
		// 	steps: 81,
		// },
	}
	for id, test := range tests {
		file := "examples/" + test.file
		expected := test.steps
		for i := 0; i < 1; i++ {
			t.Run(fmt.Sprintf("Part1 Test%d Iter%d", i, id), func(t *testing.T) {
				crawl := loadMap(file)
				t.Logf("Start:\n%v", crawl)
				steps := crawl.collectKeys()
				assert.Greater(t, utils.MaxInt, steps)
				assert.Equal(t, expected, steps)
			})
		}
	}
}

func TestPart2Examples(t *testing.T) {

}

func TestAnswers(t *testing.T) {
	// assert.Equal(t, 0, part1())
	// assert.Equal(t, 0, part2())
}

// func ExampleMain() {
// 	main()
// 	//Output:
// 	//Part 1: 0
// 	//Part 2: 0
// }

func BenchmarkPart1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		part1()
	}
}

func BenchmarkPart2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		part2()
	}
}
