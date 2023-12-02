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
		maxRed := 0
		maxGreen := 0
		maxBlue := 0
		for _, pull := range res {
			if pull["red"] > maxRed {
				maxRed = pull["red"]
			}
			if pull["green"] > maxGreen {
				maxGreen = pull["green"]
			}
			if pull["blue"] > maxBlue {
				maxBlue = pull["blue"]
			}
		}
		power := maxRed * maxGreen * maxBlue
		fmt.Println("line result", maxRed, maxGreen, maxBlue, power, gameId)
		s = s + power

	}

	fmt.Println("result", s)
}
