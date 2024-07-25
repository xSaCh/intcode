package vm

import (
	"fmt"
	"log"
	"strconv"
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
	Mode             []int
	Halt             bool

	InputFunc  func() int
	OutputFunc func(int)
}

func CreateVM() IntcodeVM {
	return IntcodeVM{PcRegister: 0, RelativeRegister: 0, Memory: []int{}, Mode: []int{}, Halt: false}
}

func (vm *IntcodeVM) LoadProgram(program []int) {
	vm.Memory = program
	vm.PcRegister = 0
}

// Return when vm is on halt with value at address '0' of memory
func (vm *IntcodeVM) Run() (int, error) {
	for vm.PcRegister < len(vm.Memory) {
		opcode := 0
		opcode, vm.Mode = vm.ReadNextOpcode()
		switch opcode {
		case OPCODE_ADD:
			aI := vm.ReadNext(vm.Mode[0])
			bI := vm.ReadNext(vm.Mode[1])
			ans := aI + bI

			vm.WriteNext(ans, vm.Mode[2])
			vm.Next()

		case OPCODE_MUL:
			aI := vm.ReadNext(vm.Mode[0])
			bI := vm.ReadNext(vm.Mode[1])
			ans := aI * bI

			vm.WriteNext(ans, vm.Mode[2])
			vm.Next()
		case OPCODE_INPUT:
			inp := vm.ReadInput()
			vm.WriteNext(inp, vm.Mode[0])
			vm.Next()
		case OPCODE_OUTPUT:
			aI := vm.ReadNext(vm.Mode[0])
			vm.WriteOutput(aI)
			vm.Next()

		case OPCODE_JMP_T:
			aI := vm.ReadNext(vm.Mode[0])
			bI := vm.ReadNext(vm.Mode[1])
			if aI != 0 {
				vm.Jump(bI)
			} else {
				vm.Next()
			}
		case OPCODE_JMP_F:
			aI := vm.ReadNext(vm.Mode[0])
			bI := vm.ReadNext(MODE_IMMEDIATE)
			if aI == 0 {
				vm.Jump(bI)
			} else {
				vm.Next()
			}
		case OPCODE_LESS_THAN:
			aI := vm.ReadNext(vm.Mode[0])
			bI := vm.ReadNext(vm.Mode[1])
			if aI < bI {
				vm.WriteNext(1, vm.Mode[2])
			} else {
				vm.WriteNext(0, vm.Mode[2])
			}
			vm.Next()
		case OPCODE_EQUALS:
			aI := vm.ReadNext(vm.Mode[0])
			bI := vm.ReadNext(vm.Mode[1])
			if aI == bI {
				vm.WriteNext(1, vm.Mode[2])
			} else {
				vm.WriteNext(0, vm.Mode[2])
			}
			vm.Next()
		case OPCODE_INC_RELV:
			aI := vm.ReadNext(vm.Mode[0])
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

func (vm *IntcodeVM) ReadNextOpcode() (int, []int) {
	val := vm.Memory[vm.PcRegister]
	n := len(strconv.Itoa(val))
	mode := []int{0, 0, 0}
	opc := val

	// Explicit speficed Mode
	if n > 2 {
		opc = val % 100
		val /= 100

		mode[0] = val % 10
		val /= 10
		mode[1] = val % 10
		val /= 10

		// For handling padding at starting
		mode[2] = val % 10
		return opc, mode

	}
	// fmt.Printf("VAL: %d OP: %d mode: %v\n", val, opc, mode)
	// Set default Modes
	// switch opc {
	// case OPCODE_JMP_T:
	// 	mode[1] = MODE_IMMEDIATE
	// case OPCODE_JMP_F:
	// 	mode[1] = MODE_IMMEDIATE
	// }
	return opc, mode

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
		vm.Next()
		return vm.Memory[vm.RelativeRegister+vm.Memory[vm.PcRegister]]
	default:
		log.Fatalf("Unknown mode: %d\n", mode)
		return -1
	}
}

func (vm *IntcodeVM) WriteNext(value int, mode int) {
	switch mode {
	case MODE_POSITION:
		vm.Next()

		// Dynamic Memory Increasing
		if vm.Memory[vm.PcRegister] >= len(vm.Memory) {
			n := vm.Memory[vm.PcRegister] - len(vm.Memory)
			vm.Memory = append(vm.Memory, make([]int, n+1)...)
		}

		vm.Memory[vm.Memory[vm.PcRegister]] = value
	case MODE_RELATIVE:
		if vm.Memory[vm.PcRegister] > len(vm.Memory) {
			n := vm.Memory[vm.PcRegister] - len(vm.Memory)
			vm.Memory = append(vm.Memory, make([]int, n+1)...)
		}
		vm.Memory[vm.RelativeRegister+vm.ReadNext(mode)] = value
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
