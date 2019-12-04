package main

import (
	"fmt"
	"log"
	"strconv"
)

func hasAdjacentDigits(guess int) bool {
	str := strconv.Itoa(guess)
	var counter [10]int
	for _, c := range str {
		val, err := strconv.Atoi(string(c))
		if err != nil {
			log.Fatal(err)
		}
		counter[val]++
	}
	for _, value := range counter {
		if value == 2 {
			return true
		}
	}
	return false
}

func digitsIncrease(guess int) bool {
	str := strconv.Itoa(guess)
	prev := 0
	for _, c := range str {
		val, err := strconv.Atoi(string(c))
		if err != nil {
			log.Fatal(err)
		}
		if val >= prev {
			prev = val
		} else {
			return false
		}
	}
	return true
}

func main() {
	count := 0
	start, end := 402328, 864247
	for guess := start; guess <= end; guess++ {
		if digitsIncrease(guess) {
			if hasAdjacentDigits(guess) {
				fmt.Println(guess)
				count++
			}
		}
	}
	fmt.Println("Possible Passwords: ", count)
}
