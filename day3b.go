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

type Point struct {
	x int32
	y int32
}

type Segment struct {
	start    Point
	end      Point
	distance int32
	origin   Point
}

func findIntersectionPoint(seg1 Segment, seg2 Segment) (bool, *Point) {
	if seg1.start.x >= seg2.start.x && seg1.start.x <= seg2.end.x {
		if seg2.start.y >= seg1.start.y && seg2.start.y <= seg1.end.y {
			return true, &Point{x: seg1.start.x, y: seg2.start.y}
		}
	}
	if seg2.start.x >= seg1.start.x && seg2.start.x <= seg1.end.x {
		if seg1.start.y >= seg2.start.y && seg1.start.y <= seg2.end.y {
			return true, &Point{x: seg2.start.x, y: seg1.start.y}
		}
	}
	return false, nil
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
			newSegment := Segment{origin: current, distance: int32(distance), start: current, end: Point{x: current.x + int32(distance), y: current.y}}
			segments = append(segments, newSegment)
			current = newSegment.end
			break
		case "D":
			newSegment := Segment{origin: current, distance: int32(distance), start: Point{x: current.x, y: current.y - int32(distance)}, end: current}
			segments = append(segments, newSegment)
			current = newSegment.start
			break
		case "U":
			newSegment := Segment{origin: current, distance: int32(distance), start: current, end: Point{x: current.x, y: current.y + int32(distance)}}
			segments = append(segments, newSegment)
			current = newSegment.end
			break
		case "L":
			newSegment := Segment{origin: current, distance: int32(distance), start: Point{x: current.x - int32(distance), y: current.y}, end: current}
			segments = append(segments, newSegment)
			current = newSegment.start
			break
		default:
			fmt.Println("Unknown instruction: ", instruction)
		}
	}
	return segments
}

func partialDistance(p Point, s Segment) int32 {
	if p.x == s.start.x {
		return p.y - s.origin.y
	} else {
		return p.x - s.origin.x
	}
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
	var segDist1 int32 = 0
	// iterate over wires checking each segment against all the others
	for _, segment1 := range wires[0] {
		var segDist2 int32 = 0
		for _, segment2 := range wires[1] {
			intersected, intersection := findIntersectionPoint(segment1, segment2)
			if intersected && intersection.x != 0 && intersection.y != 0 {
				distance1 := partialDistance(*intersection, segment1)
				distance2 := partialDistance(*intersection, segment2)
				totalDistance := distance1 + distance2 + segDist1 + segDist2
				if totalDistance < shortestDistance {
					shortestDistance = totalDistance
				}
			}
			segDist2 += segment2.distance
		}
		segDist1 += segment1.distance
	}
	if shortestDistance == math.MaxInt32 {
		fmt.Println("No shortest distance found?")
	} else {
		fmt.Printf("Shortest Distance: %d\n", shortestDistance)
	}
}
