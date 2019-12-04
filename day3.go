package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	vertical   = iota // c0 == 0
	horizontal = iota // c1 == 1
)

type Point struct {
	x int32
	y int32
}

type Segment struct {
	start       Point
	end         Point
	orientation int
}

func findIntersectionPoint(seg1 Segment, seg2 Segment) (bool, *Point) {
	if seg1.start.x >= seg2.start.x && seg1.start.x <= seg2.end.x {
		if seg2.start.y >= seg1.start.y && seg2.start.y <= seg1.end.y {
			fmt.Println(seg1, seg2)
			return true, &Point{x: seg1.start.x, y: seg2.start.y}
		}
	}
	if seg2.start.x >= seg1.start.x && seg2.start.x <= seg1.end.x {
		if seg1.start.y >= seg2.start.y && seg1.start.y <= seg2.end.y {
			fmt.Println(seg1, seg2)
			return true, &Point{x: seg2.start.x, y: seg1.start.y}
		}
	}
	return false, nil
}

//func Min(s1 int32, e1 int32, s2 int32, e2 int32) int32 {
//
//}

func Abs(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}

func manhattanLength(inter Point) int32 {
	return Abs(inter.x) + Abs(inter.y)
}

func createLineSegments(directions []string) []Segment {
	var segments []Segment
	current := Point{x: 0, y: 0}
	for _, instruction := range directions {
		distance, err := strconv.Atoi(instruction[1:])
		if err != nil {
			log.Fatal(err)
		}
		switch instruction[0:1] {
		case "R":
			newSegment := Segment{orientation: horizontal, start: current, end: Point{x: current.x + int32(distance), y: current.y}}
			segments = append(segments, newSegment)
			current = newSegment.end
			break
		case "D":
			newSegment := Segment{orientation: vertical, start: Point{x: current.x, y: current.y - int32(distance)}, end: current}
			segments = append(segments, newSegment)
			current = newSegment.start
			break
		case "U":
			newSegment := Segment{orientation: vertical, start: current, end: Point{x: current.x, y: current.y + int32(distance)}}
			segments = append(segments, newSegment)
			current = newSegment.end
			break
		case "L":
			newSegment := Segment{orientation: horizontal, start: Point{x: current.x - int32(distance), y: current.y}, end: current}
			segments = append(segments, newSegment)
			current = newSegment.start
			break
		default:
			fmt.Println("Unknown instruction: ", instruction)
		}
	}
	return segments
}

func main() {
	argsWithoutProg := os.Args[1:]
	file, err := os.Open(argsWithoutProg[0])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var wires [][]Segment
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data := scanner.Text()
		wires = append(wires, createLineSegments(strings.Split(data, ",")))
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	var shortestDistance int32 = math.MaxInt32
	// iterate over wires checking each segment against all the others
	for _, segment1 := range wires[0] {
		for _, segment2 := range wires[1] {
			intersected, intersection := findIntersectionPoint(segment1, segment2)
			if intersected && intersection.x != 0 && intersection.y != 0 {
				distance := manhattanLength(*intersection)
				if distance < shortestDistance {
					shortestDistance = distance
				}
			}
		}
	}
	if shortestDistance == math.MaxInt32 {
		fmt.Println("No shortest distance found?")
	} else {
		fmt.Printf("Shortest Distance: %d\n", shortestDistance)
	}
}
