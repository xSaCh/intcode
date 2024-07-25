package main

import (
	"fmt"

	"github.com/xSaCh/intcode/vm"
)

func main() {

	vm := vm.CreateVM()
	vm.InputFunc = func() int { return 5 }
	a := []int{}
	vm.OutputFunc = func(i int) {
		a = append(a, i)
	}
	// vm.LoadProgram([]int{3, 9, 3, 10, 4, 9, 4, 10, 99, 0, 0})
	vm.LoadProgram([]int{3, 34, 1007, 34, 1, 35, 1005, 35, 30, 1001, 34, 0, 33, 1001, 33, -1, 33, 1006, 33, 27, 2, 34, 33, 34, 1005, 33, 13, 4, 34, 99, 104, 1, 99})
	vm.Run()
	fmt.Printf("a: %v\n", a)
	fmt.Printf("vm.DumpMemory(): %v\n", vm.DumpMemory())
	// v := 0
	// fmt.Scanf("%d", &v)
	// fmt.Printf("A: %d\n", v)
}
