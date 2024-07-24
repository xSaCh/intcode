package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/xSaCh/intcode/vm"
)

func main() {
	d22()
}

func d21() {
	data, err := os.ReadFile("./inputs/day2_1.txt")

	if err != nil {
		log.Fatalf("Error %v\n", err)
		return

	}

	val := strings.Split(string(data), ",")
	instructions := make([]int, len(val))
	for i, v := range val {
		instructions[i], _ = strconv.Atoi(v)
	}
	// fmt.Printf("d: %v\n", instructions)
	instructions[1] = 12
	instructions[2] = 2
	vm := vm.CreateVM()
	vm.LoadProgram(instructions)
	ans, err := vm.Run()
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	fmt.Println(ans)

}
func d22() {
	data, err := os.ReadFile("./inputs/day2_1.txt")

	if err != nil {
		log.Fatalf("Error %v\n", err)
		return

	}

	val := strings.Split(string(data), ",")
	instructions := make([]int, len(val))
	for i, v := range val {
		instructions[i], _ = strconv.Atoi(v)
	}
	vm := vm.CreateVM()
	for i := 0; i < 100; i++ {
		for j := 0; j < 100; j++ {
			clone := make([]int, len(instructions))
			copy(clone, instructions)

			clone[1] = i
			clone[2] = j
			vm.LoadProgram(clone)
			ans, err := vm.Run()
			if err != nil {
				log.Fatalf("Error: %v\n", err)
			}

			// fmt.Printf("%d %d\n", ans, (100*i)+j)
			if ans == 19690720 {
				fmt.Println((100 * i) + j)
			}
		}
	}
	// fmt.Printf("d: %v\n", instructions)

}
