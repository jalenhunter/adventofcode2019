package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

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
	fmt.Println(len(data))
	pixel := 0
	layers := len(data) / (25 * 6)
	leastZero := 10000000
	checksum := 0
	for layer := 0; layer < layers; layer++ {
		zeros := 0
		ones := 0
		twos := 0
		for i := 0; i < 25; i++ {
			for j := 0; j < 6; j++ {
				char := data[pixel]
				switch char {
				case '0':
					zeros++
					break
				case '1':
					ones++
					break
				case '2':
					twos++
					break
				default:
					log.Fatal(char)
				}
				pixel++
			}
		}
		if zeros < leastZero {
			leastZero = zeros
			checksum = ones * twos
		}
	}
	fmt.Println(checksum)

}
