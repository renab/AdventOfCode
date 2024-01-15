package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type Node struct {
	left  *Node
	right *Node
	name  string
}

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

	nameNodeMap := map[string]*Node{}

	startingNodes := []*Node{}

	directions := []string{}

	start := &Node{name: ""}

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Trim(line, " ")
		if len(strings.Fields(line)) == 1 {
			directions = strings.Split(line, "")
			continue
		} else if len(line) == 0 {
			continue
		}

		line = strings.Replace(line, "(", "", -1)
		line = strings.Replace(line, ")", "", -1)
		line = strings.Replace(line, ",", "", -1)
		spline := strings.Split(line, "=")
		name := strings.Trim(spline[0], " ")

		node, found := nameNodeMap[name]

		nodes := strings.Fields(spline[1])

		left, lfound := nameNodeMap[nodes[0]]
		if !lfound {
			nl := Node{name: nodes[0]}
			nameNodeMap[nodes[0]] = &nl
			left = &nl
			if strings.HasSuffix(nodes[0], "A") {
				startingNodes = append(startingNodes, &nl)
			}
		}

		right, rfound := nameNodeMap[nodes[1]]
		if !rfound {
			nr := Node{name: nodes[1]}
			nameNodeMap[nodes[1]] = &nr
			right = &nr
			if strings.HasSuffix(nodes[1], "A") {
				startingNodes = append(startingNodes, &nr)
			}
		}

		if found {
			node.left = left
			node.right = right
		} else {
			node := Node{name: name, left: left, right: right}
			nameNodeMap[name] = &node
			if node.name == "AAA" {
				start = &node
			}
			if strings.HasSuffix(node.name, "A") {
				startingNodes = append(startingNodes, &node)
			}
		}
	}

	p1Start := time.Now()
	stepCount := GetNumSteps(start, false, directions)
	fmt.Println("Part 1 - ", stepCount, " steps are required to reach the destination. Calculated in ", time.Since(p1Start))

	p2Start := time.Now()
	pathLengths := []int{}
	for _, v := range startingNodes {
		pathLengths = append(pathLengths, GetNumSteps(v, true, directions))
	}
	stepCount2 := ComputeLowestCommonMultiple(pathLengths[0], pathLengths[1], pathLengths[2:]...)
	fmt.Println("Part 2 - ", stepCount2, " steps are required to reach all destinations at the same time. Calculated in ", time.Since(p2Start))
}

func GetNumSteps(start *Node, checkLast bool, directions []string) int {
	destinationReached := false
	stepCount := 0

	var direction string
	nextDirectionPos := 0

	var curNode *Node
	curNode = start

	for !destinationReached {
		direction, nextDirectionPos = GetNextDirection(directions, nextDirectionPos)
		if direction == "L" {
			curNode = curNode.left
		} else {
			curNode = curNode.right
		}
		stepCount++
		if checkLast {
			if strings.HasSuffix(curNode.name, "Z") {
				destinationReached = true
			}
		} else {
			if curNode.name == "ZZZ" {
				destinationReached = true
			}
		}
	}
	return stepCount
}

func GetNextDirection(directions []string, pos int) (string, int) {
	direction := directions[pos]
	nextPos := pos + 1
	if nextPos == len(directions) {
		nextPos = 0
	}
	return direction, nextPos
}

func ComputeLowestCommonMultiple(a, b int, numbers ...int) int {
	result := a * b / ComputeGreatestCommonDivisor(a, b)
	for i := 0; i < len(numbers); i++ {
		result = ComputeLowestCommonMultiple(result, numbers[i])
	}
	return result
}

func ComputeGreatestCommonDivisor(a, b int) int {
	divisor := a
	dividend := b
	if a > b {
		divisor = b
		dividend = a
	}
	r := -1
	for r != 0 {
		r = dividend % divisor
		if r == 0 {
			break
		}
		dividend = divisor
		divisor = r
	}
	return divisor
}
