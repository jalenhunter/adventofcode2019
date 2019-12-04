package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func simpleConvert(str string) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		log.Fatal(err)
	}
	return val
}

func do1Command(commands []string, startPosition int) {
	parameter1 := simpleConvert(commands[startPosition+1])
	parameter2 := simpleConvert(commands[startPosition+2])
	sum := simpleConvert(commands[parameter1]) + simpleConvert(commands[parameter2])
	parameter3 := simpleConvert(commands[startPosition+3])
	commands[parameter3] = strconv.Itoa(sum)
}

func do2Command(commands []string, startPosition int) {
	parameter1 := simpleConvert(commands[startPosition+1])
	parameter2 := simpleConvert(commands[startPosition+2])
	prod := simpleConvert(commands[parameter1]) * simpleConvert(commands[parameter2])
	parameter3 := simpleConvert(commands[startPosition+3])
	commands[parameter3] = strconv.Itoa(prod)
}

func doCommand(commands []string, instructionStart int) (bool, int) {
	//switch on opcode
	switch commands[instructionStart] {
	case "1":
		do1Command(commands, instructionStart)
		return true, 4
	case "2":
		do2Command(commands, instructionStart)
		return true, 4
	case "99":
		return false, 1
	}
	return false, 0
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
	memory := strings.Split(data, ",")
	// overwrite per instructions
	for noun := 1; noun <= 99; noun++ {
		for verb := 1; verb <= 99; verb++ {
			testMemory := append(memory[:0:0], memory...)
			testMemory[1] = strconv.Itoa(noun)
			testMemory[2] = strconv.Itoa(verb)
			instructionPointer := 0
			for {
				cont, step := doCommand(testMemory, instructionPointer)
				if cont {
					instructionPointer += step
				} else {
					break
				}
			}
			if simpleConvert(testMemory[0]) == 19690720 {
				fmt.Println(100*noun + verb)
				return
			}
		}
	}
}
