package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

const (
	CUT  int = 1
	NEW  int = 2
	DEAL int = 3
)

type Command struct {
	comm  int
	count int
}

func parseCommand(command string) Command {
	cutParser := regexp.MustCompile(`^cut (.*\d+)$`)
	dealParser := regexp.MustCompile(`^deal with increment (\d+)$`)
	if cutParser.MatchString(command) {
		data := cutParser.FindStringSubmatch(command)
		count, _ := strconv.Atoi(data[1])
		return Command{CUT, count}
	}
	if "deal into new stack" == command {
		return Command{NEW, -1}
	}
	if dealParser.MatchString(command) {
		data := dealParser.FindStringSubmatch(command)
		count, _ := strconv.Atoi(data[1])
		return Command{DEAL, count}
	}
	log.Fatal("Unknown command in parser")
	return Command{}
}

func main() {
	argsWithoutProg := os.Args[1:]
	file, err := os.Open(argsWithoutProg[0])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	commands := []Command{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		command := parseCommand(scanner.Text())
		commands = append(commands, command)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	deck := make([]int, 10007)
	for i := 0; i < len(deck); i++ {
		deck[i] = i
	}
	for _, command := range commands {
		fmt.Println("Running: ", command)
		switch command.comm {
		case CUT: //
			if command.count > 0 {
				deck = append(deck[command.count:], deck[0:command.count]...)
			} else {
				deck = append(deck[len(deck)+command.count:], deck[0:len(deck)+command.count]...)
			}
			break
		case NEW: //reverse
			for i, j := 0, len(deck)-1; i < j; i, j = i+1, j-1 {
				deck[i], deck[j] = deck[j], deck[i]
			}
			break
		case DEAL: //skip many with loop around
			newDeck := make([]int, len(deck))
			for i := 0; i < len(deck); i++ {
				spot := (i * command.count) % len(deck)
				newDeck[spot] = deck[i]
			}
			deck = newDeck
			break
		default:
			log.Fatal("Unknown command")
		}
	}
	for index, card := range deck {
		if card == 2019 {
			fmt.Println(index)
			break
		}
	}
}
