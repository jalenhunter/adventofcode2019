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

func do3Command(commands []int64, startPosition int, modes []int, position int) int {
	parameter1 := commands[startPosition+1]
	var i int
	fmt.Print("Enter a number: ")
	_, err := fmt.Scanf("%d", &i)
	if err != nil {
		log.Fatal(err)
	}
	handleModeWrite(commands, parameter1, modes[0], position, int64(i))
	return startPosition + 2
}

func do4Command(commands []int64, startPosition int, modes []int, position int) int {
	parameter1 := commands[startPosition+1]
	fmt.Println(handleModeRead(commands, parameter1, modes[0], position))
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

func doCommands(commands []int64) {
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
			instructionPointer = do3Command(commands, instructionPointer, modes, position)
		case 4:
			instructionPointer = do4Command(commands, instructionPointer, modes, position)
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
	// overwrite per instructions
	memory = append(memory, make([]int64, 1000)...)
	doCommands(memory)
}
