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
}

type AssembleOpcode struct {
	Opcode      int
	Mode        []int
	Params      []int
	FinalOpcode int
}

var lastAdrr int
var variables map[string]int
var labels map[string]int

func Start(filePath string) error {
	data, _ := os.ReadFile(filePath)

	lines := strings.Split(string(data), "\n")

	// Tokenizer
	tokens := []string{}
	for _, l := range lines {
		tk := strings.Fields(l)
		if len(tk) > 0 {
			tokens = append(tokens, tk...)
		}
	}
	// Preprocessor #1 (Locate labels addrs)
	labels = make(map[string]int)
	variables = make(map[string]int)
	toBeRemovedLabels := []string{}
	vCnt := 0

	for i, t := range tokens {
		if strings.HasSuffix(t, ":") {
			labels[t[:len(t)-1]] = i - len(labels)
			toBeRemovedLabels = append(toBeRemovedLabels, t)
			continue
		}

		if isOpcode(t) || !isVariable(t) {
			continue
		}

		if _, ok := variables[t]; !ok {
			variables[t] = vCnt
			vCnt++
		}

	}

	fmt.Printf("variables: %v\n", variables)

	// Preprocessor #2 (Remove label declarations)
	newTokens := slices.DeleteFunc(tokens, func(e string) bool {
		return slices.Contains(toBeRemovedLabels, e)
	})

	oldTokens := slices.Clone(tokens)
	tokens = newTokens
	_ = oldTokens
	// Preprocessor #3 (Replace labels with actual Addrs)
	fmt.Printf("labels: %v\n", labels)
	// for i, t := range tokens {
	// 	if strings.HasPrefix(t, ":") {
	// 		ind, ok := labels[t[1:]]
	// 		if !ok {
	// 			fmt.Printf("Invalid Label %s", t)
	// 			return fmt.Errorf("invalid Label %s", t)
	// 		}

	// 		// tokens[i] = fmt.Sprintf("%d->%s", ind, tokens[ind])
	// 		tokens[i] = strconv.Itoa(ind)
	// 		continue
	// 	}
	// }

	lastAdrr = len(tokens)
	fmt.Printf("tokens: %v\n", tokens)

	// Preprocessor #4 (Replace variables with actual memory addr)

	// Tokens to bytecodes
	byteCode := []int{}
	for i := 0; i < len(tokens); i++ {

		t := tokens[i]

		switch t {
		case "ADD":
			op := parseOpcodeLine(vm.OPCODE_ADD, tokens[i+1:i+4])

			// if err != nil || err1 != nil || err2 != nil {
			// 	// return fmt.Errorf("invalid ADD instruction")
			// }
			byteCode = append(byteCode, op.FinalOpcode)
			byteCode = append(byteCode, op.Params...)
			i += 3
		case "MUL":
			op := parseOpcodeLine(vm.OPCODE_MUL, tokens[i+1:i+4])

			// if err != nil || err1 != nil || err2 != nil {
			// 	// return fmt.Errorf("invalid MUL instruction")
			// }
			byteCode = append(byteCode, op.FinalOpcode)
			byteCode = append(byteCode, op.Params...)

			i += 3

		case "IN":
			op := parseOpcodeLine(vm.OPCODE_INPUT, tokens[i+1:i+2])

			// if err != nil {
			// 	// return fmt.Errorf("invalid IN instruction %v", tokens[i+1])
			// }
			byteCode = append(byteCode, op.FinalOpcode)
			byteCode = append(byteCode, op.Params...)

			i++
		case "OUT":
			op := parseOpcodeLine(vm.OPCODE_OUTPUT, tokens[i+1:i+2])

			// if err != nil {
			// 	// return fmt.Errorf("invalid OUT instruction")
			// }

			byteCode = append(byteCode, op.FinalOpcode)
			byteCode = append(byteCode, op.Params...)

			i++
		case "JNZ":
			op := parseOpcodeLine(vm.OPCODE_JMP_T, tokens[i+1:i+3])

			// if err != nil || err1 != nil {
			// 	// return fmt.Errorf("invalid JNZ instruction")
			// }
			byteCode = append(byteCode, op.FinalOpcode)
			byteCode = append(byteCode, op.Params...)
			i += 2
		case "JZ":
			op := parseOpcodeLine(vm.OPCODE_JMP_F, tokens[i+1:i+3])
			// if err != nil || err1 != nil {
			// 	// return fmt.Errorf("invalid JZ instruction")
			// }
			byteCode = append(byteCode, op.FinalOpcode)
			byteCode = append(byteCode, op.Params...)
			i += 2
		case "SLT":
			op := parseOpcodeLine(vm.OPCODE_LESS_THAN, tokens[i+1:i+4])

			// if err != nil || err1 != nil || err2 != nil {
			// 	// return fmt.Errorf("invalid SLT instruction")
			// }
			byteCode = append(byteCode, op.FinalOpcode)
			byteCode = append(byteCode, op.Params...)
			i += 3

		case "SEQ":
			op := parseOpcodeLine(vm.OPCODE_EQUALS, tokens[i+1:i+4])

			// if err != nil || err1 != nil || err2 != nil {
			// 	// return fmt.Errorf("invalid SEQ instruction")
			// }
			byteCode = append(byteCode, op.FinalOpcode)
			byteCode = append(byteCode, op.Params...)

			i += 3

		case "INCR":
			op := parseOpcodeLine(vm.OPCODE_INC_RELV, tokens[i+1:i+2])

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

	fmt.Printf("byteCode: %v\n", byteCode)
	return nil
}

func isOpcode(val string) bool {
	opcodes := []string{"ADD", "MUL", "IN", "OUT", "JNZ", "JZ", "SLT", "SEQ", "INCR", "HALT"}

	return slices.Contains(opcodes, val)
}

func isVariable(val string) bool {
	var validIdentifier = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
	return validIdentifier.MatchString(val)
}

func parseOpcodeLine(opcode int, params []string) AssembleOpcode {
	op := AssembleOpcode{
		Opcode: opcode,
		Params: []int{},
		Mode:   []int{},
	}
	mode := strings.Builder{}
	for mi := 0; mi < len(params); mi++ {

		if q, ok := variables[params[mi]]; ok {
			// tokens[i+mi] =
			op.Params = append(op.Params, lastAdrr+q)
			mode.WriteRune('0')
			op.Mode = append(op.Mode, 0)
		} else if strings.HasPrefix(params[mi], ":") {
			ind, ok := labels[params[mi][1:]]
			if !ok {
				fmt.Printf("Invalid Label %s", params[mi])
				continue
				// return fmt.Errorf("invalid Label %s", params[mi])
			}

			op.Params = append(op.Params, ind)
			mode.WriteRune('1')
			op.Mode = append(op.Mode, 1)
			// tokens[i] = fmt.Sprintf("%d->%s", ind, tokens[ind])
			// tokens[i] = strconv.Itoa(ind)

		} else {
			a, err := strconv.Atoi(params[mi])
			if err != nil {
				fmt.Printf("IDK SOMETING WENT WRONG %d %s", err, params[mi])
			}
			op.Params = append(op.Params, a)
			mode.WriteRune('1')
			op.Mode = append(op.Mode, 1)
		}
	}
	op.Opcode = opcode

	// mode.WriteString(fmt.Sprintf("0%d", op.Opcode))
	m := mode.String()
	revMode := ""
	for _, v := range m {
		revMode = string(v) + revMode
	}
	fo, err := strconv.Atoi(revMode + fmt.Sprintf("0%d", op.Opcode))
	op.FinalOpcode = fo

	if err != nil {
		fmt.Printf("ERROR %v\n", err)
	}

	fmt.Printf("OP: %d | %v | %v | %v\n", opcode, op.Params, op.Mode, op.FinalOpcode)
	return op
}
