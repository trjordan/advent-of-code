package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func update(cur int, chr int) int {
	return ((cur + chr) * 17) % 256
}

func main() {
	f, _ := os.Open("./input.txt")

	scanner := bufio.NewScanner(f)

	// Just one line today!
	steps := []string{}
	for scanner.Scan() {
		steps = strings.Split(scanner.Text(), ",")
	}

	total := 0
	for _, step := range steps {
		cur := 0
		for i := 0; i < len(step); i++ {
			cur = update(cur, int(step[i]))
		}
		fmt.Println(step, cur)
		total += cur
	}

	fmt.Println("total", total)

}
