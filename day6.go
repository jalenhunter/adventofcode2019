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

func computePath(start string, orbits map[string]string) []string {
	var path []string
	current := start
	for {
		if current == "COM" {
			return path
		}
		path = append(path, current)
		current = orbits[current]
	}
}

func shortestPath(path1 []string, path2 []string) int {
	path1Len := 0
	for _, step1 := range path1 {
		path2Len := 0
		for _, step2 := range path2 {
			//fmt.Printf("%s ? %s", step1, step2)
			if step1 == step2 {
				return path1Len + path2Len
			}
			path2Len += 1
		}
		path1Len += 1
	}
	return -1
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
