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
		opcode, mod := ParseOpcode(memory[i])

		switch opcode {
		case vm.OPCODE_ADD:

			curLine = []string{"ADD", FormattedMem(memory, i+1, mod[0]),
				FormattedMem(memory, i+2, mod[1]), FormattedMem(memory, i+3, mod[2])}
			i += 3
		case vm.OPCODE_MUL:

			curLine = []string{"MUL", FormattedMem(memory, i+1, mod[0]),
				FormattedMem(memory, i+2, mod[1]), FormattedMem(memory, i+3, mod[2])}
			i += 3
		case vm.OPCODE_INPUT:

			curLine = []string{"IN", FormattedMem(memory, i+1, mod[0])}
			i++
		case vm.OPCODE_OUTPUT:

			curLine = []string{"OUT", FormattedMem(memory, i+1, mod[0])}
			i++
		case vm.OPCODE_JMP_T:

			curLine = []string{"JZ", FormattedMem(memory, i+1, mod[0]),
				FormattedMem(memory, i+2, mod[1])}
			i += 2
		case vm.OPCODE_JMP_F:

			curLine = []string{"JNZ", FormattedMem(memory, i+1, mod[0]),
				FormattedMem(memory, i+2, mod[1])}
			i += 2
		case vm.OPCODE_LESS_THAN:

			curLine = []string{"SLT", FormattedMem(memory, i+1, mod[0]),
				FormattedMem(memory, i+2, mod[1]), FormattedMem(memory, i+3, mod[2])}
			i += 3
		case vm.OPCODE_EQUALS:

			curLine = []string{"SEQ", FormattedMem(memory, i+1, mod[0]),
				FormattedMem(memory, i+2, mod[1]), FormattedMem(memory, i+3, mod[2])}
			i += 3
		case vm.OPCODE_INC_RELV:

			curLine = []string{"INCB", FormattedMem(memory, i+1, mod[0])}
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

func FormattedMem(memory []int, ind, mod int) string {
	switch mod {
	case 0:
		return fmt.Sprintf("#%d", memory[ind])
	case 1:
		return fmt.Sprintf("%d", memory[ind])
	case 2:
		// Handle relative case
		return fmt.Sprintf("%d", memory[ind])
	default:
		return ""
	}
}
