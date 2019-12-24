package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

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

func bug(str string, offset int) bool {
	return str[offset] == '1'
}

func adjacentBugs(str string, x, y int) int {
	count := 0
	if left := offset(x-1, y); left >= 0 && bug(str, left) {
		count++
	}
	if right := offset(x+1, y); right >= 0 && bug(str, right) {
		count++
	}
	if top := offset(x, y-1); top >= 0 && bug(str, top) {
		count++
	}
	if bottom := offset(x, y+1); bottom >= 0 && bug(str, bottom) {
		count++
	}
	return count
}

func do1Minute(str string) string {
	next := ""
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			adj := adjacentBugs(str, x, y)
			if bug(str, offset(x, y)) {
				if adj != 1 {
					next += "0"
				} else {
					next += "1"
				}
			} else {
				if adj == 1 || adj == 2 {
					next += "1"
				} else {
					next += string(str[offset(x, y)])
				}
			}
		}
	}
	return next
}

func printBugs(str string) {
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if bug(str, offset(x, y)) {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func biodiversity(str string) int64 {
	var bd int64 = 0
	for i, c := range str {
		if c == '1' {
			bd += int64(math.Pow(2, float64(i)))
		}
	}
	return bd
}

func main() {
	argsWithoutProg := os.Args[1:]
	file, err := os.Open(argsWithoutProg[0])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	seen := make(map[int64]bool)
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
	fmt.Println(str)
	printBugs(str)
	for {
		i, err := strconv.ParseInt(str, 2, 64)
		if err != nil {
			log.Fatal(err)
		}
		if seen[i] {
			fmt.Println(str)
			fmt.Println(biodiversity(str))
			break
		}
		seen[i] = true
		str = do1Minute(str)

		printBugs(str)
	}
}
