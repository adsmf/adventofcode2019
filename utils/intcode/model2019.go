package intcode

import (
	"fmt"
	"strconv"
	"strings"
)

// Model2019 sets the behaviour of the intcode machine to AoC 2019 rules
func Model2019(input <-chan int, output chan<- int) MachineOption {
	return func(m *Machine) { m.model = &model2019{machine: m} }
}

type model2019 struct {
	machine *Machine

	input  <-chan int
	output chan<- int
}

func (m *model2019) name() string {
	return "Model2019"
}

func (m *model2019) parse(program string) error {
	programIntStrings := strings.Split(strings.TrimSpace(program), ",")

	for pos, valString := range programIntStrings {
		value, err := strconv.Atoi(valString)
		if err != nil {
			panic(err)
		}
		decode := &baseInteger{
			machine: m.machine,
			address: address(pos),
			value:   value,
		}
		m.machine.ram[address(pos)] = decode
	}
	for addr := address(0); int(addr) < len(programIntStrings); addr++ {
		op := model2019operation{
			baseInteger: m.machine.ram[addr].(*baseInteger),
		}
		opCode := op.Value() % 100
		// opMode := op.Value() / 100
		// fmt.Printf("Decoding operation %d: %d / %d\n", addr, opCode, opMode)
		switch opCode {
		case 1:
			op.repr = "ADD"
			op.numParams = 3
			op.handler = m.add
		case 2:
			op.repr = "MUL"
			op.numParams = 3
			op.handler = m.multiply
		case 3:
			op.repr = "SAV"
			op.numParams = 1
			op.handler = m.saveInput
		case 4:
			op.repr = "OUT"
			op.numParams = 1
			op.handler = m.writeOutput
		case 5:
		case 6:
		case 7:
		case 8:
		case 99:
			op.repr = "HCF"
			op.numParams = 0
			op.handler = m.hcf
		default:
			op.repr = fmt.Sprintf("UNK-%d", opCode)
			// return fmt.Errorf("Op code %d, not implemented", op.Value())
		}
		m.machine.ram[addr] = op
		m.machine.operations[addr] = op
		addr += address(op.numParams)
	}
	// return fmt.Errorf("Not implemented")
	return nil
}

func (m *model2019) add(op *model2019operation) bool {
	a := op.getValue(1, 0)
	b := op.getValue(2, 0)
	c := op.getParamAddress(3, 0)

	m.machine.write(c, a.Value()+b.Value())
	m.machine.instructionPointer += 4
	return false
}

func (m *model2019) multiply(op *model2019operation) bool {
	a := op.getValue(1, 0)
	b := op.getValue(2, 0)
	c := op.getParamAddress(3, 0)

	m.machine.write(c, a.Value()*b.Value())
	m.machine.instructionPointer += 4
	return false
}

func (m *model2019) saveInput(op *model2019operation) bool {

	loc := op.getParamAddress(1, 0)

	m.machine.write(loc, <-m.input)
	m.machine.instructionPointer += 2
	return false
}

func (m *model2019) writeOutput(op *model2019operation) bool {
	val := op.getValue(1, 0)
	m.output <- val.Value()
	m.machine.instructionPointer += 2
	return false
}

func (m *model2019) hcf(op *model2019operation) bool {
	return true
}

type model2019opHandler func(*model2019operation) bool
type model2019operation struct {
	baseInteger *baseInteger

	repr       string
	numParams  int
	paramModes map[int]paramMode

	handler model2019opHandler
}

func (mo model2019operation) Address() address         { return mo.baseInteger.Address() }
func (mo model2019operation) IntegerType() integerType { return mo.baseInteger.IntegerType() }
func (mo model2019operation) Value() int               { return mo.baseInteger.Value() }
func (mo model2019operation) SetValue(newVal int)      { mo.baseInteger.SetValue(newVal) }

func (mo model2019operation) Name() string   { return mo.repr }
func (mo model2019operation) NumParams() int { return mo.numParams }

func (mo model2019operation) Exec() bool {
	if mo.handler == nil {
		panic(fmt.Errorf("Handler not defined for %#v", mo))
	}
	halt := mo.handler(&mo)
	fmt.Printf("IP: %d\n", mo.baseInteger.machine.instructionPointer)
	return halt
}

func (mo model2019operation) getValue(offset address, mode paramMode) integer {
	return mo.baseInteger.machine.readAddress(mo.getParamAddress(offset, mode))
}

func (mo model2019operation) getParamAddress(offset address, mode paramMode) address {
	// response := address(0)
	baseAddress := mo.baseInteger.machine.instructionPointer
	switch mode {
	case paramModePosition:
		addr := baseAddress + offset
		dereferenced := mo.baseInteger.machine.readAddress(addr).Value()
		return address(dereferenced)
	case paramModeImmediate:
		return baseAddress
	case paramModeRelative:
		addr := baseAddress + offset
		dereferenced := mo.baseInteger.machine.readAddress(addr).Value()
		return address(dereferenced)
	default:
		panic(fmt.Sprintf("Unknown parameter mode %v", mode))
	}
}

func (mo model2019operation) String() string { return strconv.Itoa(mo.baseInteger.value) }
func (mo model2019operation) GoString() string {
	retString := fmt.Sprintf("%s", mo.repr)
	for i := 0; i < mo.numParams; i++ {
		paramAddress := mo.baseInteger.address + address(i+1)
		paramInteger := mo.baseInteger.machine.readAddress(paramAddress)
		switch mo.paramModes[i] {
		case paramModePosition:
			addr := address(paramInteger.Value())
			dereferenced := mo.baseInteger.machine.readAddress(addr).Value()
			retString = fmt.Sprintf("%s #%#v (%d)", retString, paramInteger, dereferenced)
		case paramModeImmediate:
			retString = fmt.Sprintf("%s '%v'", retString, paramInteger)
		case paramModeRelative:
			addr := address(paramInteger.Value())
			dereferenced := mo.baseInteger.machine.readAddress(addr).Value()
			retString = fmt.Sprintf("%s ~#%#v (%d)", retString, paramInteger, dereferenced)
		}
	}
	return retString
}

type paramMode int

const (
	paramModePosition paramMode = iota
	paramModeImmediate
	paramModeRelative
)
