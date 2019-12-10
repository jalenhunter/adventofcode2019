package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Point struct {
	x int
	y int
}

type Line struct {
	deltaX int
	deltaY int
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func gcdTwoNumbers(x int, y int) int {
	x = Abs(x)
	y = Abs(y)
	for y > 0 {
		var temp = y
		y = x % y
		x = temp
	}
	return x
}

func computeLine(start Point, end Point) (Line, int) {
	distance := Abs(end.x-start.x) + Abs(end.y-start.y)
	deltaX := end.x - start.x
	deltaY := end.y - start.y
	gcd := gcdTwoNumbers(deltaX, deltaY)
	return Line{deltaX / gcd, deltaY / gcd}, distance
}

func inline(asteroid Point, line Line, base Point, check Point) (bool, int) {
	line2, distance := computeLine(base, asteroid)
	return line2.deltaY == line.deltaY && line2.deltaX == line.deltaX, distance
}

func computeVisible(base Point, asteroids []Point, myIndex int) int {
	count := 0
	//fmt.Println("Checking asteroid at :", base)
	for index, asteroid := range asteroids {
		if index == myIndex {
			continue
		}
		line, distance := computeLine(base, asteroid)
		isVisible := true
		for index2, asteroid2 := range asteroids {
			if index2 == index || index2 == myIndex {
				continue
			}
			inline, distance2 := inline(asteroid2, line, base, asteroid)
			if inline && distance2 < distance {
				isVisible = false
				break
			}
		}
		if isVisible {
			count++
		}
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
	var asteroids []Point
	scanner := bufio.NewScanner(file)
	y := 0
	for scanner.Scan() {
		line := scanner.Text()
		for index, c := range line {
			if c == '#' {
				asteroids = append(asteroids, Point{x: index, y: y})
			}
		}
		y++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	maxVisible := 0
	var bestBase Point
	for index, base := range asteroids {
		count := computeVisible(base, asteroids, index)
		if count > maxVisible {
			maxVisible = count
			bestBase = base
		}
	}
	fmt.Println(maxVisible)
	fmt.Println(bestBase)
}
