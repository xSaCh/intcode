package vm_test

import (
	"fmt"
	"os"
	"slices"
	"testing"

	"github.com/xSaCh/intcode/vm"
)

func loadNRun(memory []int) (*vm.IntcodeVM, error) {
	vm := vm.CreateVM()
	if vm.PcRegister != 0 {
		return nil, fmt.Errorf("Expected vm.Pc to be 0, got %v", vm.PcRegister)
	}

	if len(vm.Memory) != 0 {
		return nil, fmt.Errorf("Expected vm.Memory to be empty, got %v", vm.Memory)
	}

	vm.LoadProgram(memory)

	_, err := vm.Run()
	if err != nil {
		return nil, fmt.Errorf("Error: %v", err)
	}
	return &vm, nil
}

func expectedDump(memory []int, expectedDump []int) (bool, error) {
	vm, err := loadNRun(memory)
	if err != nil {
		return false, err
	}
	dump := vm.DumpMemory()
	if !slices.Equal(dump, expectedDump) {
		return false, fmt.Errorf("Error: Expected %v, got %v\n", expectedDump, dump)
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
	_, w, err := os.Pipe()
	defer func(v *os.File) { os.Stdin = v }(os.Stdin)
	w.Write([]byte{'5'})
	w.Close()

	pass, err := expectedDump([]int{3, 0, 4, 0, 99}, []int{5, 0, 4, 0, 99})
	if !pass {
		t.Fatal(err)
	}
}
func TestJmpIfTrue(t *testing.T) {
	pass, err := expectedDump([]int{5, 1, 4, 99, 1, 0, 0, 0, 99}, []int{10, 1, 4, 99, 1, 0, 0, 0, 99})
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

func TestC1(t *testing.T) {

	pass, err := expectedDump([]int{1, 9, 10, 3, 2, 3, 11, 0, 99, 30, 40, 50}, []int{3500, 9, 10, 70, 2, 3, 11, 0, 99, 30, 40, 50})
	if !pass {
		t.Fatal(err)
	}
}
