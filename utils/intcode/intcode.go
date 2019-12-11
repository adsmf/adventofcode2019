package intcode

import (
	"fmt"
	"sort"
	"strings"
)

// NewMachine creates a new intcode machine
func NewMachine(options ...MachineOption) *Machine {
	m := &Machine{
		ram:        map[address]integer{},
		operations: map[address]operation{},
	}
	for _, option := range options {
		option(m)
	}
	return m
}

// Machine is a virtual machine capable of running intcode (e.g. https://adventofcode.com/2019/day/2)
type Machine struct {
	model              model
	ram                ram
	operations         operationMap
	instructionPointer address
}

// LoadProgram wipes the machine and loads a new program from an input string
func (m *Machine) LoadProgram(program string) error {
	if m.model == nil {
		return fmt.Errorf("Cannot parse program: No intcode machine model defined")
	}
	return m.model.parse(program)
}

func (m *Machine) Run() {
	for {
		halt := m.Step()
		if halt {
			return
		}
	}
}

func (m *Machine) Step() bool {
	operation, isOp := m.operations[m.instructionPointer]
	if !isOp {
		panic(fmt.Errorf("Current instruction pointer %v is not an operation", m.instructionPointer))
	}
	return operation.Exec()
}

func (m Machine) GoString() string {
	state := []string{
		"Model: " + m.model.name(),
		"Instruction Pointer: " + m.instructionPointer.String(),
		fmt.Sprintf("Program:\n%#v", m.operations),
	}
	stateString := ""
	for _, line := range state {
		stateString += line + "\n"
	}

	return strings.TrimSpace(stateString)
}

func (m *Machine) readAddress(addr address) integer {
	return m.ram[addr]
}

func (m *Machine) write(addr address, newVal int) {
	m.ram[addr].SetValue(newVal)
}

type ram map[address]integer

func (r ram) String() string {
	// state := []string{}
	addresses := []int{}
	for addr := range r {
		addresses = append(addresses, int(addr))
	}
	sort.Ints(addresses)
	stateString := ""
	for _, addr := range addresses {
		if stateString != "" {
			stateString += ","
		}
		stateString += fmt.Sprintf("%v", r[address(addr)])
	}
	return stateString
}

type operationMap map[address]operation

func (om operationMap) GoString() string {
	state := []string{
		fmt.Sprintf("\tNum operations: %d", len(om)),
	}

	operationAddresses := []int{}
	for addr := range om {
		operationAddresses = append(operationAddresses, int(addr))
	}
	sort.Ints(operationAddresses)

	for _, addr := range operationAddresses {
		state = append(state, fmt.Sprintf("\t%#v: %#v", address(addr), om[address(addr)]))
	}

	stateString := ""
	for _, line := range state {
		stateString += line + "\n"
	}
	return stateString

}

type model interface {
	name() string
	parse(program string) error
}

// MachineOption defines configuration options that can be applied to an intcode machine
type MachineOption func(*Machine)

type address int

func (a address) String() string {
	return fmt.Sprintf("#%04d", a)
}
func (a address) GoString() string {
	return fmt.Sprintf("#%04d", a)
}

type integer interface {
	Address() address
	IntegerType() integerType
	Value() int
	SetValue(int)
}

type baseInteger struct {
	machine     *Machine
	address     address
	integerType integerType
	value       int
}

func (i baseInteger) Address() address {
	return i.address
}

func (i baseInteger) IntegerType() integerType {
	return i.integerType
}

func (i baseInteger) Value() int {
	return i.value
}
func (i *baseInteger) SetValue(newVal int) {
	i.value = newVal
}

func (i baseInteger) String() string {
	return fmt.Sprintf("%d", i.value)
}
func (i baseInteger) GoString() string {
	return fmt.Sprintf("%d", i.value)
}

type integerType int

const (
	integerTypeInstruction integerType = iota
	integerTypeData
)

type operation interface {
	Exec() (halt bool)
	Name() string
	NumParams() int
}
