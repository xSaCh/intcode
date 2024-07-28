package debugger

import (
	"fmt"
	"strconv"

	"github.com/xSaCh/intcode/vm"
)

func GetFormattedMemory(memory []int) ([][]string, error) {
	lines := [][]string{}
	for i := 0; i < len(memory); i++ {
		var curLine []string
		opcode, _ := ParseOpcode(memory[i])

		switch opcode {
		case vm.OPCODE_ADD:
			a := memory[i+1]
			b := memory[i+2]
			c := memory[i+3]
			curLine = []string{"ADD", fmt.Sprintf("#%d", a), fmt.Sprintf("#%d", b), fmt.Sprintf("#%d", c)}
			i += 3
		case vm.OPCODE_MUL:
			a := memory[i+1]
			b := memory[i+2]
			c := memory[i+3]
			curLine = []string{"MUL", fmt.Sprintf("#%d", a), fmt.Sprintf("#%d", b), fmt.Sprintf("#%d", c)}
			i += 3
		case vm.OPCODE_INPUT:
			a := memory[i+1]
			curLine = []string{"IN", fmt.Sprintf("#%d", a)}
			i++
		case vm.OPCODE_OUTPUT:
			a := memory[i+1]
			curLine = []string{"OUT", fmt.Sprintf("#%d", a)}
			i++
		case vm.OPCODE_JMP_T:
			a := memory[i+1]
			b := memory[i+2]
			curLine = []string{"JZ", fmt.Sprintf("#%d", a), fmt.Sprintf("#%d", b)}
			i += 2
		case vm.OPCODE_JMP_F:
			a := memory[i+1]
			b := memory[i+2]
			curLine = []string{"JNZ", fmt.Sprintf("#%d", a), fmt.Sprintf("#%d", b)}
			i += 2
		case vm.OPCODE_LESS_THAN:
			a := memory[i+1]
			b := memory[i+2]
			c := memory[i+3]
			curLine = []string{"SLT", fmt.Sprintf("#%d", a), fmt.Sprintf("#%d", b), fmt.Sprintf("#%d", c)}
			i += 3
		case vm.OPCODE_EQUALS:
			a := memory[i+1]
			b := memory[i+2]
			c := memory[i+3]
			curLine = []string{"SEQ", fmt.Sprintf("#%d", a), fmt.Sprintf("#%d", b), fmt.Sprintf("#%d", c)}
			i += 3
		case vm.OPCODE_INC_RELV:
			a := memory[i+1]
			curLine = []string{"INCB", fmt.Sprintf("#%d", a)}
			i++
		case vm.OPCODE_HALT:
			curLine = []string{"HALT"}
		default:
			curLine = []string{fmt.Sprintf("%d", opcode)}
			// return lines, fmt.Errorf("Unknown opcode: %d at index: %d", opcode, i)
		}
		lines = append(lines, curLine)
	}
	return lines, nil
}

func ParseOpcode(val int) (int, []int) {
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
	return opc, mode
}
