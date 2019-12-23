package main

import (
	"bufio"
	"fmt"
	"log"
	"math/big"
	"os"
	"regexp"
	"strconv"
)

var DECK = big.NewInt(119315717514047)
var TIMES = big.NewInt(101741582076661)

//card i always in position ai+b (mod N)
var a = big.NewInt(1)
var b = big.NewInt(0)

func matrixMult(a []*big.Int, b []*big.Int) []*big.Int {
	TEMP := big.NewInt(0)
	A := big.NewInt(0)
	B := big.NewInt(0)
	C := big.NewInt(0)
	D := big.NewInt(0)

	TEMP.Mul(a[1], b[2])
	A.Mul(a[0], b[0]).Add(A, TEMP).Mod(A, DECK)
	TEMP.Mul(a[1], b[3])
	B.Mul(a[0], b[1]).Add(B, TEMP).Mod(B, DECK)
	TEMP.Mul(a[3], b[2])
	C.Mul(a[2], b[0]).Add(C, TEMP).Mod(C, DECK)
	TEMP.Mul(a[3], b[3])
	D.Mul(a[2], b[1]).Add(D, TEMP).Mod(D, DECK)

	return []*big.Int{A, B, C, D}
}

var ZERO = big.NewInt(0)

func matrixPower(mat []*big.Int, exp *big.Int) []*big.Int {
	mul := mat[:]
	var ans = []*big.Int{big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(1)}
	for exp.Cmp(ZERO) > 0 {
		mod := big.NewInt(0)
		mod.Mod(exp, TWO)
		if mod.Cmp(ONE) == 0 {
			ans = matrixMult(mul, ans)
		}
		exp.Div(exp, big.NewInt(2))
		mul = matrixMult(mul, mul)
	}
	return ans
}

var TWO = big.NewInt(2)
var ONE = big.NewInt(1)

func inversePrime(num *big.Int, p *big.Int) *big.Int {
	exp := big.NewInt(0).Sub(p, TWO)
	var ans = big.NewInt(1)
	for exp.Cmp(ZERO) > 0 {
		mod := big.NewInt(0)
		mod.Mod(exp, TWO)
		if mod.Cmp(ONE) == 0 {
			ans.Mul(ans, num).Mod(ans, p)
		}
		exp.Div(exp, big.NewInt(2))
		num.Mul(num, num).Mod(num, p)
	}
	return ans
}

func parseCommand(command string) {
	cutParser := regexp.MustCompile(`^cut (.*\d+)$`)
	dealParser := regexp.MustCompile(`^deal with increment (\d+)$`)
	if cutParser.MatchString(command) {
		data := cutParser.FindStringSubmatch(command)
		count, _ := strconv.Atoi(data[1])
		b.Sub(b, big.NewInt(int64(count))).Mod(b, DECK)
	}
	if "deal into new stack" == command {
		a.Neg(a).Mod(a, DECK)

		b.Add(b, ONE).Neg(b).Mod(b, DECK)
		fmt.Println("B2 = ", b)
	}
	if dealParser.MatchString(command) {
		data := dealParser.FindStringSubmatch(command)
		count, _ := strconv.Atoi(data[1])
		a.Mul(a, big.NewInt(int64(count))).Mod(a, DECK)

		b.Mul(b, big.NewInt(int64(count))).Mod(b, DECK)
		fmt.Println("B3 = ", b)
	}
}

func main() {
	argsWithoutProg := os.Args[1:]
	file, err := os.Open(argsWithoutProg[0])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parseCommand(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(a)
	fmt.Println(b)
	res := matrixPower([]*big.Int{a, b, big.NewInt(0), big.NewInt(1)}, TIMES)
	ansA, ansB := res[0], res[1]
	fmt.Println(ansA)
	fmt.Println(ansB)
	var ans = big.NewInt(0)
	offset := big.NewInt(2020)
	ans.Mul(inversePrime(ansA, DECK), offset.Sub(offset, ansB)).Mod(ans, DECK)
	fmt.Println(ans)
}
