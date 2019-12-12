package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Moon struct {
	x, y, z    int
	vX, vY, vZ int
	name       string
}

func getValue(input string) int {
	val, err := strconv.Atoi(strings.Split(strings.Trim(input, " "), "=")[1])
	if err != nil {
		log.Fatal(err)
	}
	return val
}

func deltaVel(a int, b int) int {
	if a < b {
		return 1
	} else if a > b {
		return -1
	}
	return 0
}

func computeVel(moon Moon, moons []Moon) (int, int, int) {
	dx, dy, dz := 0, 0, 0
	for _, other := range moons {
		if moon.name == other.name {
			continue
		}
		dx += deltaVel(moon.x, other.x)
		dy += deltaVel(moon.y, other.y)
		dz += deltaVel(moon.z, other.z)
	}
	//fmt.Println(dx, dy, dz)
	return dx, dy, dz
}

func hashX(moons []Moon) string {
	str := ""
	for index, _ := range moons {
		str += strconv.Itoa(moons[index].x) + "," + strconv.Itoa(moons[index].vX)
		if index != 3 {
			str += ","
		}
	}
	return str
}
func hashY(moons []Moon) string {
	str := ""
	for index, _ := range moons {
		str += strconv.Itoa(moons[index].y) + "," + strconv.Itoa(moons[index].vY)
		if index != 3 {
			str += ","
		}
	}
	return str
}
func hashZ(moons []Moon) string {
	str := ""
	for index, _ := range moons {
		str += strconv.Itoa(moons[index].z) + "," + strconv.Itoa(moons[index].vZ)
		if index != 3 {
			str += ","
		}
	}
	return str
}

func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func LCM(a, b int, integers ...int) int {
	result := a * b / GCD(a, b)

	for i := 0; i < len(integers); i++ {
		result = LCM(result, integers[i])
	}

	return result
}

func main() {
	argsWithoutProg := os.Args[1:]
	file, err := os.Open(argsWithoutProg[0])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	names := [4]string{"Io", "Europa", "Ganymede", "Callisto"}
	index := 0
	scanner := bufio.NewScanner(file)
	var moons []Moon
	for scanner.Scan() {
		line := scanner.Text()
		data := strings.Split(line[1:len(line)-1], ",")
		moons = append(moons, Moon{getValue(data[0]), getValue(data[1]), getValue(data[2]), 0, 0, 0, names[index]})
		index++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	var steps int = 0
	var xStates = make(map[string]bool)
	var yStates = make(map[string]bool)
	var zStates = make(map[string]bool)
	xStep, yStep, zStep := 0, 0, 0
	for {
		//do the velocity
		for index, moon := range moons {
			dx, dy, dz := computeVel(moon, moons)
			moons[index].vX += dx
			moons[index].vY += dy
			moons[index].vZ += dz
		}
		//do the movement
		for index, _ := range moons {
			moons[index].x, moons[index].y, moons[index].z = moons[index].x+moons[index].vX, moons[index].y+moons[index].vY, moons[index].z+moons[index].vZ
		}
		xStr := hashX(moons)
		if _, ok := xStates[xStr]; ok {
			if xStep == 0 {
				xStep = steps
			}
		} else {
			xStates[xStr] = true
		}
		yStr := hashY(moons)
		if _, ok := yStates[yStr]; ok {
			if yStep == 0 {
				yStep = steps
			}
		} else {
			yStates[yStr] = true
		}
		zStr := hashZ(moons)
		if _, ok := zStates[zStr]; ok {
			if zStep == 0 {
				zStep = steps
			}
		} else {
			zStates[zStr] = true
		}
		if xStep > 0 && yStep > 0 && zStep > 0 {
			break
		}
		steps++
	}
	fmt.Println(xStep, yStep, zStep)
	fmt.Println(LCM(xStep, yStep, zStep))
}
