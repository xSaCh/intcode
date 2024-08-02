package assembler

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/xSaCh/intcode/vm"
)

type Assembler struct {
	Tokens    []string
	Labels    map[string]int
	Variables map[string]int
	ByteCode  []int

	ParsedOpcodes []AssembleOpcode

	lastAddr          int
	toBeRemovedLabels []string
}

type AssembleOpcode struct {
	Opcode      int
	Mode        []int
	Params      []int
	FinalOpcode int
}

func NewAssembler() *Assembler {
	return &Assembler{
		Tokens:    []string{},
		Labels:    map[string]int{},
		Variables: map[string]int{},
		ByteCode:  []int{},

		ParsedOpcodes: []AssembleOpcode{},

		lastAddr:          0,
		toBeRemovedLabels: []string{},
	}
}

func (a *Assembler) AssembleFromFile(filePath string) error {
	data, _ := os.ReadFile(filePath)

	return a.Assemble(string(data))
}

func (a *Assembler) Assemble(data string) error {
	a.tokenize(data)
	a.symbolGenerator()

	// Preprocessor #2 (Remove label declarations)
	a.Tokens = slices.DeleteFunc(a.Tokens, func(e string) bool {
		return slices.Contains(a.toBeRemovedLabels, e)
	})

	a.lastAddr = len(a.Tokens)

	a.translation()
	return nil
}

// Tokens to bytecodes
func (a *Assembler) translation() error {
	byteCode := []int{}
	for i := 0; i < len(a.Tokens); i++ {

		t := a.Tokens[i]

		switch t {
		case "ADD":
			op := a.parseOpcodeLine(vm.OPCODE_ADD, a.Tokens[i+1:i+4])

			// if err != nil || err1 != nil || err2 != nil {
			// 	// return fmt.Errorf("invalid ADD instruction")
			// }
			byteCode = append(byteCode, op.FinalOpcode)
			byteCode = append(byteCode, op.Params...)
			i += 3
		case "MUL":
			op := a.parseOpcodeLine(vm.OPCODE_MUL, a.Tokens[i+1:i+4])

			// if err != nil || err1 != nil || err2 != nil {
			// 	// return fmt.Errorf("invalid MUL instruction")
			// }
			byteCode = append(byteCode, op.FinalOpcode)
			byteCode = append(byteCode, op.Params...)

			i += 3

		case "IN":
			op := a.parseOpcodeLine(vm.OPCODE_INPUT, a.Tokens[i+1:i+2])

			// if err != nil {
			// 	// return fmt.Errorf("invalid IN instruction %v", a.Tokens[i+1])
			// }
			byteCode = append(byteCode, op.FinalOpcode)
			byteCode = append(byteCode, op.Params...)

			i++
		case "OUT":
			op := a.parseOpcodeLine(vm.OPCODE_OUTPUT, a.Tokens[i+1:i+2])

			// if err != nil {
			// 	// return fmt.Errorf("invalid OUT instruction")
			// }

			byteCode = append(byteCode, op.FinalOpcode)
			byteCode = append(byteCode, op.Params...)

			i++
		case "JNZ":
			op := a.parseOpcodeLine(vm.OPCODE_JMP_T, a.Tokens[i+1:i+3])

			// if err != nil || err1 != nil {
			// 	// return fmt.Errorf("invalid JNZ instruction")
			// }
			byteCode = append(byteCode, op.FinalOpcode)
			byteCode = append(byteCode, op.Params...)
			i += 2
		case "JZ":
			op := a.parseOpcodeLine(vm.OPCODE_JMP_F, a.Tokens[i+1:i+3])
			// if err != nil || err1 != nil {
			// 	// return fmt.Errorf("invalid JZ instruction")
			// }
			byteCode = append(byteCode, op.FinalOpcode)
			byteCode = append(byteCode, op.Params...)
			i += 2
		case "SLT":
			op := a.parseOpcodeLine(vm.OPCODE_LESS_THAN, a.Tokens[i+1:i+4])

			// if err != nil || err1 != nil || err2 != nil {
			// 	// return fmt.Errorf("invalid SLT instruction")
			// }
			byteCode = append(byteCode, op.FinalOpcode)
			byteCode = append(byteCode, op.Params...)
			i += 3

		case "SEQ":
			op := a.parseOpcodeLine(vm.OPCODE_EQUALS, a.Tokens[i+1:i+4])

			// if err != nil || err1 != nil || err2 != nil {
			// 	// return fmt.Errorf("invalid SEQ instruction")
			// }
			byteCode = append(byteCode, op.FinalOpcode)
			byteCode = append(byteCode, op.Params...)

			i += 3

		case "INCR":
			op := a.parseOpcodeLine(vm.OPCODE_INC_RELV, a.Tokens[i+1:i+2])

			// if err != nil {
			// 	// return fmt.Errorf("invalid INCR instruction")
			// }
			byteCode = append(byteCode, op.FinalOpcode)
			byteCode = append(byteCode, op.Params...)
			i++
		case "HALT":
			byteCode = append(byteCode, vm.OPCODE_HALT)
		}
	}

	a.ByteCode = byteCode
	return nil
}

