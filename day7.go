package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
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

func do3Command(index int, commands []int, startPosition int, input chan int) int {
	parameter1 := commands[startPosition+1]
	commands[parameter1] = <-input
	return startPosition + 2
}

func do4Command(index int, commands []int, startPosition int, modes []int, output chan int) int {
	parameter1 := commands[startPosition+1]
	output <- handleMode(commands, parameter1, modes[0])
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

func doCommands(index int, commands []int, input chan int, output chan int, wg *sync.WaitGroup) {
	instructionPointer := 0
	for {
		opCode, modes := parseCommand(commands[instructionPointer])
		switch opCode {
		case 1:
			instructionPointer = do1Command(commands, instructionPointer, modes)
		case 2:
			instructionPointer = do2Command(commands, instructionPointer, modes)
		case 3:
			instructionPointer = do3Command(index, commands, instructionPointer, input)
		case 4:
			instructionPointer = do4Command(index, commands, instructionPointer, modes, output)
		case 5:
			instructionPointer = do5Command(commands, instructionPointer, modes)
		case 6:
			instructionPointer = do6Command(commands, instructionPointer, modes)
		case 7:
			instructionPointer = do7Command(commands, instructionPointer, modes)
		case 8:
			instructionPointer = do8Command(commands, instructionPointer, modes)
		case 99:
			if wg != nil {
				wg.Done()
			}
			close(output)
			//fmt.Printf("Closing Amp %d\n", index)
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
	Perm([]rune("56789"), func(a []rune) {
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
		var chans [5]chan int
		for i := range chans {
			chans[i] = make(chan int)
		}
		var wg sync.WaitGroup
		for index, phase := range permutation {
			testMemory := append(memory[:0:0], memory...)
			v, err := strconv.Atoi(string(phase))
			if err != nil {
				log.Fatal(err)
			}
			if index == 4 {
				go doCommands(index, testMemory, chans[index], chans[0], nil)
			} else {
				wg.Add(1)
				go doCommands(index, testMemory, chans[index], chans[index+1], &wg)
			}
			chans[index] <- v
		}
		chans[0] <- 0
		wg.Wait()
		value := <-chans[0]
		//fmt.Println("Got Answer:", value)
		if value > maxThrust {
			maxThrust = value
			bestSetting = permutation
		}
	}
	fmt.Println(bestSetting)
	fmt.Println(maxThrust)
}
