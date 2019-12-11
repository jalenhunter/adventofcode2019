package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	x     int
	y     int
	color int
}

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

func do4Command(commands []int64, startPosition int, modes []int, position int, output chan int64) int {
	parameter1 := commands[startPosition+1]
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

func newDirection(current string, turn int64) string {
	switch current {
	case "U":
		if turn == 0 {
			return "L"
		} else {
			return "R"
		}
		break
	case "L":
		if turn == 0 {
			return "D"
		} else {
			return "U"
		}
		break
	case "R":
		if turn == 0 {
			return "U"
		} else {
			return "D"
		}
		break
	case "D":
		if turn == 0 {
			return "R"
		} else {
			return "L"
		}
		break
	default:
		log.Fatal("Unknown Direction", current, turn)
	}
	fmt.Println("Odd Command", current, turn)
	return "U"
}

func pointExists(points []*Point, x int, y int) *Point {
	for _, point := range points {
		if point.x == x && point.y == y {
			return point
		}
	}
	return nil
}

func moveOne(points []*Point, direction string, position *Point) (*Point, []*Point) {
	x := position.x
	y := position.y

	switch direction {
	case "U":
		y -= 1
		break
	case "L":
		x -= 1
		break
	case "R":
		x += 1
		break
	case "D":
		y += 1
		break
	}
	newPoint := pointExists(points, x, y)
	if newPoint != nil {
		return newPoint, points
	} else {
		addedPoint := Point{x, y, 0}
		points = append(points, &addedPoint)
		return &addedPoint, points
	}
}

func main() {
	argsWithoutProg := os.Args[1:]
	file, err := os.Open(argsWithoutProg[0])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	input := make(chan int)
	state := make(chan int)
	output := make(chan int64)

	var points []*Point
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
	go doCommands(memory, input, output, state)
	rectangle := image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 100, Y: 100}}
	registrationImage := image.NewRGBA(rectangle)
	for j := 0; j < 100; j++ {
		for i := 0; i < 100; i++ {
			registrationImage.Set(i, j, color.Black)
		}
	}
	position := &Point{20, 20, 1}
	points = append(points, position)
	currentDirection := "U"
	for {
		progState, ok := <-state
		if !ok {
			break
		}
		if progState == 3 {
			input <- position.color
			//read number 1 is color
			pixelColor := <-output
			if pixelColor == 1 {
				registrationImage.Set(position.x, position.y, color.White)
			} else {
				registrationImage.Set(position.x, position.y, color.Black)
			}
			//fmt.Println(color)
			position.color = int(pixelColor)
			// read number 2 is turn
			turn := <-output
			currentDirection = newDirection(currentDirection, turn)
			position, points = moveOne(points, currentDirection, position)

		}
	}
	close(input)
	close(output)
	outFile, err := os.Create("registration.png")
	if err != nil {
		log.Fatal(err)
	}
	err = png.Encode(outFile, registrationImage)
	if err != nil {
		log.Fatal(err)
	}
}
