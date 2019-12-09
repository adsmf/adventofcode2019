package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// var debug = debugPrintf
var debug = noOut

func main() {
	fmt.Printf("Part 1: %d\n", part1())
	fmt.Printf("Part 2: %d\n", part2())
}

func part1() int64 {
	inputString := loadInputString()
	outputs := gatherOutputs(inputString, -1, 1)
	fmt.Printf("Outputs: %#v\n", outputs)
	return outputs[len(outputs)-1]
}

func part2() []int64 {
	inputString := loadInputString()
	outputs := gatherOutputs(inputString, -1, 2)
	return outputs
}

func gatherOutputs(program string, maxSteps, in int64) []int64 {
	outputs := []int64{}
	output := make(chan int64)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for op := range output {
			outputs = append(outputs, op)
		}
		wg.Done()
	}()

	input := make(chan int64, 1)
	input <- in
	tape := newMachine(program, input, output)
	tape.run(maxSteps)

	wg.Wait()
	debug("Final state: %#v\n", tape)
	return outputs
}

func loadInputString() string {
	inputRaw, err := ioutil.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	return string(inputRaw)

}

type machine struct {
	headPos      int64
	values       map[int64]int64
	inputs       <-chan int64
	outputs      chan<- int64
	relativeBase int64
}

func newMachine(initial string, inputs <-chan int64, output chan<- int64) machine {
	initialValueStrings := strings.Split(strings.TrimSpace(initial), ",")
	initialValues := map[int64]int64{}
	for pos, valString := range initialValueStrings {
		val, err := strconv.Atoi(valString)
		if err != nil {
			panic(err)
		}
		initialValues[int64(pos)] = int64(val)
	}
	mach := machine{
		values:  initialValues,
		headPos: 0,
		inputs:  inputs,
		outputs: output,
	}
	return mach
}

func (t *machine) run(maxSteps int64) {
	if maxSteps == -1 {
		maxSteps = 10000
	}
	steps := maxSteps
	for {
		done := t.step()
		if done {
			return
		}
		steps--
		if steps <= 0 {
			close(t.outputs)
			return
		}
	}
}

func (t *machine) step() bool {
	// debug("#100 => %d", t.values[100])
	initialHead := t.headPos
	oper := t.values[initialHead]
	debug("  -- Addr: %d; op: %d; base: %d\n", initialHead, oper, t.relativeBase)
	paramModes := int64(oper / 100)
	oper = oper % 100
	switch oper {
	case 1:
		// Add
		params := t.getParams(paramModes, 3, true)
		p1 := params[0]
		p2 := params[1]
		p3 := params[2]

		t.values[p3] = p1 + p2
	case 2:
		// Mult
		params := t.getParams(paramModes, 3, true)
		p1 := params[0]
		p2 := params[1]
		p3 := params[2]

		t.values[p3] = p1 * p2
	case 3:
		// Input
		params := t.getParams(paramModes, 1, true)
		p := params[0]

		nextInput := <-t.inputs
		debug("STO %d => #%d", nextInput, p)
		t.values[p] = nextInput
	case 4:
		// Output
		params := t.getParams(paramModes, 1, false)
		p := params[0]
		debug("READ %d", p)
		t.outputs <- p
	case 5:
		// JNZ
		params := t.getParams(paramModes, 2, false)
		p1 := params[0]
		p2 := params[1]

		if p1 != 0 {
			t.headPos = p2
		}
	case 6:
		// JEZ
		params := t.getParams(paramModes, 2, false)
		p1 := params[0]
		p2 := params[1]

		if p1 == 0 {
			t.headPos = p2
		}
	case 7:
		// CLT
		params := t.getParams(paramModes, 3, true)
		p1 := params[0]
		p2 := params[1]
		p3 := params[2]

		if p1 < p2 {
			t.values[p3] = 1
		} else {
			t.values[p3] = 0
		}
	case 8:
		// CMP
		params := t.getParams(paramModes, 3, true)
		p1 := params[0]
		p2 := params[1]
		p3 := params[2]

		// debug("CMP %d == %d => %v", p1, p2, p1 == p2)
		if p1 == p2 {
			t.values[p3] = 1
		} else {
			t.values[p3] = 0
		}
	case 9:
		// Update relative base
		params := t.getParams(paramModes, 1, false)
		p1 := params[0]
		t.relativeBase += p1
		// debug("Adjusted relative base by %d to %d\n", p1, t.relativeBase)
	case 99:
		// HCF
		close(t.outputs)
		return true
	default:
		panic(fmt.Errorf("Invalid opcode %d at position %d: %#v", oper, t.headPos, t))
	}
	return false
}

func (t *machine) getParams(paramModes, numParams int64, lastIsAddress bool) []int64 {
	params := []int64{}
	for param := int64(0); param < numParams; param++ {
		lastParam := (param == numParams-1)
		p := t.getVal(t.headPos + param + 1)
		mode := paramMode(paramModes, param)
		returnAddress := lastParam && lastIsAddress
		switch mode {
		case 0:
			if !returnAddress {
				p = t.getVal(p)
			}
		case 1:
			// Immediate mode, nothing further to do
		case 2:
			p += t.relativeBase
			if !returnAddress {
				p = t.getVal(p)
			} else {
			}
		default:
			panic(fmt.Errorf("Unknown parameter mode %d", mode))
		}
		params = append(params, p)
	}

	t.headPos = t.headPos + numParams + 1
	return params
}

func paramMode(modes, pos int64) int64 {
	mask := int64(math.Pow(10, float64(pos)))
	return (modes / mask) % 10
}

func (t *machine) getVal(pos int64) int64 {
	if pos < 0 {
		panic("Cannot read negative address")
	}
	return t.values[pos]
}

func (t *machine) String() string {
	valueStrings := []string{}
	keys := []int{}
	for key := range t.values {
		keys = append(keys, int(key))
	}
	sort.Ints(keys)
	for _, key := range keys {
		valueStrings = append(valueStrings, strconv.Itoa(int(t.values[int64(key)])))
	}
	return strings.Join(valueStrings, ",")
}

func debugPrintf(format string, params ...interface{}) {
	fmt.Printf(format, params...)
}
func noOut(format string, params ...interface{}) {}