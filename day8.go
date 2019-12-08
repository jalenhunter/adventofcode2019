package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

type Layer [6][25]int

func printLayer(layer Layer) {
	for j := 0; j < 6; j++ {
		for i := 0; i < 25; i++ {
			if layer[j][i] == 1 {
				fmt.Print("\xE2\x96\xA0")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
	fmt.Println()
	fmt.Println()
}

func makeImage(layer Layer) {
	rectangle := image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: 25, Y: 6}}
	passwordImage := image.NewRGBA(rectangle)
	for j := 0; j < 6; j++ {
		for i := 0; i < 25; i++ {
			if layer[j][i] == 1 {
				passwordImage.Set(i, j, color.White)
			} else {
				passwordImage.Set(i, j, color.Black)
			}
		}
	}
	file, err := os.Create("password.png")
	if err != nil {
		log.Fatal(err)
	}
	err = png.Encode(file, passwordImage)
	if err != nil {
		log.Fatal(err)
	}
}

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
	pixel := 0
	layers := len(data) / (25 * 6)
	layersData := make([]Layer, layers)

	for layer := 0; layer < layers; layer++ {
		var layerData Layer
		for j := 0; j < 6; j++ {
			for i := 0; i < 25; i++ {
				char := data[pixel]
				switch char {
				case '0':
					layerData[j][i] = 0
					break
				case '1':
					layerData[j][i] = 1
					break
				case '2':
					layerData[j][i] = 2
					break
				default:
					log.Fatal(char)
				}
				pixel++
			}
		}
		layersData[layer] = layerData
	}
	var final [6][25]int
	for i := 0; i < 25; i++ {
		for j := 0; j < 6; j++ {
			final[j][i] = 2
		}
	}

	for layer := 0; layer < layers; layer++ {
		for j := 0; j < 6; j++ {
			for i := 0; i < 25; i++ {
				if layersData[layer][j][i] != 2 && final[j][i] == 2 {
					final[j][i] = layersData[layer][j][i]
				}
			}
		}
	}

	printLayer(final)
	makeImage(final)
}
