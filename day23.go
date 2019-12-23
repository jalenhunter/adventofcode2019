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
type result func() Packet

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
	var queues = make([][]Packet, 50)
	for j := 0; j < 50; j++ {
		queues[j] = []Packet{}
	}
	memory := toInt(strings.Split(data, ","))
	var answers = []result{}
	memory = append(memory, make([]int64, 1000)...)
	for i := 0; i < 50; i++ {
		testMemory := append(memory[:0:0], memory...)
		in, out, answer := handleIO(i, &queues)
		answers = append(answers, answer)
		go doCommands(testMemory, in, out)
	}
	var lastY int64 = -1
	var answer Packet
	for {
		allEmpty := true
		for k := 0; k < 50; k++ {
			if len(queues[k]) != 0 {
				fmt.Println("Not Empty:", k)
				allEmpty = false
				break
			}
			temp := answers[k]()
			if temp.y != -1 {
				answer = temp
			}
		}
		if allEmpty {
			if lastY == answer.y {
				fmt.Println(lastY)
				return
			}
			lastY = answer.y
			//fmt.Println("Adding to queue 0, ", answer)
			queues[0] = append(queues[0], answer)
		}
	}
}

type Packet struct {
	x, y int64
}

func handleIO(address int, queues *[][]Packet) (input, output, result) {
	sentAddress := false
	sendX := true
	in := func() int {
		if !sentAddress {
			//fmt.Println("Writing Address:", address)
			sentAddress = true
			return address
		}
		if len((*queues)[address]) == 0 {
			fmt.Println("Empty Address:", address)
			return -1
		}
		packet := (*queues)[address][0]
		if sendX {
			//fmt.Println("Writing: X from ", packet, " to ", address)
			sendX = false
			return int(packet.x)
		} else {
			//fmt.Println("Writing: Y from ", packet, " to ", address)
			sendX = true
			(*queues)[address] = (*queues)[address][1:]
			return int(packet.y)
		}
	}
	answer := Packet{-1, -1}
	gotX := false
	toAddress := -1
	var x int64 = 0
	var y int64 = 0
	out := func(r int64) {
		if toAddress == -1 {
			toAddress = int(r)
		} else if !gotX {
			gotX = true
			x = r
		} else {
			y = r
			if toAddress == 255 {
				answer = Packet{x, y}
			} else {
				(*queues)[toAddress] = append((*queues)[toAddress], Packet{x, y})
			}
			toAddress = -1
			gotX = false
		}
	}

	result := func() Packet {
		return answer
	}
	return in, out, result
}
