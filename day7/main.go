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
		if charMap[r] > maxValue && r != 'J' {
			maxValue = charMap[r]
		}
	}

	numJokers := charMap['J']

	if maxValue == 5 || maxValue == 4 {
		// 4 or 5 of a kind -- add 1 to make space for full house
		//
		// We can just add jokers here because we found 4 or 5 of a single
		// non-joker anyway
		return maxValue + 1 + numJokers
	} else if maxValue == 3 && len(charMap) == 2 {
		// full house
		//
		// Jokers upgrade to either 4 or 5 of a kind
		return 4 + numJokers
	} else if maxValue == 3 {
		// 3 of a kind
		//
		// Jokers upgrade to either 4 or 5 of a kind
		if numJokers > 0 {
			return numJokers + 4
		} else {
			return 3
		}
	} else if maxValue == 2 && len(charMap) == 3 {
		// 2 pair
		//
		// Since we found a non-joker pair:
		// - A joker pair upgrades to 4-of-a-kind (5 value)
		// - A single joker upgrades to full house (4 value)
		if numJokers > 0 {
			return numJokers + 3
		} else {
			return 2
		}
	} else if maxValue == 2 {
		// Pair
		//
		// - Single joker upgrades to 3 of a kind (value 3)
		// - 2 jokers upgrades to 4 of a kind (value 5)
		// - 3 jokers upgrades to 5 of a kind (value 6)
		if numJokers > 1 {
			return numJokers + 3
		} else if numJokers == 1 {
			return 3
		} else {
			return 1
		}
	}
	// High card
	//
	// Jokers upgrade to N of a kind
	if numJokers > 3 {
		// 4 or 5 jokers are both 5 of a kind
		return 6
	} else if numJokers == 3 {
		return 5
	} else if numJokers == 2 {
		return 3
	} else if numJokers == 1 {
		return 1
	}
	return 0
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
		byte('J'): -1,
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

	// Compute the bids
	totalWinnings := 0
	for i, hand := range hands {
		totalWinnings += (i + 1) * hand.bid
		//fmt.Println(hand, (i+1)*hand.bid, findTypeValue(hand.cards))
	}
	fmt.Println("winnings", totalWinnings)
}
