package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type hand struct {
	cards string
	bid   int
}

func compareHand(h1 hand, h2 hand) int {
	// First check for hand strength, then compare rank
	return compareHandStrings(h1.cards, h2.cards)

}

func compareHandStrings(s1 string, s2 string) int {
	handTypeRes := compareHandType(s1, s2)
	if handTypeRes != 0 {
		return handTypeRes
	} else {
		return compareHandCard(s1, s2)
	}
}

func compareHandType(h1 string, h2 string) int {
	return findTypeValue(h1) - findTypeValue(h2)
}

func findTypeValue(m string) int {
	// Count unique runs in each string
	charMap := map[rune]int{}
	maxValue := 0
	for _, r := range m {
		charMap[r] += 1
		if charMap[r] > maxValue {
			maxValue = charMap[r]
		}
	}

	if maxValue == 5 || maxValue == 4 {
		// 4 or 5 of a kind -- add 1 to make space for full house
		return maxValue + 1
	} else if maxValue == 3 && len(charMap) == 2 {
		// full house
		return 4
	} else if maxValue == 3 {
		// 3 of a kind
		return 3
	} else if maxValue == 2 && len(charMap) == 3 {
		// 2 pair
		return 2
	} else if maxValue == 2 {
		// Pair
		return 1
	} else if len(charMap) == 5 {
		// High card
		return 0
	}
	// Nothing
	return -1
}

func compareHandCard(h1 string, h2 string) int {
	values := map[byte]int{
		byte('0'): 0,
		byte('1'): 1,
		byte('2'): 2,
		byte('3'): 3,
		byte('4'): 4,
		byte('5'): 5,
		byte('6'): 6,
		byte('7'): 7,
		byte('8'): 8,
		byte('9'): 9,
		byte('T'): 10,
		byte('J'): 11,
		byte('Q'): 12,
		byte('K'): 13,
		byte('A'): 14,
	}
	for i := 0; i < len(h2); i++ {
		c1 := values[h1[i]]
		c2 := values[h2[i]]
		if c1 != c2 {
			return c1 - c2
		}
	}
	return 0
}

func main() {
	f, _ := os.Open("./input.txt")

	scanner := bufio.NewScanner(f)

	// Read the whole thing in, map N map inputs
	hands := []hand{}
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		bid, _ := strconv.Atoi(line[1])
		hand := hand{
			cards: line[0],
			bid:   bid,
		}
		hands = append(hands, hand)
	}

	// sort!
	slices.SortStableFunc(hands, compareHand)
	//fmt.Println(hands)

	// Compute the bids
	totalWinnings := 0
	for i, hand := range hands {
		totalWinnings += (i + 1) * hand.bid
		fmt.Println(hand, (i+1)*hand.bid, findTypeValue(hand.cards))
	}
	fmt.Println("winnings", totalWinnings)
}
