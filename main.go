package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

const (
	winPoolCommission    = 0.15
	placePoolCommission  = 0.12
	exactaPoolCommission = 0.18
)

var winPool = make(map[string]float64)
var placePool = make(map[string]float64)
var exactaPool = make(map[string]float64)

// validateInput returns true if input line format matches to designed format.
// Otherwise returns false
func validInput(line string) bool {
	var validWinAndPlace = regexp.MustCompile(`^bet:[w|p]:\d+:\d+(\.\d+)?$`)
	var validExacta = regexp.MustCompile(`^bet:e:\d+\,\d+:\d+(\.\d+)?$`)
	var validResult = regexp.MustCompile(`^result:\d+:\d+:\d+$`)

	if strings.HasPrefix(line, "bet:w") || strings.HasPrefix(line, "bet:p") {
		return validWinAndPlace.MatchString(line)
	}

	if strings.HasPrefix(line, "bet:e") {
		return validExacta.MatchString(line)
	}

	if strings.HasPrefix(line, "result") {
		return validResult.MatchString(line)
	}

	return false
}

// purify removes whitespace, newline, tabs, converts string to lowercase.
func purify(inStr string) string {
	var outStr string
	for _, c := range inStr {
		if !unicode.IsSpace(c) {
			outStr = outStr + string(c)
		}
	}
	outStr = strings.ToLower(outStr)
	return outStr
}

// processBet parses input line and proccess it
func processBet(line string) {
	parts := strings.Split(line, ":")
	product := parts[1]
	selections := parts[2]
	stake, _ := strconv.ParseFloat(parts[3], 64)
	stake = round(stake, 0.01)

	switch product {
	case "w":
		addToPool(winPool, selections, stake)
	case "p":
		addToPool(placePool, selections, stake)
	case "e":
		addToPool(exactaPool, selections, stake)
	}
}

// addToPool adds stake to the pool
func addToPool(pool map[string]float64, selections string, stake float64) {
	v, ok := pool[selections]
	if ok {
		pool[selections] = v + stake
	} else {
		pool[selections] = stake
	}
}

// round float to nearest value. e.g 0.01
func round(f, nearest float64) float64 {
	return float64(int64(f/nearest+0.5)) * nearest
}

func getSum(pool map[string]float64) float64 {
	var sum float64
	sum = 0
	for _, value := range pool {
		sum = sum + value
	}
	return sum
}

func printResult(line string) {
	parts := strings.Split(line, ":")
	first := parts[1]
	second := parts[2]
	third := parts[3]

	fmt.Printf("Win:%s:$%.2f\n", first, calculateWinResult(first))

	firstResult, secondResult, thirdResult := calculatePlaceResult(first, second, third)
	fmt.Printf("Place:%s:$%.2f\n", first, firstResult)
	fmt.Printf("Place:%s:$%.2f\n", second, secondResult)
	fmt.Printf("Place:%s:$%.2f\n", third, thirdResult)

	fmt.Printf("Exacta:%s,%s:$%.2f\n", first, second, calculateExactaResult(first, second))
}

func calculateWinResult(first string) float64 {
	stake, ok := winPool[first]
	if !ok {
		return 0
	}
	sum := getSum(winPool)
	amount := round(sum-sum*winPoolCommission, 0.01)
	dividends := getDividends(amount, stake)
	return dividends
}

func calculatePlaceResult(first, second, third string) (float64, float64, float64) {
	var firstDividends, secondDividends, thirdDividends float64 = 0, 0, 0

	sum := getSum(placePool)
	amount := sum * (1 - placePoolCommission) / 3

	firstStake, firstOk := placePool[first]
	if firstOk {
		firstDividends = getDividends(amount, firstStake)
	}

	secondStake, secondOk := placePool[second]
	if secondOk {
		secondDividends = getDividends(amount, secondStake)
	}

	thirdStake, thirdOk := placePool[third]
	if thirdOk {
		thirdDividends = getDividends(amount, thirdStake)
	}

	return firstDividends, secondDividends, thirdDividends
}

func calculateExactaResult(first string, second string) float64 {
	stake, ok := exactaPool[first+","+second]
	if !ok {
		return 0
	}
	sum := getSum(exactaPool)
	amount := sum * (1 - exactaPoolCommission)
	dividends := getDividends(amount, stake)
	return dividends
}

func getDividends(amount float64, stake float64) float64 {
	return round(amount/stake, 0.01)
}

func main() {

	reader := bufio.NewReader(os.Stdin)

	for {
		line, _ := reader.ReadString('\n')
		line = purify(line)

		// validate input to match designed format
		if !validInput(line) {
			fmt.Printf("Invalid input, ignore %s\n", line)
			break
		}

		if strings.HasPrefix(line, "result") {
			printResult(line)
			break
		}

		processBet(line)
	}
}
