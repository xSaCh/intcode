package main

import (
	"fmt"

	"github.com/xSaCh/intcode/assembler"
	"github.com/xSaCh/intcode/debugger"
	"github.com/xSaCh/intcode/vm"
)

func mainVM() {

	vm := vm.CreateVM()
	i := 0
	vm.InputFunc = func() int { i++; return 5 }
	// a := []int{}
	// vm.LoadProgram([]int{3, 9, 3, 10, 4, 9, 4, 10, 99, 0, 0})
	// vm.LoadProgram([]int{3, 34, 1007, 34, 1, 35, 1005, 35, 30, 1001, 34, 0, 33, 1001, 33, -1, 33, 1006, 33, 27, 2, 34, 33, 34, 1005, 33, 13, 4, 34, 99, 104, 1, 99})
	vm.LoadProgram([]int{4, 3, 101, 72, 14, 3, 101, 1, 4, 4, 5, 3, 16, 99, 29, 7, 0, 3, -67, -12, 87, -8, 3, -6, -8, -67, -23, -10})
	// vm.Run()
	// fmt.Printf("a: %v\n", a)
	// fmt.Printf("vm.DumpMemory(): %v\n", vm.DumpMemory())
	a, _, err := debugger.GetFormattedMemory(vm.Memory, vm.PcRegister)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	// vm.LoadProgram([]int{3, 9, 3, 10, 4, 9, 4, 10, 99, 0, 0})
	// vm.LoadProgram([]int{3, 34, 1007, 34, 1, 35, 1005, 35, 30, 1001, 34, 0, 33, 1001, 33, -1, 33, 1006, 33, 27, 2, 34, 33, 34, 1005, 33, 13, 4, 34, 99, 104, 1, 99})
	// vm.Run()
	fmt.Printf("a: %v\n", a)
	debugger.Run(&vm)

}

func mainAs() {
	a := assembler.NewAssembler()
	a.AssembleFromFile("test.asm")
	fmt.Printf("a.Tokens: %v\n", a.Tokens)
	fmt.Printf("a.ParsedOpcodes: %v\n", a.ParsedOpcodes)
	fmt.Printf("a.ByteCode: %v\n", a.ByteCode)
	// fmt.Printf("assembler.Start(\"test.asm\"): %v\n", a.ByteCode)
}

func main() {
	mainAs()
	// mainVM()
}
