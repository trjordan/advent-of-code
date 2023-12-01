package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

func main() {
	f, _ := os.Open("./input.txt")

	scanner := bufio.NewScanner(f)

	// regexps!!
	firstNum, _ := regexp.Compile("[0-9]")
	lastNum, _ := regexp.Compile(".*([0-9])")
	s := 0

	// replacements!
	replacements := map[string]string{
		"one":   "1",
		"two":   "2",
		"three": "3",
		"four":  "4",
		"five":  "5",
		"six":   "6",
		"seven": "7",
		"eight": "8",
		"nine":  "9",
	}

	for scanner.Scan() {
		line := scanner.Text()

		// fix-up for the second part
		fmt.Println(line)
		replacedLine := ""
		for len(line) > 0 {
			replaced := false
			for word, digit := range replacements {
				if strings.HasPrefix(line, word) {
					replacedLine = replacedLine + digit
					//line = strings.Replace(line, word, "", 1)
					line = string(line[1:])
					replaced = true
					break
				}
			}
			if replaced {
				continue
			}
			r, size := utf8.DecodeRuneInString(line)
			replacedLine = replacedLine + string(r)
			line = line[size:]
		}
		fmt.Println("ok!", replacedLine)

		// do that matching
		firstMatch := firstNum.FindString(replacedLine)
		lastMatch := lastNum.FindStringSubmatch(replacedLine)
		num := firstMatch + lastMatch[1]
		fmt.Println(num)
		numStr, _ := strconv.Atoi(num)
		s = s + numStr
		fmt.Println("sum", s)
		fmt.Println("")
	}
	fmt.Println(s)
}
