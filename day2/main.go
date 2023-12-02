package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type M map[string]int

func parseLine(line string) (gameId int, ret []M) {
	// get the game id
	gameIdMatcher, _ := regexp.Compile(".*?([0-9]+):")
	gameIdMatch := gameIdMatcher.FindStringSubmatch(line)
	gameId, _ = strconv.Atoi(gameIdMatch[1])

	// parse the pulls into a list of maps
	subsetsStr := strings.Split(line, ":")[1]
	subsets := strings.Split(subsetsStr, ";")
	for _, subset := range subsets {
		pulls := strings.Split(subset, ",")
		pullColorMap := M{}
		for _, pull := range pulls {
			kv := strings.Split(strings.Trim(pull, " "), " ")
			v, _ := strconv.Atoi(kv[0])
			pullColorMap[kv[1]] = v
		}
		ret = append(ret, pullColorMap)
	}

	return
}

func main() {
	f, _ := os.Open("./input.txt")

	scanner := bufio.NewScanner(f)

	s := 0
	for scanner.Scan() {
		line := scanner.Text()
		gameId, res := parseLine(line)

		fmt.Println(line)
		// myJson, _ := json.MarshalIndent(res, "", "    ")
		// fmt.Println(string(myJson))

		// Logic!!
		valid := true
		for _, pull := range res {
			if pull["red"] > 12 {
				valid = false
			}
			if pull["green"] > 13 {
				valid = false
			}
			if pull["blue"] > 14 {
				valid = false
			}
		}
		fmt.Println("line result", valid, gameId)
		if valid {
			s = s + gameId
		}
	}

	fmt.Println("result", s)
}
