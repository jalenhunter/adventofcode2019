package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func convertScan(scan string) []int {
	data := make([]int, len(scan))
	for index, c := range scan {
		data[index] = toInt(string(c))
	}
	return data
}

func combine(start int, end int, data []int) int {
	var ret = 0
	for i := start; i < end; i++ {
		ret = ret*10 + data[i]
	}
	return ret
}

func main() {
	argsWithoutProg := os.Args[1:]
	file, err := os.Open(argsWithoutProg[0])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var scanned string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		scanned = scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	data := convertScan(scanned)
	start := combine(0, 7, data)
	newData := make([]int, len(data)*10000)
	for i := 0; i < 10000; i++ {
		for index, _ := range data {
			newData[(i*len(data))+index] = data[index]
		}
	}
	for p := 0; p < 100; p++ {
		newData = doPhase2(&newData)
	}
	fmt.Println(combine(start, start+8, newData))
}

func doPhase2(data *[]int) []int {
	max := len(*data)
	ret := make([]int, max)
	var ans int
	for row := 0; row < max/2; row++ {
		ans += (*data)[max-(row+1)]
		ret[max-(row+1)] = Abs(ans) % 10
	}
	return ret
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func toInt(c string) int {
	x, err := strconv.Atoi(c)
	if err != nil {
		log.Fatal(err)
	}
	return x
}

func doPhase1(data *[]int) []int {
	base := []int8{0, 1, 0, -1}
	max := len(*data)
	ret := make([]int, max)
	for row := 0; row < max; row++ {
		var sum = 0
		for col := 0; col < max; col++ {
			bp := ((col + 1) / (row + 1)) % 4
			multi := base[bp]
			if multi == -1 {
				sum -= int((*data)[col])
			} else if multi == 1 {
				sum += int((*data)[col])
			}
		}
		sum %= 10
		ret[row] = Abs(sum)
	}
	return ret
}
