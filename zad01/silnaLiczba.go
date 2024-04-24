package main

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func createNick(name, surname string) string {
	name = strings.ToLower(name)
	surname = strings.ToLower(surname)
	name = normalizePolishCharacters(name)
	surname = normalizePolishCharacters(surname)
	return  name[:3] + surname[:3]
}

func normalizePolishCharacters(input string) string {
	replacements := map[rune]rune{
		'ą': 'a',
		'ć': 'c',
		'ę': 'e',
		'ł': 'l',
		'ń': 'n',
		'ó': 'o',
		'ś': 's',
		'ż': 'z',
		'ź': 'z',
	}
	result := strings.Map(func(r rune) rune {
		if replacement, ok := replacements[unicode.ToLower(r)]; ok {
			return replacement
		}
		return r
	}, input)

	return result
}

func stringToByteArrays(s string) []int {
	byteArray := []int{}
	for _, char := range s {
		byteArray = append(byteArray, int(char))
	}
	return byteArray
}

func factorial(n int) *big.Int {
	result := big.NewInt(1)
	for i := 2; i <= n; i++ {
		result = result.Mul(result, big.NewInt(int64(i)))
	}
	return result
}

func containsAllDigits(factorialResult *big.Int, digits []int) bool {
	factorialStr := factorialResult.String()
	for _, digit := range digits {
		digitStr := strconv.Itoa(digit)
		if !strings.Contains(factorialStr, digitStr) {
			return false
		}
	}
	return true
}

func findStrongNumber(asciiCodes []int) int {
	i := 0
	for {
		factorialResult := factorial(i)
		if containsAllDigits(factorialResult, asciiCodes) {
			return i
		}
		i++
	}
}

var calls [31]int
func fibonacci(n int) int {
	calls[n]++
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return fibonacci(n-1) + fibonacci(n-2)
	}
}

func fibonacciTest(n int) int {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return fibonacciTest(n-1) + fibonacciTest(n-2)
	}
}

func findWeakNumber(target int, n int) int {
	fibonacci(n)
	closest := math.MaxInt32
	closestKey := 0
	for key, callCount := range calls {
		if diff := int(math.Abs(float64(target - callCount))); diff < closest {
			closest = diff
			closestKey = key
		}
	}
	return closestKey
}

func ackermann(m, n int) int {
    if m == 0 {
        return n + 1
    } else if n == 0 {
        return ackermann(m-1, 1)
    } else {
        return ackermann(m-1, ackermann(m, n-1))
    }
}


func main() {
	// start := time.Now()
	// ackermann(4, 2)
	// elapsed := time.Since(start)
	// fmt.Println(elapsed)

	var name, surname string

	fmt.Print("Podaj imię: ")
	fmt.Scanln(&name)

	fmt.Print("Podaj nazwisko: ")
	fmt.Scanln(&surname)

	nick := createNick(name, surname)

	asciiCodes := stringToByteArrays(nick)

	fmt.Println("Nick:", nick)

	strongNumber := findStrongNumber(asciiCodes)
	weakNumber := findWeakNumber(strongNumber, 30)
	fmt.Printf("Twoja silna liczba to: %d.\n", strongNumber)
	fmt.Printf("Twoja słaba liczba to: %d.\n", weakNumber)

	start := time.Now()
	fibonacciTest(48)
	elapsed := time.Since(start)
	fmt.Println(elapsed)
}

// dla liczb 138, 19 złożoność:
// (2 -> (n + 3) -> (m-2)) - 3