package debugger

import (
	"fmt"
	"strconv"

	"github.com/xSaCh/intcode/vm"
)

func GetFormattedMemory(memory []int, pc int) ([][]string, int, error) {
	lines := [][]string{}
	activeInd := -1
	for i := 0; i < len(memory); i++ {
		var curLine []string
		opcode, mod := ParseOpcode(memory[i])

		if i == pc {
			activeInd = len(lines)
		}
		switch opcode {
		case vm.OPCODE_ADD:
			if i+3 >= len(memory) {
				// curLine = []string{fmt.Sprintf("%d", opcode)}
				curLine = []string{}
				break
			}
			curLine = []string{"ADD", FormattedMem(memory, i+1, mod[0]),
				FormattedMem(memory, i+2, mod[1]), FormattedMem(memory, i+3, mod[2])}
			i += 3
		case vm.OPCODE_MUL:
			if i+3 >= len(memory) {
				// curLine = []string{fmt.Sprintf("%d", opcode)}
				curLine = []string{}
				break
			}
			curLine = []string{"MUL", FormattedMem(memory, i+1, mod[0]),
				FormattedMem(memory, i+2, mod[1]), FormattedMem(memory, i+3, mod[2])}
			i += 3
		case vm.OPCODE_INPUT:
			if i+1 >= len(memory) {
				// curLine = []string{fmt.Sprintf("%d", opcode)}
				curLine = []string{}
				break
			}
			curLine = []string{"IN", FormattedMem(memory, i+1, mod[0])}
			i++
		case vm.OPCODE_OUTPUT:
			if i+1 >= len(memory) {
				// curLine = []string{fmt.Sprintf("%d", opcode)}
				curLine = []string{}
				break
			}
			curLine = []string{"OUT", FormattedMem(memory, i+1, mod[0])}
			i++
		case vm.OPCODE_JMP_T:
			if i+2 >= len(memory) {
				// curLine = []string{fmt.Sprintf("%d", opcode)}
				curLine = []string{}
				break
			}
			curLine = []string{"JZ", FormattedMem(memory, i+1, mod[0]),
				FormattedMem(memory, i+2, mod[1])}
			i += 2
		case vm.OPCODE_JMP_F:
			if i+2 >= len(memory) {
				// curLine = []string{fmt.Sprintf("%d", opcode)}
				curLine = []string{}
				break
			}
			curLine = []string{"JNZ", FormattedMem(memory, i+1, mod[0]),
				FormattedMem(memory, i+2, mod[1])}
			i += 2
		case vm.OPCODE_LESS_THAN:
			if i+3 >= len(memory) {
				// curLine = []string{fmt.Sprintf("%d", opcode)}
				curLine = []string{}
				break
			}
			curLine = []string{"SLT", FormattedMem(memory, i+1, mod[0]),
				FormattedMem(memory, i+2, mod[1]), FormattedMem(memory, i+3, mod[2])}
			i += 3
		case vm.OPCODE_EQUALS:
			if i+3 >= len(memory) {
				// curLine = []string{fmt.Sprintf("%d", opcode)}
				curLine = []string{}
				break
			}
			curLine = []string{"SEQ", FormattedMem(memory, i+1, mod[0]),
				FormattedMem(memory, i+2, mod[1]), FormattedMem(memory, i+3, mod[2])}
			i += 3
		case vm.OPCODE_INC_RELV:
			if i+1 >= len(memory) {
				// curLine = []string{fmt.Sprintf("%d", opcode)}
				curLine = []string{}
				break
			}
			curLine = []string{"INCB", FormattedMem(memory, i+1, mod[0])}
			i++
		case vm.OPCODE_HALT:
			curLine = []string{"HALT"}
		default:
			// curLine = []string{fmt.Sprintf("%d", opcode)}
			curLine = []string{}

			// return lines, fmt.Errorf("Unknown opcode: %d at index: %d", opcode, i)
		}
		lines = append(lines, curLine)
	}
	return lines, activeInd, nil
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
