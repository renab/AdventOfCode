package main

import (
	"bufio"
	"fmt"
	"os"
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

	times := []int{}
	distances := []int{}
	var time int
	var distance int

	for scanner.Scan() {
		line := scanner.Text()
		spLine := strings.Split(line, ":")
		if spLine[0] == "Time" {
			sTimes := strings.Fields(spLine[1])
			for _, sTime := range sTimes {
				time, _ := strconv.Atoi(string(sTime))
				times = append(times, time)
			}
			sTime := strings.Join(sTimes, "")
			time, _ = strconv.Atoi(sTime)
		} else {
			sDists := strings.Fields(spLine[1])
			for _, sDist := range sDists {
				dist, _ := strconv.Atoi(string(sDist))
				distances = append(distances, dist)
			}
			sDist := strings.Join(sDists, "")
			distance, _ = strconv.Atoi(sDist)
		}
	}
	winningResults := [][]int{}
	for i, time := range times {
		winningResults = append(winningResults, CalculatePossibleVictory(distances[i], time))
	}
	totalMarginOfError := 1
	for _, winningResultSlice := range winningResults {
		totalMarginOfError = len(winningResultSlice) * totalMarginOfError
	}
	fmt.Println("Part 1 - Margin of Error: ", totalMarginOfError)
	waysToWinP2 := CalculatePossibleVictory(distance, time)
	fmt.Println("Part 2 - Number of ways to win the race: ", len(waysToWinP2))
}

func CalculatePossibleVictory(minDist int, time int) []int {
	ret := []int{}
	for i := 1; i < time; i++ {
		// d = x(y-x) where d = distance, x = time held down, y = available time
		if i*(time-i) > minDist {
			ret = append(ret, i*(time-i))
		}
	}
	return ret
}
