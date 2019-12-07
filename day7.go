package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var silent = true

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
	return startPosition + 4
}

func do2Command(commands []int, startPosition int, modes []int) int {
	parameter1 := commands[startPosition+1]
	parameter2 := commands[startPosition+2]
	parameter3 := commands[startPosition+3]
	prod := handleMode(commands, parameter1, modes[0]) * handleMode(commands, parameter2, modes[1])
	commands[parameter3] = prod //writes are always mode 0
	return startPosition + 4
}

func do3Command(commands []int, startPosition int, _ []int) int {
	if !silent {
		parameter1 := commands[startPosition+1]
		var i int
		fmt.Print("Enter a number: ")
		_, err := fmt.Scanf("%d", &i)
		if err != nil {
			log.Fatal(err)
		}
		commands[parameter1] = i //writes are always mode 0
	}
	return startPosition + 2
}

func do4Command(commands []int, startPosition int, modes []int) int {
	if !silent {
		parameter1 := commands[startPosition+1]
		fmt.Println(handleMode(commands, parameter1, modes[0]))
	}
	return startPosition + 2
}

func do5Command(commands []int, startPosition int, modes []int) int {
	parameter1 := commands[startPosition+1]
	parameter2 := commands[startPosition+2]
	if handleMode(commands, parameter1, modes[0]) > 0 {
		return handleMode(commands, parameter2, modes[1])
	} else {
		return startPosition + 3
	}
}
func do6Command(commands []int, startPosition int, modes []int) int {
	parameter1 := commands[startPosition+1]
	parameter2 := commands[startPosition+2]
	if handleMode(commands, parameter1, modes[0]) == 0 {
		return handleMode(commands, parameter2, modes[1])
	} else {
		return startPosition + 3
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
	return startPosition + 4
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
	return startPosition + 4
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
			instructionPointer = do1Command(commands, instructionPointer, modes)
		case 2:
			instructionPointer = do2Command(commands, instructionPointer, modes)
		case 3:
			instructionPointer = do3Command(commands, instructionPointer, modes)
		case 4:
			instructionPointer = do4Command(commands, instructionPointer, modes)
		case 5:
			instructionPointer = do5Command(commands, instructionPointer, modes)
		case 6:
			instructionPointer = do6Command(commands, instructionPointer, modes)
		case 7:
			instructionPointer = do7Command(commands, instructionPointer, modes)
		case 8:
			instructionPointer = do8Command(commands, instructionPointer, modes)
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

func Perm(a []rune, f func([]rune)) {
	perm(a, f, 0)
}

// Permute the values at index i to len(a)-1.
func perm(a []rune, f func([]rune), i int) {
	if i > len(a) {
		f(a)
		return
	}
	perm(a, f, i+1)
	for j := i + 1; j < len(a); j++ {
		a[i], a[j] = a[j], a[i]
		perm(a, f, i+1)
		a[i], a[j] = a[j], a[i]
	}
}

func createPermutations() []string {
	var result []string
	Perm([]rune("01234"), func(a []rune) {
		result = append(result, string(a))
	})
	return result
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
	maxThrust := 0
	var bestSetting string
	for _, permutation := range createPermutations() {
		inputSignal := 0
		for _, phase := range permutation {
			testMemory := append(memory[:0:0], memory...)
			v, err := strconv.Atoi(string(phase))
			if err != nil {
				log.Fatal(err)
			}
			testMemory[8] = v
			testMemory[9] = inputSignal
			doCommands(testMemory)
			inputSignal = testMemory[9]
		}
		if inputSignal > maxThrust {
			maxThrust = inputSignal
			bestSetting = permutation
		}
	}
	fmt.Println(bestSetting)
	fmt.Println(maxThrust)
}
