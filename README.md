# Intcode Assember, Visualizer & VM

Assembler, Visualizer and VM for the [Intcode](https://adventofcode.com/2019/day/9) from Advent of Code 2019.
Written in Golang just for fun.

There are four packages in this project:

- `vm`: An intcode VM that can run intcode files.
- `assembler`: An intcode assembler that can assemble programs from intcode assembly.
- `debugger`: Visualizer that shows formatted intcode, memory and output.

## Building

```
$ git clone https://github.com/xSaCh/intcode/
$ cd intcode
$ go run .
```
### Testing VM
```
$ cd vm
$ go test -v
```
## Assembler
### Instructions
| Operation  |   | Note  |
| :------------ | :------------: | :------------: |
| ADD a b c | c = a + b  |   |
| MUL a b c | c = a * b  |   |
| IN a |   | Read value and store it in 'a' |
| OUT a  |   | Output 'a' value  |
| JNZ  a label | if a !=0; goto label  |   |
| JZ  a label | if a ==0; goto label  |   |
| SLT a b c | c = a < b  |   |
|  EQ a b c | c = a == b  |   |
|  INCR a | rel_reg  += a  | Increment relative base register by 'a' |
|  HALT |   | Stop program execution |

### Example
Factorial of n
```asm
IN      n                   ; Read n
SLT     n 1   can_exit      ;
JNZ     can_exit :exit      ; Jump to label 'exit' if 'can_exit' is true
ADD     n 0   i 

loop:
ADD     i -1  i
JZ      i :print            ; Jump if 'i' == 0
MUL     n i n
JNZ     i :loop             ; Jump to label 'loop' till 'i' != 0

print:
OUT     n                   ; Output value of 'n'
HALT                        ; Stop program    

exit:
OUT     1
HALT
```
