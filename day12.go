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

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
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
	for i := 0; i < 1000; i++ {
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
		//fmt.Println(moons)
	}
	//compute
	total := 0
	for _, moon := range moons {
		total += (Abs(moon.x) + Abs(moon.y) + Abs(moon.z)) * (Abs(moon.vX) + Abs(moon.vY) + Abs(moon.vZ))
		//fmt.Println(total)
	}
	fmt.Println(total)

}
