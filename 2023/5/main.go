package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	directionMap = map[string]string{"soil": "fertilizer", "fertilizer": "water", "water": "light", "light": "temperature", "temperature": "humidity", "humidity": "location", "location": ""}
	BATCH_SIZE   = 1000000
)

type SeedRange struct {
	start  int
	length int
}

// type Map struct {
// 	destinationStart int
// 	destinationEnd   int
// 	sourceStart      int
// 	sourceEnd        int
// 	length           int
// 	target           string
// }

type Map struct {
	sourceStart int
	sourceEnd   int
	distance    int
	target      string
}

type Almanac struct {
	Seeds      []int
	SeedRanges []SeedRange
	Maps       map[string][]Map
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

	almanac := Almanac{[]int{}, []SeedRange{}, map[string][]Map{}}
	state := ""

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "seeds:") {
			ProcessSeeds(line, &almanac)
			continue
		}
		if strings.Contains(line, "seed-to-soil") {
			state = "soil"
			continue
		}
		if strings.Contains(line, "soil-to-fertilizer") {
			state = "fertilizer"
			continue
		}
		if strings.Contains(line, "fertilizer-to-water") {
			state = "water"
			continue
		}
		if strings.Contains(line, "water-to-light") {
			state = "light"
			continue
		}
		if strings.Contains(line, "light-to-temp") {
			state = "temperature"
			continue
		}
		if strings.Contains(line, "temperature-to-humidity") {
			state = "humidity"
			continue
		}
		if strings.Contains(line, "humidity-to-location") {
			state = "location"
			continue
		}
		if strings.TrimSpace(line) == "" {
			continue
		}
		ProcessRow(line, state, &almanac)
	}

	// Part 1
	p1Start := time.Now()
	part1Locations := []int{}
	for _, seed := range almanac.Seeds {
		localSeed := seed
		RecursiveMapLookup(&localSeed, &almanac, "soil")
		part1Locations = append(part1Locations, localSeed)
	}
	part1LowestLoc := slices.Min(part1Locations)
	fmt.Println("Part 1 Lowest Location: ", part1LowestLoc, " Calulated in ", time.Since(p1Start))

	// Part 2
	p2Start := time.Now()
	sendChan := make(chan int, len(almanac.SeedRanges))
	var workerWaitGroup sync.WaitGroup
	for _, seedRange := range almanac.SeedRanges {
		localSeedRange := seedRange
		workerWaitGroup.Add(1)
		go Worker(sendChan, &localSeedRange, &almanac, &workerWaitGroup)
	}
	result := -1
	var processorWaitGroup sync.WaitGroup
	processorWaitGroup.Add(1)
	go ProcessResult(sendChan, &result, &processorWaitGroup)
	workerWaitGroup.Wait()
	close(sendChan)
	processorWaitGroup.Wait()

	fmt.Println("Part 2 Lowest Location: ", result, " Calulated in ", time.Since(p2Start))
}

func ProcessSeeds(line string, almanac *Almanac) {
	re := regexp.MustCompile(`(\d+)`)
	seeds := re.FindAllString(line, -1)
	seedRanges := []SeedRange{}
	for _, seed := range seeds {
		seedInt, _ := strconv.Atoi(seed)
		almanac.Seeds = append(almanac.Seeds, seedInt)
	}
	for i := 0; i < len(seeds); i = i + 2 {
		start, _ := strconv.Atoi(seeds[i])
		remaining, _ := strconv.Atoi(seeds[i+1])
		for remaining > BATCH_SIZE {
			seedRanges = append(seedRanges, SeedRange{start, BATCH_SIZE})
			remaining = remaining - BATCH_SIZE
			start = start + BATCH_SIZE
		}
		if remaining > 0 {
			seedRanges = append(seedRanges, SeedRange{start, remaining})
		}
	}
	fmt.Println(seedRanges)
	almanac.SeedRanges = seedRanges
}

func ProcessRow(row string, whichMap string, almanac *Almanac) {
	re := regexp.MustCompile(`(\d+)`)
	entries := re.FindAllString(row, -1)
	nums := []int{}
	for _, entry := range entries {
		num, _ := strconv.Atoi(entry)
		nums = append(nums, num)
	}
	almanac.Maps[whichMap] = append(almanac.Maps[whichMap], Map{nums[1], nums[1] + nums[2], nums[0] - nums[1], directionMap[whichMap]})
}

func RecursiveMapLookup(source *int, almanac *Almanac, destinationType string) {
	for _, v := range almanac.Maps[destinationType] {
		if *source >= v.sourceStart && *source < v.sourceEnd {
			*source = *source + v.distance
			if v.target != "" {
				RecursiveMapLookup(source, almanac, v.target)
				return
			}
		}
	}
	if directionMap[destinationType] != "" {
		RecursiveMapLookup(source, almanac, directionMap[destinationType])
	}
}

func Worker(result chan<- int, seedRange *SeedRange, almanac *Almanac, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	results := []int{}
	for i := seedRange.start; i <= seedRange.start+seedRange.length; i++ {
		iterationResult := i
		RecursiveMapLookup(&iterationResult, almanac, "soil")
		results = append(results, iterationResult)
	}
	location := slices.Min(results)
	result <- location
}

func ProcessResult(receive <-chan int, result *int, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	locations := []int{}
	for location := range receive {
		locations = append(locations, location)
	}
	*result = slices.Min(locations)
}
