package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type input func() int
type output func(int64)
type result func() int64

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

func do3Command(commands []int64, startPosition int, modes []int, position int, in input) int {
	parameter1 := commands[startPosition+1]
	handleModeWrite(commands, parameter1, modes[0], position, int64(in()))
	return startPosition + 2
}

func do4Command(commands []int64, startPosition int, modes []int, position int, out output) int {
	parameter1 := commands[startPosition+1]
	out(handleModeRead(commands, parameter1, modes[0], position))
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

func doCommands(commands []int64, in input, out output) {
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
			instructionPointer = do3Command(commands, instructionPointer, modes, position, in)
		case 4:
			instructionPointer = do4Command(commands, instructionPointer, modes, position, out)
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
	data := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data += scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	memory := toInt(strings.Split(data, ","))
	memory = append(memory, make([]int64, 1000)...)
	testMemory := append(memory[:0:0], memory...)
	in, out := handleIO()
	doCommands(testMemory, in, out)
}

func handleIO() (input, output) {
	//AND := []int{65, 78, 68}
	//OR := []int{79, 78, 82}
	//NOT := []int{78, 79, 84}
	//WALK := []int{87, 65, 76, 75}
	//RUN := []int{82, 85, 78}
	pos := 0

	//!(A && B && C) && D
	//program := []int{78, 79, 84, SPACE, J, SPACE, T, NEWLINE, 65, 78, 68, SPACE, A, SPACE, T, NEWLINE,
	//	65, 78, 68, SPACE, B, SPACE, T, NEWLINE, 65, 78, 68, SPACE, C, SPACE, T, NEWLINE,
	//	78, 79, 84, SPACE, T, SPACE, J, NEWLINE, 65, 78, 68, SPACE, D, SPACE, J, NEWLINE,
	//	87, 65, 76, 75, NEWLINE}

	//!(A && B && C) && D && (E || H)
	program := []int{
		78, 79, 84, SPACE, J, SPACE, T, NEWLINE, 65, 78, 68, SPACE, A, SPACE, T, NEWLINE,
		65, 78, 68, SPACE, B, SPACE, T, NEWLINE, 65, 78, 68, SPACE, C, SPACE, T, NEWLINE,
		78, 79, 84, SPACE, T, SPACE, J, NEWLINE, 65, 78, 68, SPACE, D, SPACE, J, NEWLINE,
		78, 79, 84, SPACE, E, SPACE, T, NEWLINE, 78, 79, 84, SPACE, T, SPACE, T, NEWLINE,
		79, 82, SPACE, H, SPACE, T, NEWLINE, 65, 78, 68, SPACE, T, SPACE, J, NEWLINE,
		82, 85, 78, NEWLINE}

	in := func() int {
		instruct := program[pos]
		pos++
		return instruct
	}

	out := func(r int64) {
		if r >= 256 {
			fmt.Println("Damage = ", r)
		} else {
			fmt.Print(string(rune(r)))
		}
	}
	return in, out
}

const (
	A       int = 65
	B       int = 66
	C       int = 67
	D       int = 68
	E       int = 69
	F       int = 70
	G       int = 71
	H       int = 72
	I       int = 73
	T       int = 84
	J       int = 74
	SPACE   int = 32
	NEWLINE int = 10
)
