package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func fuelNeedForModule(weight int) int {
	return (weight / 3) - 2
}

func totalFueldNeededForModule(weight int) int {
	total := 0
	fuel := fuelNeedForModule(weight)
	for fuel > 0 {
		total += fuel
		fuel = fuelNeedForModule(fuel)
	}
	return total
}

func main() {
	argsWithoutProg := os.Args[1:]
	total := 0
	file, err := os.Open(argsWithoutProg[0])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		value, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		total += totalFueldNeededForModule(value)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total  %d\n ", total)
}
