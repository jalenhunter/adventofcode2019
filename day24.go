package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
)

type Level struct {
	location int
	data     string
}

func NewLevel(i int) Level {
	return Level{i, "0000000000000000000000000"}
}

const size = 5

func offset(x, y int) int {
	if x < 0 || y < 0 {
		return -1
	}
	if x >= size || y >= size {
		return -1
	}
	return x + y*size
}

func middle(x, y int) bool {
	return (x == 2 && y == 1) ||
		(x == 2 && y == 3) ||
		(x == 1 && y == 2) ||
		(x == 3 && y == 2)
}

func edge(x, y int) bool {
	return x == 0 || x == size-1 || y == 0 || y == size-1
}

func bug(str string, offset int) bool {
	if offset < 0 {
		return false
	}
	return str[offset] == '1'
}

func middleBugCount(level string, x, y int) int {
	count := 0
	if x == 1 && y == 2 { //left
		for y1 := 0; y1 < 5; y1++ {
			if bug(level, offset(0, y1)) {
				count++
			}
		}
	}
	if x == 3 && y == 2 { //right
		for y1 := 0; y1 < 5; y1++ {
			if bug(level, offset(size-1, y1)) {
				count++
			}
		}
	}
	if x == 2 && y == 1 { //up
		for x1 := 0; x1 < 5; x1++ {
			if bug(level, offset(x1, 0)) {
				count++
			}
		}
	}
	if x == 2 && y == 3 { //down
		for x1 := 0; x1 < 5; x1++ {
			if bug(level, offset(x1, size-1)) {
				count++
			}
		}
	}
	return count
}

func edgeBugCount(level string, x, y int) int {
	count := 0
	if x == 0 && bug(level, offset(1, 2)) {
		count++
	}
	if x == size-1 && bug(level, offset(3, 2)) {
		count++
	}
	if y == 0 && bug(level, offset(2, 1)) {
		count++
	}
	if y == size-1 && bug(level, offset(2, 3)) {
		count++
	}
	return count
}

func findLevel(levels []Level, location int) Level {
	for _, level := range levels {
		if level.location == location {
			return level
		}
	}
	return NewLevel(location)
}

func adjacentBugs(level Level, x, y int, levels []Level) int {
	count := 0
	if left := offset(x-1, y); left >= 0 && bug(level.data, left) {
		count++
	}
	if right := offset(x+1, y); right >= 0 && bug(level.data, right) {
		count++
	}
	if top := offset(x, y-1); top >= 0 && bug(level.data, top) {
		count++
	}
	if bottom := offset(x, y+1); bottom >= 0 && bug(level.data, bottom) {
		count++
	}
	if edge(x, y) {
		outLevel := findLevel(levels, level.location-1)
		count += edgeBugCount(outLevel.data, x, y)
	} else if middle(x, y) {
		inLevel := findLevel(levels, level.location+1)
		count += middleBugCount(inLevel.data, x, y)

	}
	return count
}

func hasEdgeBugs(level Level) bool {
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if middle(x, y) && bug(level.data, offset(x, y)) {
				return true
			}
		}
	}
	return false
}

func hasMiddleBugs(level Level) bool {
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if edge(x, y) && bug(level.data, offset(x, y)) {
				return true
			}
		}
	}
	return false
}

func do1Minute(levels []Level) []Level {
	next := []Level{}
	for _, level := range levels { //do all current levels
		newLevel := ""
		for y := 0; y < size; y++ {
			for x := 0; x < size; x++ {
				adj := adjacentBugs(level, x, y, levels)
				if x == 2 && y == 2 {
					newLevel += "0"
				} else if bug(level.data, offset(x, y)) {
					if adj != 1 {
						newLevel += "0"
					} else {
						newLevel += "1"
					}
				} else {
					if adj == 1 || adj == 2 {
						newLevel += "1"
					} else {
						newLevel += string(level.data[offset(x, y)])
					}
				}
			}
		}
		next = append(next, Level{level.location, newLevel})
	}
	return next
}

func printBugs(level Level) {
	fmt.Println("Level: ", level.location)
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if x == 2 && y == 2 {
				fmt.Print("?")
			} else if bug(level.data, offset(x, y)) {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func countSingleLevel(level string) int {
	count := 0
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if bug(level, offset(x, y)) {
				count++
			}
		}
	}
	return count
}

func countBugs(levels []Level) int {
	count := 0
	for _, level := range levels {
		count += countSingleLevel(level.data)
	}
	return count
}

func main() {
	argsWithoutProg := os.Args[1:]
	file, err := os.Open(argsWithoutProg[0])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var str = ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		for _, c := range line {
			if c == '.' {
				str += "0"
			} else {
				str += "1"
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	levels := []Level{Level{0, str}}
	currentIn, currentOut := 0, 0
	for j := 0; j < 200; j++ {
		if hasEdgeBugs(findLevel(levels, currentOut)) {
			currentOut--
			levels = append(levels, NewLevel(currentOut))
		}
		if hasMiddleBugs(findLevel(levels, currentIn)) {
			currentIn++
			levels = append(levels, NewLevel(currentIn))
		}
		levels = do1Minute(levels)
	}
	sort.Slice(levels, func(i, j int) bool {
		return levels[i].location < levels[j].location
	})
	for _, level := range levels {
		printBugs(level)
	}
	fmt.Println("Bug count = ", countBugs(levels))
}
