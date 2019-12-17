package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func handleModeRead(commands []int64, param int64, mode int, position int) int64 {
	if mode == 0 { //reference
		return commands[param]
	} else if mode == 1 {
		return param
	} else if mode == 2 {
		return commands[int64(position)+param]
	} else {
		log.Fatal("Unknown mode:", mode)
	}
	return 0
}

func handleModeWrite(commands []int64, param int64, mode int, position int, result int64) {
	if mode == 0 { //reference
		commands[param] = result
	} else if mode == 1 {
		log.Fatal("Illegal Write")
	} else if mode == 2 {
		commands[int64(position)+param] = result
	} else {
		log.Fatal("Unknown write mode:", mode)
	}
}

func do1Command(commands []int64, startPosition int, modes []int, position int) int {
	parameter1 := commands[startPosition+1]
	parameter2 := commands[startPosition+2]
	parameter3 := commands[startPosition+3]
	sum := handleModeRead(commands, parameter1, modes[0], position) + handleModeRead(commands, parameter2, modes[1], position)
	handleModeWrite(commands, parameter3, modes[2], position, sum)
	return startPosition + 4
}

func do2Command(commands []int64, startPosition int, modes []int, position int) int {
	parameter1 := commands[startPosition+1]
	parameter2 := commands[startPosition+2]
	parameter3 := commands[startPosition+3]
	prod := handleModeRead(commands, parameter1, modes[0], position) * handleModeRead(commands, parameter2, modes[1], position)
	handleModeWrite(commands, parameter3, modes[2], position, prod)
	return startPosition + 4
}

func do3Command(commands []int64, startPosition int, modes []int, position int, input chan int, state chan int) int {
	parameter1 := commands[startPosition+1]
	state <- 3
	var i = <-input
	handleModeWrite(commands, parameter1, modes[0], position, int64(i))
	return startPosition + 2
}

func do4Command(commands []int64, startPosition int, modes []int, position int, output chan int64, state chan int) int {
	parameter1 := commands[startPosition+1]
	state <- 4
	output <- handleModeRead(commands, parameter1, modes[0], position)
	return startPosition + 2
}

func do5Command(commands []int64, startPosition int, modes []int, position int) int {
	parameter1 := commands[startPosition+1]
	parameter2 := commands[startPosition+2]
	if handleModeRead(commands, parameter1, modes[0], position) > 0 {
		return int(handleModeRead(commands, parameter2, modes[1], position))
	} else {
		return startPosition + 3
	}
}
func do6Command(commands []int64, startPosition int, modes []int, position int) int {
	parameter1 := commands[startPosition+1]
	parameter2 := commands[startPosition+2]
	if handleModeRead(commands, parameter1, modes[0], position) == 0 {
		return int(handleModeRead(commands, parameter2, modes[1], position))
	} else {
		return startPosition + 3
	}
}
func do7Command(commands []int64, startPosition int, modes []int, position int) int {
	parameter1 := commands[startPosition+1]
	parameter2 := commands[startPosition+2]
	parameter3 := commands[startPosition+3]
	if handleModeRead(commands, parameter1, modes[0], position) < handleModeRead(commands, parameter2, modes[1], position) {
		handleModeWrite(commands, parameter3, modes[2], position, 1)
	} else {
		handleModeWrite(commands, parameter3, modes[2], position, 0)
	}
	return startPosition + 4
}
func do8Command(commands []int64, startPosition int, modes []int, position int) int {
	parameter1 := commands[startPosition+1]
	parameter2 := commands[startPosition+2]
	parameter3 := commands[startPosition+3]
	if handleModeRead(commands, parameter1, modes[0], position) == handleModeRead(commands, parameter2, modes[1], position) {
		handleModeWrite(commands, parameter3, modes[2], position, 1)
	} else {
		handleModeWrite(commands, parameter3, modes[2], position, 0)
	}
	return startPosition + 4
}

