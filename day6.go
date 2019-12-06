package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func parseOrbit(value string) []string {
	return strings.Split(value, ")")
}

func computePath(start string, orbits map[string]string) map[string]bool {
	var path = make(map[string]bool)
	current := start
	for {
		if current == "COM" {
			return path
		}
		path[current] = true
		current = orbits[current]
	}
}

func shortestPath(path1 map[string]bool, path2 map[string]bool) int {
	total := 0
	for star, _ := range path1 {
		if _, ok := path2[star]; !ok {
			total++
		}
	}
	for star, _ := range path2 {
		if _, ok := path1[star]; !ok {
			total++
		}
	}
	return total
}

func main() {
	argsWithoutProg := os.Args[1:]
	file, err := os.Open(argsWithoutProg[0])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	orbits := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data := scanner.Text()
		values := parseOrbit(data)
		orbits[values[1]] = values[0]
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	myPath := computePath(orbits["YOU"], orbits)
	santaPath := computePath(orbits["SAN"], orbits)

	fmt.Println(shortestPath(myPath, santaPath))
}