func (a *Assembler) tokenize(data string) {
	lines := strings.Split(string(data), "\n")

	for _, l := range lines {
		tk := strings.Fields(l)
		if len(tk) > 0 {
			a.Tokens = append(a.Tokens, tk...)
		}
	}
}

func (a *Assembler) symbolGenerator() {
	vCnt := 0

	for i, t := range a.Tokens {
		if strings.HasSuffix(t, ":") {
			a.Labels[t[:len(t)-1]] = i - len(a.Labels)
			a.toBeRemovedLabels = append(a.toBeRemovedLabels, t)
			continue
		}

		if isOpcode(t) || !isVariable(t) {
			continue
		}

		if _, ok := a.Variables[t]; !ok {
			a.Variables[t] = vCnt
			vCnt++
		}
	}

}

func (a *Assembler) parseOpcodeLine(opcode int, params []string) AssembleOpcode {
	op := AssembleOpcode{
		Opcode:      opcode,
		Params:      []int{},
		Mode:        make([]int, 0, 3),
		FinalOpcode: 0,
	}
	for mi := 0; mi < len(params); mi++ {

		// Memory Mode
		if q, ok := a.Variables[params[mi]]; ok {
			op.Params = append(op.Params, a.lastAddr+q)
			op.Mode = append(op.Mode, 0)
		} else if strings.HasPrefix(params[mi], ":") { // Label Mode
			ind, ok := a.Labels[params[mi][1:]]
			if !ok {
				fmt.Printf("Invalid Label %s", params[mi])
				continue
			}

			op.Params = append(op.Params, ind)
			op.Mode = append(op.Mode, 1)

		} else { // Immediate Mode
			a, err := strconv.Atoi(params[mi])
			if err != nil {
				fmt.Printf("IDK SOMETING WENT WRONG %d %s", err, params[mi])
			}
			op.Params = append(op.Params, a)
			op.Mode = append(op.Mode, 1)
		}
	}
	op.Opcode = opcode

	for i := len(op.Mode) - 1; i >= 0; i-- {
		op.FinalOpcode = op.FinalOpcode*10 + op.Mode[i]
	}
	op.FinalOpcode = op.FinalOpcode*100 + op.Opcode
	// fmt.Printf("OP: %d | %v | %v | %v\n", opcode, op.Params, op.Mode, op.FinalOpcode)

	a.ParsedOpcodes = append(a.ParsedOpcodes, op)
	return op
}

func isOpcode(val string) bool {
	opcodes := []string{"ADD", "MUL", "IN", "OUT", "JNZ", "JZ", "SLT", "SEQ", "INCR", "HALT"}

	return slices.Contains(opcodes, val)
}

func isVariable(val string) bool {
	var validIdentifier = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
	return validIdentifier.MatchString(val)
}
