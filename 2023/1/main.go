package main

import (
    "fmt"
    "os"
    "bufio"
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
    if (err != nil) {
        fmt.Println("Error:", err)
        os.Exit(1)
    }

    scanner := bufio.NewScanner(file)
    numLines := 0
    lineMap := make(map[int][]string)
    // regex, err := regexp.Compile(`(?=(one|two|three|four|five|six|seven|eight|nine|\d))`)
    // Golang doesn't support look* or overlaps, so we have to do some string replacement.

    regex, err := regexp.Compile(`(\d)`)
    for scanner.Scan() {        
        lineMap[numLines] = regex.FindAllString(NumberReplace(scanner.Text()), -1)
        numLines++
    }
    defer file.Close()
    result := 0
    for _, line := range lineMap {
        if len (line) != 0 {
            numStr := ""
            if (len (line[0]) == 1) {
                numStr = line[0] 
            } else {
                numStr = getNumStr(line[0])   
            }
            if (len (line[len (line)-1]) == 1 ) {
                numStr = numStr + line[len(line)-1]
            } else {
                numStr = numStr + getNumStr(line[len(line)-1]) 
            }
            num, err := strconv.Atoi(numStr)
            if err != nil {
               fmt.Println("Error:", err)
               os.Exit(1)
            }
            // fmt.Println(line[0], " ", line[len(line)-1], " ", numStr, " ", num)
            result += num
        }
    }
    fmt.Println(result)
}

func NumberReplace (str string) string {
    localStr := strings.Replace(str, "one", "o1e", -1)
    localStr = strings.Replace(localStr, "two", "t2o", -1)
    localStr = strings.Replace(localStr, "three", "t3e", -1)
    localStr = strings.Replace(localStr, "four", "f4r", -1)
    localStr = strings.Replace(localStr, "five", "f5e", -1)
    localStr = strings.Replace(localStr, "six", "s6x", -1)
    localStr = strings.Replace(localStr, "seven", "s7n", -1)
    localStr = strings.Replace(localStr, "eight", "e8t", -1)
    localStr = strings.Replace(localStr, "nine", "n9e", -1)
    return localStr
}

func getNumStr (numStr string) string {
    switch numStr {
    case "one":
        return "1"
    case "two":
        return "2"
    case "three":
        return "3"
    case "four":
        return "4"
    case "five":
        return "5"
    case "six":
        return "6"
    case "seven":
        return "7"
    case "eight":
        return "8"
    case "nine":
        return "9"
    default:
        fmt.Println("Default: ", numStr)
        return "-1"
    }
}
