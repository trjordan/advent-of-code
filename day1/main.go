package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func main() {
	f, _ := os.Open("./input.txt")

	scanner := bufio.NewScanner(f)

	// regexps!!
	firstNum, _ := regexp.Compile("[0-9]")
	lastNum, _ := regexp.Compile(".*([0-9])")
	s := 0

	for scanner.Scan() {
		line := scanner.Text()
		//fmt.Println(line)
		firstMatch := firstNum.FindString(line)
		lastMatch := lastNum.FindStringSubmatch(line)
		num := firstMatch + lastMatch[1]
		//fmt.Println(num)
		//fmt.Println("")
		numStr, _ := strconv.Atoi(num)
		s += numStr
	}
	fmt.Println(s)
}
