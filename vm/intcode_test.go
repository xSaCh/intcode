package vm_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/xSaCh/intcode/vm"
)

func load(memory []int) (*vm.IntcodeVM, error) {
	vm := vm.CreateVM()

	if vm.PcRegister != 0 {
		return nil, fmt.Errorf("Expected vm.Pc to be 0, got %v", vm.PcRegister)
	}

	if len(vm.Memory) != 0 {
		return nil, fmt.Errorf("Expected vm.Memory to be empty, got %v", vm.Memory)
	}

	vm.LoadProgram(memory)

	return &vm, nil
}

func expectedDump(memory, expectedDump []int) (bool, error) {
	vm, err := load(memory)
	if err != nil {
		return false, err
	}
	_, err = vm.Run()
	if err != nil {
		return false, fmt.Errorf("Error: %v", err)
	}

	dump := vm.DumpMemory()
	if !slices.Equal(dump, expectedDump) {
		return false, fmt.Errorf("Error: Expected %v, got %v\n", expectedDump, dump)
	}
	return true, nil
}

func expectedOutput(memory, inputs, expectedOutput []int) (bool, error) {
	vm, err := load(memory)
	if err != nil {
		return false, err
	}
	ii := 0
	op := []int{}
	vm.InputFunc = func() int {
		ii += 1
		return inputs[ii-1]
	}

	vm.OutputFunc = func(i int) {
		op = append(op, i)
	}
	_, err = vm.Run()
	if err != nil {
		return false, fmt.Errorf("Error: %v", err)
	}

	// dump := vm.DumpMemory()
	if !slices.Equal(op, expectedOutput) {
		return false, fmt.Errorf("Error: Expected %v, got %v\n", expectedOutput, op)
	}
	return true, nil
}

func TestAdd(t *testing.T) {
	pass, err := expectedDump([]int{1, 0, 0, 0, 99}, []int{2, 0, 0, 0, 99})
	if !pass {
		t.Fatal(err)
	}

}
func TestMul(t *testing.T) {
	pass, err := expectedDump([]int{2, 3, 0, 3, 99}, []int{2, 3, 0, 6, 99})
	if !pass {
		t.Fatal(err)
	}
}
func TestIO(t *testing.T) {
	pass, err := expectedOutput([]int{3, 0, 4, 0, 99}, []int{6}, []int{6})
	if !pass {
		t.Fatal(err)
	}
}
func TestJmpIfTrue(t *testing.T) {
	pass, err := expectedDump([]int{1005, 1, 4, 99, 1, 0, 0, 0, 99}, []int{2010, 1, 4, 99, 1, 0, 0, 0, 99})
	if !pass {
		t.Fatal(err)
	}
}

func TestJmpIfFalse(t *testing.T) {
	pass, err := expectedDump([]int{6, 7, 4, 99, 1, 7, 0, 0, 99}, []int{6, 7, 4, 99, 1, 7, 0, 0, 99})
	if !pass {
		t.Fatal(err)
	}
}
func TestLessThan(t *testing.T) {
	pass, err := expectedDump([]int{7, 3, 4, 5, 99, -1}, []int{7, 3, 4, 5, 99, 1})
	if !pass {
		t.Fatal(err)
	}
}
func TestEqual(t *testing.T) {
	pass, err := expectedDump([]int{8, 3, 4, 5, 99, 0}, []int{8, 3, 4, 5, 99, 0})
	if !pass {
		t.Fatal(err)
	}
}

func TestImmediate(t *testing.T) {
	pass, err := expectedDump([]int{1002, 4, 3, 4, 33}, []int{1002, 4, 3, 4, 99})
	if !pass {
		t.Fatal(err)
	}
}

func TestC1(t *testing.T) {

	pass, err := expectedDump([]int{1, 9, 10, 3, 2, 3, 11, 0, 99, 30, 40, 50}, []int{3500, 9, 10, 70, 2, 3, 11, 0, 99, 30, 40, 50})
	if !pass {
		t.Fatal(err)
	}
}
func TestC2(t *testing.T) {

	pass, err := expectedOutput([]int{3, 9, 3, 10, 4, 9, 4, 10, 99, 0, 0},
		[]int{7, 42}, []int{7, 42})
	if !pass {
		t.Fatal(err)
	}
}
