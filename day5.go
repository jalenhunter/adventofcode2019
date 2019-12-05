package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func handleMode(commands []int, param int, mode int) int {
	if mode == 0 { //reference
		return commands[param]
	} else if mode == 1 {
		return param
	} else {
		log.Fatal("Unknown mode:", mode)
	}
	return 0
}

func do1Command(commands []int, startPosition int, modes []int) int {
	parameter1 := commands[startPosition+1]
	parameter2 := commands[startPosition+2]
	parameter3 := commands[startPosition+3]
	sum := handleMode(commands, parameter1, modes[0]) + handleMode(commands, parameter2, modes[1])
	commands[parameter3] = sum //writes are always mode 0
	return 4
}

func do2Command(commands []int, startPosition int, modes []int) int {
	parameter1 := commands[startPosition+1]
	parameter2 := commands[startPosition+2]
	parameter3 := commands[startPosition+3]
	prod := handleMode(commands, parameter1, modes[0]) * handleMode(commands, parameter2, modes[1])
	commands[parameter3] = prod //writes are always mode 0
	return 4
}

func do3Command(commands []int, startPosition int, _ []int) int {
	parameter1 := commands[startPosition+1]
	var i int
	_, err := fmt.Scanf("%d", &i)
	if err != nil {
		log.Fatal(err)
	}
	commands[parameter1] = i //writes are always mode 0
	return 2
}

func do4Command(commands []int, startPosition int, modes []int) int {
	parameter1 := commands[startPosition+1]
	fmt.Println(handleMode(commands, parameter1, modes[0]))
	return 2
}

func do5Command(commands []int, startPosition int, modes []int) (bool, int) {
	parameter1 := commands[startPosition+1]
	parameter2 := commands[startPosition+2]
	if handleMode(commands, parameter1, modes[0]) > 0 {
		return true, handleMode(commands, parameter2, modes[1])
	} else {
		return false, 3
	}
}
func do6Command(commands []int, startPosition int, modes []int) (bool, int) {
	parameter1 := commands[startPosition+1]
	parameter2 := commands[startPosition+2]
	if handleMode(commands, parameter1, modes[0]) == 0 {
		return true, handleMode(commands, parameter2, modes[1])
	} else {
		return false, 3
	}
}
func do7Command(commands []int, startPosition int, modes []int) int {
	parameter1 := commands[startPosition+1]
	parameter2 := commands[startPosition+2]
	parameter3 := commands[startPosition+3]
	if handleMode(commands, parameter1, modes[0]) < handleMode(commands, parameter2, modes[1]) {
		commands[parameter3] = 1
	} else {
		commands[parameter3] = 0
	}
	return 4
}
func do8Command(commands []int, startPosition int, modes []int) int {
	parameter1 := commands[startPosition+1]
	parameter2 := commands[startPosition+2]
	parameter3 := commands[startPosition+3]
	if handleMode(commands, parameter1, modes[0]) == handleMode(commands, parameter2, modes[1]) {
		commands[parameter3] = 1
	} else {
		commands[parameter3] = 0
	}
	return 4
}

func parseCommand(command int) (opCode int, mode []int) {
	code := command % 100
	modes := command / 100
	paramMode := make([]int, 3)
	for i := 0; i < 3; i++ {
		paramMode[i] = modes % 10
		modes = modes / 10
	}
	return code, paramMode
}

func doCommands(commands []int) {
	instructionPointer := 0
	for {
		opCode, modes := parseCommand(commands[instructionPointer])
		switch opCode {
		case 1:
			instructionPointer += do1Command(commands, instructionPointer, modes)
		case 2:
			instructionPointer += do2Command(commands, instructionPointer, modes)
		case 3:
			instructionPointer += do3Command(commands, instructionPointer, modes)
		case 4:
			instructionPointer += do4Command(commands, instructionPointer, modes)
		case 5:
			jump, value := do5Command(commands, instructionPointer, modes)
			if jump {
				instructionPointer = value
			} else {
				instructionPointer += value
			}
		case 6:
			jump, value := do6Command(commands, instructionPointer, modes)
			if jump {
				instructionPointer = value
			} else {
				instructionPointer += value
			}
		case 7:
			instructionPointer += do7Command(commands, instructionPointer, modes)
		case 8:
			instructionPointer += do8Command(commands, instructionPointer, modes)
		case 99:
			return
		}
	}
}

func toInt(input []string) []int {
	data := make([]int, len(input))
	for index, code := range input {
		i, err := strconv.Atoi(code)
		if err != nil {
			log.Fatal(err)
		}
		data[index] = i
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
	doCommands(memory)
}
