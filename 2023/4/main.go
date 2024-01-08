package main

import (
    "os"
    "fmt"
    "bufio"
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
    for scanner.Scan() {
        line := scanner.Text()
        numbers := strings.Split(line, ":")[1]
        winningNumbersString := strings.Trim(strings.Split(numbers, "|")[0], " ")
        myNumbersString := strings.Trim(strings.Split(numbers, "|")[1], " ")
        winningNumbers := strings.Split(winningNumbersString, " ")
        myNumbers := strings.Split(myNumbersString, " ")
        winningNumbersMap := make(map[string]bool)
        for _, n := range winningNumbers {
            winningNumbersMap[n] = true
        }
        points := 0
        for _, n := range myNumbers {
            if winningNumbersMap[n] == true {
                if (points == 0) {
                    points = 1
                } else {
                    points = points * 2
                }
            }
        }
        totalPoints += points
    }
    fmt.Println("Total points: ", totalPoints)
}