func do9Command(commands []int64, startPosition int, modes []int, position int) (int, int) {
	parameter1 := commands[startPosition+1]
	//fmt.Printf("Changing Position from %d to %d\n", position, handleModeRead(commands, parameter1, modes[0], position))
	return startPosition + 2, position + int(handleModeRead(commands, parameter1, modes[0], position))
}

func parseCommand(command int64) (opCode int, mode []int) {
	code := command % 100
	modes := command / 100
	paramMode := make([]int, 3)
	for i := 0; i < 3; i++ {
		paramMode[i] = int(modes % 10)
		modes = modes / 10
	}
	return int(code), paramMode
}

func doCommands(commands []int64, input chan int, output chan int64, state chan int) {
	instructionPointer := 0
	position := 0
	for {
		opCode, modes := parseCommand(commands[instructionPointer])
		switch opCode {
		case 1:
			instructionPointer = do1Command(commands, instructionPointer, modes, position)
		case 2:
			instructionPointer = do2Command(commands, instructionPointer, modes, position)
		case 3:
			instructionPointer = do3Command(commands, instructionPointer, modes, position, input, state)
		case 4:
			instructionPointer = do4Command(commands, instructionPointer, modes, position, output, state)
		case 5:
			instructionPointer = do5Command(commands, instructionPointer, modes, position)
		case 6:
			instructionPointer = do6Command(commands, instructionPointer, modes, position)
		case 7:
			instructionPointer = do7Command(commands, instructionPointer, modes, position)
		case 8:
			instructionPointer = do8Command(commands, instructionPointer, modes, position)
		case 9:
			instructionPointer, position = do9Command(commands, instructionPointer, modes, position)
		case 99:
			close(state)
			return
		}
	}
}

func toInt(input []string) []int64 {
	data := make([]int64, len(input))
	for index, code := range input {
		n, err := strconv.ParseInt(code, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		data[index] = n
	}
	return data
}

func main() {
	argsWithoutProg := os.Args[1:]
	file, err := os.Open(argsWithoutProg[0])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	input := make(chan int)
	output := make(chan int64)
	state := make(chan int)
	data := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data += scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	memory := toInt(strings.Split(data, ","))
	memory = append(memory, make([]int64, 10000)...)
	memory[0] = 2
	go doCommands(memory, input, output, state)
	handleIO(input, output, state)
	close(output)
	close(input)
}

func handleIO(input chan int, output chan int64, state chan int) {

	main := []int{A, COMMA, A, COMMA, B, COMMA, C, COMMA, B, COMMA, A, COMMA, C, COMMA, B, COMMA, C, COMMA, A, NEWLINE}
	//L,6,R,12,L,6,L,8,L,8,8NEWLINE
	a := []int{L, COMMA, 54, COMMA, R, COMMA, 49, 50, COMMA, L, COMMA, 54, COMMA, L, COMMA, 56, COMMA, L, COMMA, 56, NEWLINE}
	//L,6,R,12,R,8,L,8,NEWLINE
	b := []int{L, COMMA, 54, COMMA, R, COMMA, 49, 50, COMMA, R, COMMA, 56, COMMA, L, COMMA, 56, NEWLINE}
	//L,4,L,4,L,6,NEWLINE
	c := []int{L, COMMA, 52, COMMA, L, COMMA, 52, COMMA, L, COMMA, 54, NEWLINE}
	noUpdate := []int{110, NEWLINE}
	prog := 0
	pos := 0
	for {
		progState, ok := <-state
		if !ok {
			break
		}
		if progState == 3 {
			var code int
			switch prog {
			case 0:
				code = main[pos]
				break
			case 1:
				code = a[pos]
				break
			case 2:
				code = b[pos]
				break
			case 3:
				code = c[pos]
				break
			case 4:
				code = noUpdate[pos]
				break
			}
			pos++
			if code == NEWLINE {
				prog++
				pos = 0
			}
			input <- code
		} else if progState == 4 {
			out := <-output
			if out > 128 {
				fmt.Println(out)
			}
		}
	}
}

const (
	NEWLINE int = 10
	A       int = 65
	B       int = 66
	C       int = 67
	L       int = 76
	R       int = 82
	COMMA   int = 44
)
