package main

import (
	"github.com/xSaCh/intcode/vm"
)

func main() {

	vm := vm.CreateVM()
	vm.InputFunc = nil
	vm.OutputFunc = nil
	vm.LoadProgram([]int{3, 0, 4, 0, 99})
	vm.Run()
	vm.DumpMemory()
	// v := 0
	// fmt.Scanf("%d", &v)
	// fmt.Printf("A: %d\n", v)
}
