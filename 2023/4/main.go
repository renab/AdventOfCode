package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <filename>")
		os.Exit(1)
	}
	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		os.Exit(1)
	}
	defer file.Close()

    scanner := bufio.NewScanner(file)

    totalPoints := 0
    totalScratchCards := 0

    scratchCardInstanceMap := make(map[int]int)

    spaceRegex, _ := regexp.Compile("^\\s*$")
    gameNumRegex, _ := regexp.Compile("(\\d+)")

    originalScratchCards := 0

    for scanner.Scan() {
        originalScratchCards++
        line := scanner.Text()
        gameNum, _ := strconv.Atoi(gameNumRegex.FindAllString(strings.Split(line, ":")[0], -1)[0])
        if _, ok := scratchCardInstanceMap[gameNum]; !ok {
            scratchCardInstanceMap[gameNum] = 1
        } else {
            scratchCardInstanceMap[gameNum]++
        }
        numbers := strings.Split(line, ":")[1]
        winningNumbersString := strings.Split(numbers, "|")[0]
        myNumbersString := strings.Split(numbers, "|")[1]
        winningNumbers := strings.Split(winningNumbersString, " ")
        myNumbers := strings.Split(myNumbersString, " ")
        winningNumbersMap := make(map[string]bool)
        for _, n := range winningNumbers {
            if !spaceRegex.MatchString(n) { // ignore entries that are just whitespace
                winningNumbersMap[n] = true
            }
        }
        points := 0
        matches := 0
        for _, n := range myNumbers {
            if winningNumbersMap[n] == true {
                matches++
                if (points == 0) {
                    points = 1
                } else {
                    points = points * 2
                }
            }
        }

        // Calculate the number of scratch cards we won and increment the instance map by the number of instances of *this* scratch card
        for i := 1; i <= matches; i++ {
            _, ok := scratchCardInstanceMap[gameNum + i]
            if !ok {
                scratchCardInstanceMap[gameNum + i] = scratchCardInstanceMap[gameNum]
            } else {
                scratchCardInstanceMap[gameNum + i] += scratchCardInstanceMap[gameNum]
            }
        }
        totalPoints += points
    }

    for k, v := range scratchCardInstanceMap {
        // We can't invent new cards, so we can only count instances of the cards we started with
        if k <= originalScratchCards {
            totalScratchCards += v
        }
    }

    fmt.Println("Total points (except there are no points): ", totalPoints)
    fmt.Println("Total scratch cards (because we need to read the rules on the back of the card): ", totalScratchCards)
}


