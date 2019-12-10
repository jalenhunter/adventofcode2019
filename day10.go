package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
)

type Point struct {
	x     int
	y     int
	angle float64
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

func getVisible(base Point, asteroids []Point, myIndex int) []Point {
	var visible []Point
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
			visible = append(visible, asteroid)
		}
	}
	return visible
}

func calculateAngle(base Point, asteroid Point) float64 {
	yCoord := float64(base.y - asteroid.y)
	xCoord := float64(base.x - asteroid.x)
	angleR := math.Atan2(yCoord, xCoord)
	angleD := angleR * 180 / math.Pi

	angleD -= 90 //rotate to top Y axis
	if angleD < 0 {
		angleD += 360
	}
	return angleD
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

	var mostVisibleAsteroids []Point
	maxVisible := 0
	var bestBase Point
	for index, base := range asteroids {
		visibleAsteroids := getVisible(base, asteroids, index)
		if len(visibleAsteroids) > maxVisible {
			maxVisible = len(visibleAsteroids)
			bestBase = base
			mostVisibleAsteroids = visibleAsteroids
		}
	}
	fmt.Println(maxVisible)
	fmt.Println(bestBase)

	var angledAsteroids []Point
	for _, asteroid := range mostVisibleAsteroids {
		angledAsteroids = append(angledAsteroids, Point{
			x:     asteroid.x,
			y:     asteroid.y,
			angle: calculateAngle(bestBase, asteroid),
		})
	}

	sort.Slice(angledAsteroids, func(a, b int) bool {
		return angledAsteroids[a].angle < angledAsteroids[b].angle
	})
	fmt.Println(100*angledAsteroids[199].x + angledAsteroids[199].y)

}
