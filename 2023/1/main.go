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
    defer file.Close()
    
    scanner := bufio.NewScanner(file)
    numLines := 0
    lineMap := make(map[int][]string)

    regex, _ := regexp.Compile(`(\d)`)
    for scanner.Scan() {        
        lineMap[numLines] = regex.FindAllString(NumberReplace(scanner.Text()), -1)
        numLines++
    }
    
    result := 0
    for _, line := range lineMap {
        if len (line) != 0 {
            numStr := ""
            numStr = line[0] 
            numStr = numStr + line[len(line)-1]
            num, err := strconv.Atoi(numStr)
            if err != nil {
               fmt.Println("Error:", err)
               os.Exit(1)
            }
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
