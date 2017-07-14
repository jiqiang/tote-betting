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

	winWagerRegExp    = `^bet:[w|p]:\d+:\d+(\.\d+)?$`
	exactaWagerRegExp = `^bet:e:\d+\,\d+:\d+(\.\d+)?$`
	resultRegExp      = `^result:\d+:\d+:\d+$`
)

var winPool = make(map[string]float64)
var placePool = make(map[string]float64)
var exactaPool = make(map[string]float64)

// valid returns true if input line format matches to designed format.
// Otherwise returns false
func valid(line string) bool {
	if strings.HasPrefix(line, "bet:w") || strings.HasPrefix(line, "bet:p") {
		return validWager(line, winWagerRegExp)
	}

	if strings.HasPrefix(line, "bet:e") {
		return validWager(line, exactaWagerRegExp)
	}

	if strings.HasPrefix(line, "result") {
		return validWager(line, resultRegExp)
	}

	return false
}

// validWager checks if given line matches given format
func validWager(line string, wagerFormat string) bool {
	var format = regexp.MustCompile(wagerFormat)
	return format.MatchString(line)
}

// isResultLine returns true if line starting with "result"
func isResultLine(line string) bool {
	return strings.HasPrefix(line, "result")
}

// purify removes whitespace, newline, tabs and converts string to lowercase.
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

// parse parses input line and proccess it
func parse(line string) {
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

// getSum gets sum of values from a map
func getSum(pool map[string]float64) float64 {
	var sum float64
	sum = 0
	for _, value := range pool {
		sum = sum + value
	}
	return sum
}

// printDividendsResult prints final dividends result
func printDividendsResult(line string) {
	parts := strings.Split(line, ":")
	first := parts[1]
	second := parts[2]
	third := parts[3]

	fmt.Printf("Win:%s:$%.2f\n", first, calculateWinDividends(first))

	firstResult, secondResult, thirdResult := calculatePlaceDividends(first, second, third)
	fmt.Printf("Place:%s:$%.2f\n", first, firstResult)
	fmt.Printf("Place:%s:$%.2f\n", second, secondResult)
	fmt.Printf("Place:%s:$%.2f\n", third, thirdResult)

	fmt.Printf("Exacta:%s,%s:$%.2f\n", first, second, calculateExactaDividends(first, second))
}

// calculateWinDividends returns dividends for Win product
func calculateWinDividends(first string) float64 {
	stake, ok := winPool[first]
	if !ok {
		return 0
	}
	sum := getSum(winPool)
	amount := round(sum-sum*winPoolCommission, 0.01)
	dividends := getDividends(amount, stake)
	return dividends
}

// calculatePlaceDividends returns dividends for Place product
func calculatePlaceDividends(first, second, third string) (float64, float64, float64) {
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

// calculateExactaDividends returns dividends for Exacta product
func calculateExactaDividends(first string, second string) float64 {
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

	// read input from stdin
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()

		line = purify(line)

		// validate input to match designed format
		if !valid(line) {
			fmt.Printf("Invalid input: [%s]\n", line)
			continue
		}

		// print dividents result and exit
		if isResultLine(line) {
			printDividendsResult(line)
			break
		}

		parse(line)
	}
}
