package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const size = 100

var locations [size][size]int

func handleModeRead(commands []int, param int, mode int, position int) int {
	if mode == 0 { //reference
		return commands[param]
	} else if mode == 1 {
		return param
	} else if mode == 2 {
		return commands[position+param]
	} else {
		log.Fatal("Unknown mode:", mode)
	}
	return 0
}

func handleModeWrite(commands []int, param int, mode int, position int, result int) {
	if mode == 0 { //reference
		commands[param] = result
	} else if mode == 1 {
		log.Fatal("Illegal Write")
	} else if mode == 2 {
		commands[position+param] = result
	} else {
		log.Fatal("Unknown write mode:", mode)
	}
}

func do1Command(commands []int, startPosition int, modes []int, position int) int {
	parameter1 := commands[startPosition+1]
	parameter2 := commands[startPosition+2]
	parameter3 := commands[startPosition+3]
	sum := handleModeRead(commands, parameter1, modes[0], position) + handleModeRead(commands, parameter2, modes[1], position)
	handleModeWrite(commands, parameter3, modes[2], position, sum)
	return startPosition + 4
}

func do2Command(commands []int, startPosition int, modes []int, position int) int {
	parameter1 := commands[startPosition+1]
	parameter2 := commands[startPosition+2]
	parameter3 := commands[startPosition+3]
	prod := handleModeRead(commands, parameter1, modes[0], position) * handleModeRead(commands, parameter2, modes[1], position)
	handleModeWrite(commands, parameter3, modes[2], position, prod)
	return startPosition + 4
}

func do3Command(commands []int, startPosition int, modes []int, position int, input chan int) int {
	parameter1 := commands[startPosition+1]
	var i = <-input
	handleModeWrite(commands, parameter1, modes[0], position, int(i))
	return startPosition + 2
}

func do4Command(commands []int, startPosition int, modes []int, position int, output chan int) int {
	parameter1 := commands[startPosition+1]
	output <- handleModeRead(commands, parameter1, modes[0], position)
	return startPosition + 2
}

func do5Command(commands []int, startPosition int, modes []int, position int) int {
	parameter1 := commands[startPosition+1]
	parameter2 := commands[startPosition+2]
	if handleModeRead(commands, parameter1, modes[0], position) > 0 {
		return int(handleModeRead(commands, parameter2, modes[1], position))
	} else {
		return startPosition + 3
	}
}
func do6Command(commands []int, startPosition int, modes []int, position int) int {
	parameter1 := commands[startPosition+1]
	parameter2 := commands[startPosition+2]
	if handleModeRead(commands, parameter1, modes[0], position) == 0 {
		return int(handleModeRead(commands, parameter2, modes[1], position))
	} else {
		return startPosition + 3
	}
}
func do7Command(commands []int, startPosition int, modes []int, position int) int {
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
func do8Command(commands []int, startPosition int, modes []int, position int) int {
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

func do9Command(commands []int, startPosition int, modes []int, position int) (int, int) {
	parameter1 := commands[startPosition+1]
	//fmt.Printf("Changing Position from %d to %d\n", position, handleModeRead(commands, parameter1, modes[0], position))
	return startPosition + 2, position + int(handleModeRead(commands, parameter1, modes[0], position))
}

func parseCommand(command int) (opCode int, mode []int) {
	code := command % 100
	modes := command / 100
	paramMode := make([]int, 3)
	for i := 0; i < 3; i++ {
		paramMode[i] = int(modes % 10)
		modes = modes / 10
	}
	return int(code), paramMode
}

func doCommands(commands []int, input chan int, output chan int) {
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
			instructionPointer = do3Command(commands, instructionPointer, modes, position, input)
		case 4:
			instructionPointer = do4Command(commands, instructionPointer, modes, position, output)
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
			close(output)
			return
		}
	}
}

func toInt(input []string) []int {
	data := make([]int, len(input))
	for index, code := range input {
		n, err := strconv.Atoi(code)
		if err != nil {
			log.Fatal(err)
		}
		data[index] = n
	}
	return data
}

func getNextDirection(x int, y int) int {
	return rand.Intn(4) + 1

	//for {
	//	testDir := rand.Intn(4) + 1
	//	var status = -1
	//	switch testDir {
	//	case 1:
	//		status = locations[x][y + 1]
	//		break
	//	case 2:
	//		status = locations[x][y - 1]
	//		break
	//	case 3:
	//		status = locations[x + 1][y]
	//		break
	//	case 4:
	//		status = locations[x - 1][y]
	//		break
	//	}
	//	if status == 1 || status == -1{
	//		return testDir
	//	}
	//}
}

func printLocations(actual bool) {
	for j := 0; j < size; j++ {
		for i := 0; i < size; i++ {
			if actual {
				fmt.Print(locations[j][i])
			} else {
				if j == size/2 && i == size/2 {
					fmt.Print("S")
				} else {
					if locations[j][i] == 0 {
						fmt.Print("O")
					} else if locations[j][i] == 2 {
						fmt.Print("X")
					} else if locations[j][i] > 1 {
						fmt.Print(" ")
					} else {
						fmt.Print("U")
					}
				}
			}
		}
		fmt.Println()
	}
	fmt.Println()
	fmt.Println()
}

func getNextPosition(x int, y int, dir int) (int, int) {
	switch dir {
	case 1:
		return x, y + 1
	case 2:
		return x, y - 1
	case 3:
		return x + 1, y
	case 4:
		return x - 1, y
	default:
		fmt.Println("Unknown Direction", dir)
		return x, y
	}
}

func main() {
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			locations[i][j] = -1
		}
	}
	//printLocations(true)
	rand.Seed(time.Now().UnixNano())
	argsWithoutProg := os.Args[1:]
	file, err := os.Open(argsWithoutProg[0])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	input := make(chan int)
	output := make(chan int)

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
	memory = append(memory, make([]int, 1000)...)
	found := false
	var x, y = size / 2, size / 2
	locations[x][y] = 1
	var dir = 0
	go doCommands(memory, input, output)
	distance := 0
	for {
		dir = getNextDirection(x, y)
		//fmt.Print(dir, ":")
		input <- dir
		status, ok := <-output
		if !ok {
			fmt.Println("What happened to the program?")
			break
		}
		//fmt.Print(status, " ")
		a, b := getNextPosition(x, y, dir)
		switch status {
		case 0: //wall
			locations[a][b] = 0
			break
		case 1: //clear
			if locations[a][b] < 0 {
				distance++
				locations[a][b] = distance
			} else {
				distance = locations[a][b]
			}
			x, y = a, b
			break
		case 2: //oxygen
			fmt.Println("Found It:", x, y)
			found = true
			distance++
			locations[a][b] = 2
			x, y = a, b
			break
		default:
			fmt.Println("Unknown value from computer:", status)
			break
		}
		if found {
			break
		}
	}
	//printLocations(false)
	fmt.Println(distance)
	close(input)

}
