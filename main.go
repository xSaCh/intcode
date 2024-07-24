package main

import (
	"fmt"

	"github.com/xSaCh/intcode/vm"
)

func main() {

	vm := vm.CreateVM()
	vm.InputFunc = nil
	vm.OutputFunc = nil
	vm.LoadProgram([]int{3, 9, 3, 10, 4, 9, 4, 10, 99, 0, 0})
	vm.Run()
	fmt.Printf("vm.DumpMemory(): %v\n", vm.DumpMemory())
	// v := 0
	// fmt.Scanf("%d", &v)
	// fmt.Printf("A: %d\n", v)
}
