package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

var (
	cardValueMap  = map[string]int{"2": 1, "3": 2, "4": 3, "5": 4, "6": 5, "7": 6, "8": 7, "9": 8, "T": 9, "J": 10, "Q": 11, "K": 12, "A": 13}
	jCardValueMap = map[string]int{"J": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9, "T": 10, "Q": 11, "K": 12, "A": 13}
	rankValueMap  = map[string]int{"HC": 1, "P": 2, "2P": 3, "3oaK": 4, "FH": 5, "4oaK": 6, "5oaK": 7}
)

type Hand struct {
	cards  string
	rank   string
	values []int
	bid    int
}

type SortedHand []Hand

func (h SortedHand) Len() int {
	return len(h)
}

func (h SortedHand) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h SortedHand) Less(i, j int) bool {
	if rankValueMap[h[i].rank] < rankValueMap[h[j].rank] {
		return true
	} else if rankValueMap[h[i].rank] > rankValueMap[h[j].rank] {
		return false
	}
	for k := 0; k < 5; k++ {
		if h[i].values[k] == h[j].values[k] {
			continue
		} else if h[i].values[k] < h[j].values[k] {
			return true
		} else {
			return false
		}
	}
	return false
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

	hands := []Hand{}
	jHands := []Hand{}

	for scanner.Scan() {
		line := scanner.Text()
		spLine := strings.Split(line, " ")
		bid, _ := strconv.Atoi(spLine[1])
		hand, jHand := CreateHand(strings.Trim(spLine[0], " "), bid)
		hands = append(hands, hand)
		jHands = append(jHands, jHand)
	}

	sort.Stable(SortedHand(hands))
	winnings := 0
	for i := 1; i <= len(hands); i++ {
		winnings = winnings + (hands[i-1].bid * i)
	}

	fmt.Println("Part 1 - Total winnings: ", winnings)

	sort.Stable(SortedHand(jHands))
	jWinnings := 0
	for i := 1; i <= len(jHands); i++ {
		jWinnings = jWinnings + (jHands[i-1].bid * i)
	}

	fmt.Println("Part 2 - Total winnings with Jokers: ", jWinnings)
}

func CreateHand(cards string, bid int) (Hand, Hand) {
	cardMap := map[string]int{}
	hand := Hand{cards: cards, bid: bid}
	jHand := Hand{cards: cards, bid: bid}
	cardValues := []int{}
	jCardValues := []int{}
	for i := 0; i < len(cards); i++ {
		num, err := cardMap[string(cards[i])]
		if !err {
			cardMap[string(cards[i])] = 1
		} else {
			cardMap[string(cards[i])] = num + 1
		}
		cardValues = append(cardValues, cardValueMap[string(cards[i])])
		jCardValues = append(jCardValues, jCardValueMap[string(cards[i])])
	}
	hand.values = cardValues
	jHand.values = jCardValues

	if len(cardMap) == 1 {
		hand.rank = "5oaK"
	} else if len(cardMap) == 2 {
		for _, v := range cardMap {
			if v == 2 || v == 3 {
				hand.rank = "FH"
			} else {
				hand.rank = "4oaK"
			}
			break
		}
	} else if len(cardMap) == 3 {
		for _, v := range cardMap {
			if v == 1 {
				continue
			} else if v == 2 {
				hand.rank = "2P"
			} else {
				hand.rank = "3oaK"
			}
			break
		}
	} else if len(cardMap) == 4 {
		hand.rank = "P"
	} else {
		hand.rank = "HC"
	}

	var numJoker int
	var found bool
	if len(cardMap) == 1 {
		jHand.rank = "5oaK"
	} else if len(cardMap) == 2 {
		_, found = cardMap["J"]
		if found {
			jHand.rank = "5oaK"
		} else {
			for _, v := range cardMap {
				if v == 2 || v == 3 {
					jHand.rank = "FH"
				} else {
					jHand.rank = "4oaK"
				}
				break
			}
		}
	} else if len(cardMap) == 3 {
		numJoker, found = cardMap["J"]
		if found {
			for k, v := range cardMap {
				if k == "J" {
					continue
				}
				if numJoker == 1 && v == 2 {
					jHand.rank = "FH"
				} else {
					jHand.rank = "4oaK"
				}
				break
			}
		} else {
			for _, v := range cardMap {
				if v == 1 {
					continue
				} else if v == 2 {
					jHand.rank = "2P"
				} else {
					jHand.rank = "3oaK"
				}
				break
			}
		}
	} else if len(cardMap) == 4 {
		_, found = cardMap["J"]
		if found {
			jHand.rank = "3oaK"
		} else {
			jHand.rank = "P"
		}
	} else {
		_, found = cardMap["J"]
		if found {
			jHand.rank = "P"
		} else {
			jHand.rank = "HC"
		}
	}
	return hand, jHand
}
