package main

import (
    "os"
    "fmt"
    "bufio"
    "regexp"
    "strconv"
    "strings"
)

func main() {
    if len(os.Args) != 5 {
        fmt.Println("Usage: go run main.go <r> <g> <b> <filename>")
        os.Exit(1)
    }

    maxRed, err := strconv.Atoi(os.Args[1])
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }

    maxGreen, err := strconv.Atoi(os.Args[2])
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }

    maxBlue, err := strconv.Atoi(os.Args[3])
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }

    maxMap := make(map[string]int)
    maxMap["red"] = maxRed
    maxMap["green"] = maxGreen
    maxMap["blue"] = maxBlue

    fmt.Println("Max Red Blocks: ", maxMap["red"], " Max Green Blocks: ", maxMap["green"], " Max Blue Blocks: ", maxMap["blue"], " Input File: ", os.Args[4])

    filename := os.Args[4]
    file, err := os.Open(filename)
    if (err != nil) {
        fmt.Println("Error:", err)
        os.Exit(1)
    }
    defer file.Close()

    gameRegex, _ := regexp.Compile(`Game (\d+)`)
    redRegex, _ := regexp.Compile(`(\d+) red`)
    greenRegex, _ := regexp.Compile(`(\d+) green`)
    blueRegex, _ := regexp.Compile(`(\d+) blue`)

    possibleMap := make(map[int]bool)
    
    minBlocksPerGameMap := make(map[int]map[string]int)

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        game, _ := strconv.Atoi(gameRegex.FindStringSubmatch(strings.Split(line, ":")[0])[1])
        results := strings.Split(strings.Split(line, ":")[1], ";")
        possible := true
        for _, result := range results {
            if redRegex.MatchString(result) {
                red, _ := strconv.Atoi(redRegex.FindStringSubmatch(result)[1])
                if red > maxMap["red"] {
                    possible = false 
                }
                gameMap := minBlocksPerGameMap[game]
                if gameMap == nil {
                    gameMap = make(map[string]int)
                    gameMap["red"] = red
                    minBlocksPerGameMap[game] = gameMap
                } else if gameMap["red"] < red {
                    gameMap["red"] = red
                }
            }
            if greenRegex.MatchString(result) {
                green, _ := strconv.Atoi(greenRegex.FindStringSubmatch(result)[1])
                if green > maxMap["green"] {
                    possible = false
                }
                gameMap := minBlocksPerGameMap[game]
                if gameMap == nil {
                    gameMap = make(map[string]int)
                    gameMap["green"] = green
                    minBlocksPerGameMap[game] = gameMap
                } else if gameMap["green"] < green {
                    gameMap["green"] = green
                }
            }
            if blueRegex.MatchString(result) {
                blue, _ := strconv.Atoi(blueRegex.FindStringSubmatch(result)[1])
                if blue > maxMap["blue"] {
                    possible = false
                }
                gameMap := minBlocksPerGameMap[game]
                if gameMap == nil {
                    gameMap = make(map[string]int)
                    gameMap["blue"] = blue
                    minBlocksPerGameMap[game] = gameMap
                } else if gameMap["blue"] < blue {
                    gameMap["blue"] = blue
                }
            }
        }
        if possible {
            possibleMap[game] = true
        } else {
            possibleMap[game] = false
        }
    }

    total := 0
    for k, v := range possibleMap {
        if v == true {
            total += k
        }
    }
    fmt.Println("Total: ", total)
    
    power := 0
    for _, v := range minBlocksPerGameMap {
        power += v["red"] * v["green"] * v["blue"]
    }
    fmt.Println("Total Power: ", power)
}
