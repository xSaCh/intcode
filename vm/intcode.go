package vm

import (
	"fmt"
	"log"
)

const (
	MODE_POSITION = iota
	MODE_IMMEDIATE
	MODE_RELATIVE
)

type IntcodeVM struct {
	PcRegister       int
	RelativeRegister int
	Memory           []int
	Mode             int
	Halt             bool

	InputFunc  func() int
	OutputFunc func(int)
}

func CreateVM() IntcodeVM {
	return IntcodeVM{PcRegister: 0, RelativeRegister: 0, Memory: []int{}, Mode: MODE_POSITION, Halt: false}
}

func (vm *IntcodeVM) LoadProgram(program []int) {
	vm.Memory = program
	vm.PcRegister = 0
}

// Return when vm is on halt with value at address '0' of memory
func (vm *IntcodeVM) Run() (int, error) {
	for vm.PcRegister < len(vm.Memory) {
		opcode := vm.ReadNextOpcode()
		switch opcode {
		case OPCODE_ADD:
			aI := vm.ReadNext(vm.Mode)
			bI := vm.ReadNext(vm.Mode)
			ans := aI + bI

			vm.WriteNext(ans, vm.Mode)
			vm.Next()

		case OPCODE_MUL:
			aI := vm.ReadNext(vm.Mode)
			bI := vm.ReadNext(vm.Mode)
			ans := aI * bI

			vm.WriteNext(ans, vm.Mode)
			vm.Next()
		case OPCODE_INPUT:
			inp := vm.ReadInput()
			vm.WriteNext(inp, vm.Mode)
			vm.Next()
		case OPCODE_OUTPUT:
			aI := vm.ReadNext(vm.Mode)
			vm.WriteOutput(aI)
			vm.Next()

		case OPCODE_JMP_T:
			aI := vm.ReadNext(vm.Mode)
			bI := vm.ReadNext(MODE_IMMEDIATE)
			if aI != 0 {
				vm.Jump(bI)
			} else {
				vm.Next()
			}
		case OPCODE_JMP_F:
			aI := vm.ReadNext(vm.Mode)
			bI := vm.ReadNext(MODE_IMMEDIATE)
			if aI == 0 {
				vm.Jump(bI)
			} else {
				vm.Next()
			}
		case OPCODE_LESS_THAN:
			aI := vm.ReadNext(vm.Mode)
			bI := vm.ReadNext(vm.Mode)
			if aI < bI {
				vm.WriteNext(1, vm.Mode)
			} else {
				vm.WriteNext(0, vm.Mode)
			}
			vm.Next()
		case OPCODE_EQUALS:
			aI := vm.ReadNext(vm.Mode)
			bI := vm.ReadNext(vm.Mode)
			if aI == bI {
				vm.WriteNext(1, vm.Mode)
			} else {
				vm.WriteNext(0, vm.Mode)
			}
			vm.Next()
		case OPCODE_INC_RELV:
			aI := vm.ReadNext(vm.Mode)
			vm.WriteRelvRegister(aI + vm.ReadRelvRegister())
			vm.Next()
		case OPCODE_HALT:
			vm.Next()
			return vm.Memory[0], nil
		}
	}

	return -1, nil
}

func (vm *IntcodeVM) Jump(address int) {
	vm.PcRegister = address
}
func (vm *IntcodeVM) Next() {
	vm.PcRegister++
}

func (vm *IntcodeVM) ReadNextOpcode() int {
	val := vm.Memory[vm.PcRegister]
	return val

}
func (vm *IntcodeVM) ReadNext(mode int) int {
	switch mode {
	case MODE_POSITION:
		vm.Next()
		return vm.Memory[vm.Memory[vm.PcRegister]]
	case MODE_IMMEDIATE:
		vm.Next()

		return vm.Memory[vm.PcRegister]
	case MODE_RELATIVE:
		// TODO: To be implemented
		vm.Next()
		return vm.Memory[vm.Memory[vm.PcRegister]]
	default:
		log.Fatalf("Unknown mode: %d\n", mode)
		return -1
	}
}

func (vm *IntcodeVM) WriteNext(value int, mode int) {
	switch mode {
	case MODE_POSITION:
		vm.Next()
		vm.Memory[vm.Memory[vm.PcRegister]] = value
	case MODE_RELATIVE:
		// TODO: To be implemented
		vm.Memory[vm.ReadNext(mode)] = value
	default:
		log.Fatalf("Unknown mode: %d\n", mode)
	}
}

func (vm *IntcodeVM) ReadInput() int {
	if vm.InputFunc != nil {
		return vm.InputFunc()
	}
	v := 0
	fmt.Scanf("%d", &v)
	return v
}
func (vm *IntcodeVM) WriteOutput(val int) {
	if vm.OutputFunc != nil {
		vm.OutputFunc(val)
		return
	}
	print(val)

}

func (vm *IntcodeVM) ReadRelvRegister() int {
	return vm.RelativeRegister
}
func (vm *IntcodeVM) WriteRelvRegister(value int) {
	vm.RelativeRegister = value
}
func (vm *IntcodeVM) DumpMemory() []int {
	return vm.Memory
}
