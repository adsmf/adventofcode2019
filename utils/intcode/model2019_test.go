package intcode

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModel2019(t *testing.T) {
	type testDef struct {
		program  string
		endState string
		inputs   []int
		outputs  []int
	}
	tests := []testDef{
		testDef{
			program:  "1,0,0,1,99",
			endState: "1,2,0,1,99",
		},
		testDef{
			program:  "2,3,0,3,99",
			endState: "2,3,0,6,99",
		},
	}
	for id, test := range tests {
		t.Run(fmt.Sprintf("Test %d", id), func(t *testing.T) {
			t.Logf("Test definition:\n%#v", test)
			inputStream := make(chan int)
			outputStream := make(chan int)
			m := NewMachine(Model2019(inputStream, outputStream))
			m.LoadProgram(test.program)
			t.Logf("Initial machine state:\n%#v", m)
			// t.Logf("Initial machine RAM:\n%v", m.ram)
			m.Run()

			t.Logf("Stepped machine state:\n%#v", m)
			if test.endState != "" {
				assert.Equal(t, test.endState, m.ram.String())
			}
		})
	}
}
