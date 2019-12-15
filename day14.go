package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Chem struct {
	name   string
	amount int64
}

type Reaction struct {
	inputs []Chem
	output Chem
}

func makeChem(str string) Chem {
	parser := regexp.MustCompile(`(\d+)\s(\w+)`)
	if parser.MatchString(str) {
		data := parser.FindStringSubmatch(str)
		amount, _ := strconv.ParseInt(data[1], 10, 64)
		return Chem{data[2], amount}
	} else {
		log.Fatal("Problem parsing chemical", str)
	}
	return Chem{}
}

var surplus = make(map[string]int64)

func getChemicals(qty int64, chem string, reactions []Reaction) map[string]int64 {
	for _, reaction := range reactions {
		if reaction.output.name == chem {
			chemQty := reaction.output.amount
			mult := int64(math.Ceil(float64(qty) / float64(chemQty)))
			adjusted := make(map[string]int64)
			for _, input := range reaction.inputs {
				adjusted[input.name] = input.amount*mult - surplus[input.name]
			}
			final := make(map[string]int64)
			for n, q := range adjusted {
				surplus[n] = 0
				if q < 0 {
					surplus[n] = -q
				} else {
					final[n] = q
				}
			}
			if chemQty*mult > qty {
				surplus[chem] += chemQty*mult - qty
			}
			return final
		} else {
			continue
		}
	}
	log.Fatal("Missing recipe: ", chem)
	return make(map[string]int64)
}

func firstKey(needs map[string]int64) string {
	for key, _ := range needs {
		return key
	}
	return ""
}

func getOre(fuel int64, reactions []Reaction) int64 {
	var ore int64 = 0
	needs := make(map[string]int64)
	needs["FUEL"] = fuel
	for {
		key := firstKey(needs)
		if len(key) == 0 {
			break
		}
		qty := needs[key]
		delete(needs, key)
		//fmt.Println(key, qty)
		chemicals := getChemicals(qty, key, reactions)
		//fmt.Println(chemicals)
		for name, value := range chemicals {
			if name == "ORE" {
				ore += value
			} else {
				needs[name] += value
			}
		}
	}
	return ore
}

func main() {
	splitter := regexp.MustCompile(`[,|=>]`)
	argsWithoutProg := os.Args[1:]
	file, err := os.Open(argsWithoutProg[0])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var reactions []Reaction
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var inputs []Chem
		var output Chem
		line := scanner.Text()
		chems := splitter.Split(line, -1)
		in := true
		for _, chem := range chems {
			chem = strings.TrimSpace(chem)
			if len(chem) > 0 {
				if in {
					inputs = append(inputs, makeChem(chem))
				} else {
					output = makeChem(chem)
				}
			} else {
				in = false
			}
		}
		reactions = append(reactions, Reaction{inputs, output})
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	simpleOre := getOre(1, reactions)
	fmt.Println(simpleOre)

	var myOre int64 = 1000000000000
	surplus = make(map[string]int64)
	var max int64 = 3061523 - 2
	for {
		max += 1
		surplus = make(map[string]int64)
		testOre := getOre(max, reactions)
		if testOre > myOre {
			fmt.Println("Greater: ", max)
			break
		} else {
			fmt.Println(max, testOre)
		}
	}
}
