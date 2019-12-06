package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func parseOrbit(value string) (center string, orbiter string) {
	i := strings.Index(value, ")")
	return value[0:i], value[i+1:]
}

func walkBackOrbit(start string, orbits map[string]string) int {
	if start == "COM" {
		return 0
	} else {
		return 1 + walkBackOrbit(orbits[start], orbits)
	}
}

func computeChecksum(orbits map[string]string) int {
	total := 0
	for k := range orbits {
		total += walkBackOrbit(k, orbits)
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
		value, key := parseOrbit(data)
		orbits[key] = value
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(computeChecksum(orbits))

}
