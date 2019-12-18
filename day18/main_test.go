package main

import (
	"fmt"
	"testing"

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
			file:  "p1ex1.txt",
			steps: 86,
		},
	}
	for id, test := range tests {
		t.Run(fmt.Sprintf("Part1Test%d", id), func(t *testing.T) {
			crawl := loadMap("examples/" + test.file)
			t.Logf("Start:\n%v", crawl)
			steps := crawl.collectKeys()
			assert.Equal(t, test.steps, steps)
		})
	}
}

func TestPart2Examples(t *testing.T) {

}

func TestAnswers(t *testing.T) {
	assert.Equal(t, 0, part1())
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
