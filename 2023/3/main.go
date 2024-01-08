package main

import (
	"bufio"
	"fmt"
	"math"
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

	matrix := make(map[int][]string)

	scanner := bufio.NewScanner(file)
	matrixWidth := -1
	row := 0
	for scanner.Scan() {
		line := scanner.Text()
		if matrixWidth == -1 {
			matrixWidth = len(line)
		}
		matrix[row] = strings.Split(line, "")
		row++
	}

	partNumberSum := 0
    gearRatioSum := 0

	for k, v := range matrix {
		if len(v) != matrixWidth {
			fmt.Println("Error: matrix is not rectangular")
			os.Exit(1)
		}
		inNumber := false
		numberStart := -1
		numberEnd := -1
		for i, s := range v {
			if !inNumber && IsNumber(s) {
				inNumber = true
				numberStart = i
			}
			if (inNumber && !IsNumber(s)) {
				inNumber = false
				numberEnd = i
            // If the number is at the end of the line
			} else if (inNumber && i + 1 == len(v)) {
                inNumber = false
                numberEnd = i + 1
            }
			if numberStart != -1 && numberEnd != -1 {
				isPartNumber := CheckPartNumber(numberStart, numberEnd, k, matrix, matrixWidth)
				if isPartNumber {
                    // PrintNumberAndSurround(numberStart, numberEnd, k, matrix, matrixWidth)
					numStr := strings.Join(v[numberStart:numberEnd], "")
					num, err := strconv.Atoi(numStr)
					if err != nil {
						fmt.Println("Error: ", err)
						os.Exit(1)
					}
					partNumberSum += num
				}
				numberStart = -1
				numberEnd = -1
			}
            if v[i] == "*" {
                // PrintNumberAndSurround(i, i, k, matrix, matrixWidth)
                isGear, ratioComponents := CheckGear(k, i, matrix, matrixWidth)
                if isGear {
                    gearRatioSum += ratioComponents[0] * ratioComponents[1]
                }
            }
		}
	}

	fmt.Println("Sum of all part numbers: ", partNumberSum)
    fmt.Println("Sum of all gear ratios: ", gearRatioSum)
}

func CheckGear(row int, col int, matrix map[int][]string, matrixWidth int) (bool, []int) {
    adjacentCount := 0
    numberArr := make([]int, 2)
    // Check left and right
    if col - 1 >= 0 && IsNumber(matrix[row][col - 1]) {
        numberArr[adjacentCount] = GetNumAtPos(row, col - 1, matrix, matrixWidth)
        adjacentCount++
    }
    if col + 1 < matrixWidth && IsNumber(matrix[row][col + 1]) {
        numberArr[adjacentCount] = GetNumAtPos(row, col + 1, matrix, matrixWidth)
        adjacentCount++
    }
    // Check up
    if row - 1 >= 0 {
        startingCol := col
        endingCol := col
        if col - 1 >= 0 {
            startingCol--
        }
        if col + 1 < matrixWidth {
            endingCol++
        }
        for i := startingCol; i <= endingCol; i++ {
            if adjacentCount > 2 {
                break
            }
            if IsNumber(matrix[row - 1][i]) {
                if i != startingCol && !IsNumber(matrix[row - 1][i - 1]) {
                    adjacentCount++
                    if adjacentCount <= 2 {
                        numberArr[adjacentCount - 1] = GetNumAtPos(row - 1, i, matrix, matrixWidth)
                    }
                } else if i == startingCol {
                    adjacentCount++
                    if adjacentCount <= 2 {
                        numberArr[adjacentCount - 1] = GetNumAtPos(row - 1, i, matrix, matrixWidth)
                    }
                }
            }
        }
    }
    // Check down
    if row + 1 < len(matrix) {
        startingCol := col
        endingCol := col
        if col - 1 >= 0 {
            startingCol--
        }
        if col + 1 < matrixWidth {
            endingCol++
        }
        for i := startingCol; i <= endingCol; i++ {
            if adjacentCount > 2 {
                break
            }
            if IsNumber(matrix[row + 1][i]) {
                if i != startingCol && !IsNumber(matrix[row + 1][i - 1]) {
                    adjacentCount++
                    if adjacentCount <= 2 {
                        numberArr[adjacentCount - 1] = GetNumAtPos(row + 1, i, matrix, matrixWidth)
                    }
                } else if i == startingCol {
                    adjacentCount++
                    if adjacentCount <= 2 {
                        numberArr[adjacentCount - 1] = GetNumAtPos(row + 1, i, matrix, matrixWidth)
                    }
                }
            }
        }

    }
    return adjacentCount == 2, numberArr
}

func GetNumAtPos(row int, col int, matrix map[int][]string, matrixWidth int) int {
    currentCol := col
    for IsNumber(matrix[row][currentCol]) {
        currentCol--
        if currentCol < 0 {
            break
        }
    }
    numStart := currentCol + 1
    currentCol = col
    for IsNumber(matrix[row][currentCol]) {
        currentCol++
        if currentCol >= matrixWidth {
            break
        }
    }
    numEnd := currentCol
    numStr := strings.Join(matrix[row][numStart:numEnd], "")
    num, _ := strconv.Atoi(numStr)
    return num
}    

func PrintNumberAndSurround(numberStart int, numberEnd int, row int, matrix map[int][]string, matrixWidth int) {
	fmt.Println()
    startingCol := numberStart - 1
    endingCol := numberEnd + 1
    if startingCol < 0 {
        startingCol += 1
    }
    if endingCol > matrixWidth {
        endingCol -= 1
    }
    if row > 0 {
		for i := startingCol; i <= endingCol; i++ {
			fmt.Print(matrix[row-1][i])
		}
	}
	fmt.Println()
	for i := startingCol; i <= endingCol; i++ {
		fmt.Print(matrix[row][i])
	}
	fmt.Println()
	if row < len(matrix)-1 {
		for i := startingCol; i <= endingCol; i++ {
			fmt.Print(matrix[row+1][i])
		}
	}
	fmt.Println()
}

func CheckPartNumber(numberStart int, numberEnd int, row int, matrix map[int][]string, matrixWidth int) bool {
	// PrintNumberAndSurround(numberStart, numberEnd, row, matrix, matrixWidth)
	colStart := math.Max(0, float64(numberStart-1))
	colEnd := math.Min(float64(matrixWidth), float64(numberEnd+1))

	// Check left and right
	if numberStart > 0 && IsSymbol(matrix[row][numberStart-1]) {
		return true
	}
	if numberEnd < matrixWidth && IsSymbol(matrix[row][numberEnd]) {
		return true
	}

	// Check up
	if row > 0 {
		for i := colStart; i < colEnd; i++ {
			if IsSymbol(matrix[row-1][int(i)]) {
				return true
			}
		}
	}

	// Check down
	if row < len(matrix)-1 {
		for i := colStart; i < colEnd; i++ {
			if IsSymbol(matrix[row+1][int(i)]) {
				// PrintNumberAndSurround(numberStart, numberEnd, row, matrix, matrixWidth)
				return true
			}
		}
	}

	// PrintNumberAndSurround(numberStart, numberEnd, row, matrix, matrixWidth)
	return false
}

func IsSymbol(s string) bool {
	symbolRegex, _ := regexp.Compile("[^\\d\\.\\s]")
	return symbolRegex.MatchString(s)
}

func IsNumber(s string) bool {
	isNumber := false
	numberRegex, _ := regexp.Compile("\\d")
	if numberRegex.MatchString(s) {
		isNumber = true
	}
	return isNumber
}
