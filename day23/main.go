package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type Point struct {
	x, y int
}

type Segment struct {
	start, end Point
	points     []Point
	length     int
}

func NextPoints(strMap []string, prev Point, cur Point) []Point {
	nextPoints := []Point{}

	// Look right
	if cur.x > 0 {
		if !(cur.x+1 == prev.x && cur.y == prev.y) &&
			(strMap[cur.y][cur.x+1] == '.' ||
				strMap[cur.y][cur.x+1] == '>') {
			nextPoints = append(nextPoints, Point{cur.x + 1, cur.y})
		}
	}

	// Look left
	if cur.x < len(strMap[0])-1 {
		if !(cur.x-1 == prev.x && cur.y == prev.y) &&
			(strMap[cur.y][cur.x-1] == '.' ||
				strMap[cur.y][cur.x-1] == '<') {
			nextPoints = append(nextPoints, Point{cur.x - 1, cur.y})
		}
	}

	// Look up
	if cur.y > 0 {
		if !(cur.x == prev.x && cur.y-1 == prev.y) &&
			(strMap[cur.y-1][cur.x] == '.' ||
				strMap[cur.y-1][cur.x] == '^') {
			nextPoints = append(nextPoints, Point{cur.x, cur.y - 1})
		}
	}

	// Look down
	if cur.y < len(strMap)-1 {
		if !(cur.x == prev.x && cur.y+1 == prev.y) &&
			(strMap[cur.y+1][cur.x] == '.' ||
				strMap[cur.y+1][cur.x] == 'v') {
			nextPoints = append(nextPoints, Point{cur.x, cur.y + 1})
		}
	}

	return nextPoints
}

func CreateSegments(strMap []string) []Segment {
	// Manual start
	nextStarts := []Point{{1, 0}, {1, 1}}
	prev := nextStarts[0]
	cur := Point{1, 1}
	segments := []Segment{}
	curPoints := []Point{}

	// Don't get stuck looping, track the seen and where they connect to
	seen := map[Point]bool{}

	// Go explore!
	for len(nextStarts) > 0 {
		curPoints = append(curPoints, prev)
		nexts := NextPoints(strMap, prev, cur)
		if len(nexts) != 1 {
			// wrap it up
			curPoints = append(curPoints, cur)
			segments = append(segments, Segment{
				start:  cur,
				end:    nextStarts[0],
				points: curPoints,
				length: len(curPoints) - 1,
			})
			// Also its reverse
			segments = append(segments, Segment{
				start:  nextStarts[0],
				end:    cur,
				points: curPoints,
				length: len(curPoints) - 1,
			})
			curPoints = []Point{}
			nextStarts = nextStarts[2:]

			if len(nexts) > 1 {
				if !seen[cur] {
					for i := 0; i < len(nexts); i++ {
						nextStarts = append(nextStarts, cur, nexts[i])
					}
				}
				seen[cur] = true
			}
			if len(nextStarts) > 0 {
				prev = nextStarts[0]
				cur = nextStarts[1]
			}
		} else {
			// keep chugging
			prev = cur
			cur = nexts[0]
		}
	}

	// fmt.Println(seen)
	// fmt.Println(segments)

	// [gross] dedup
	uniqueSegmentMap := map[[2]Point]Segment{}
	for i := 0; i < len(segments); i++ {
		k := [2]Point{segments[i].start, segments[i].end}
		uniqueSegmentMap[k] = segments[i]
	}
	uniqueSegments := []Segment{}
	for _, s := range uniqueSegmentMap {
		uniqueSegments = append(uniqueSegments, s)
	}
	sort.Slice(uniqueSegments, func(i, j int) bool {
		if uniqueSegments[i].start.x == uniqueSegments[j].start.x {
			return uniqueSegments[i].start.y < uniqueSegments[j].start.y
		} else {
			return uniqueSegments[i].start.x < uniqueSegments[j].start.x
		}
	})

	return uniqueSegments
}

func FindSegments(p Point, segments []Segment, seen []Point) []Segment {
	unexploredSegments := []Segment{}
	seenMap := map[Point]bool{}
	for i := 0; i < len(seen)-1; i++ {
		seenMap[seen[i]] = true
	}

	for i := 0; i < len(segments); i++ {
		// Walk this one if it starts from our start point AND we haven't seen
		// its end
		if segments[i].start == p && !seenMap[segments[i].end] {
			unexploredSegments = append(unexploredSegments, segments[i])
		}
	}

	return unexploredSegments
}

func FindLongestPathLen(start Point, lenSoFar int, segments []Segment, strMap []string, seen []Point) (int, []Point) {
	// Update this seen for this round
	newSeen := make([]Point, len(seen))
	copy(newSeen, seen)
	newSeen = append(newSeen, start)
	longestSeen := newSeen

	nextSegments := FindSegments(start, segments, newSeen)
	// short-circuit if we've found a path that terminates but isn't an end
	if len(nextSegments) == 0 && start.y != len(strMap)-1 {
		return 0, []Point{}
	}

	longest := lenSoFar
	for i := 0; i < len(nextSegments); i++ {
		// fmt.Println("Gonna go explore", nextSegments[i].end, longest, i, len(nextSegments))
		nextLenSoFar, nextSeenSoFar := FindLongestPathLen(
			nextSegments[i].end,
			lenSoFar+nextSegments[i].length,
			segments,
			strMap,
			newSeen)
		if nextLenSoFar > longest {
			longest = nextLenSoFar
			longestSeen = nextSeenSoFar
		}
		// if nextSegments[i].end.y == len(strMap)-1 {
		// 	fmt.Println("Found a path", nextLenSoFar, nextSeenSoFar)
		// }
	}

	return longest, longestSeen
}

func PrintFromSeen(strMap []string, segments []Segment, seen []Point) {
	seenMap := map[Point]Point{}
	for i := 0; i < len(seen)-1; i++ {
		seenMap[seen[i]] = seen[i+1]
	}
	validSegments := []Segment{}
	for i := 0; i < len(segments); i++ {
		if seenMap[segments[i].start] == segments[i].end {
			validSegments = append(validSegments, segments[i])
		}
	}

	PrintMap(strMap, validSegments)
}

func PrintMap(strMap []string, segments []Segment) {
	// Invert the segments into a map for easy lookup
	walkedPoints := map[Point]bool{}
	for _, s := range segments {
		for _, p := range s.points {
			walkedPoints[p] = true
		}
	}
	for i := 0; i < len(strMap); i++ {
		for j := 0; j < len(strMap[i]); j++ {
			if walkedPoints[Point{j, i}] {
				fmt.Print("O")
			} else {
				fmt.Printf("%c", strMap[i][j])
			}
		}
		fmt.Println("")
	}
}

func main() {
	f, _ := os.Open("./input.txt")

	scanner := bufio.NewScanner(f)

	// Read it all in
	strMap := []string{}
	r := strings.NewReplacer("^", ".", "v", ".", "<", ".", ">", ".")
	for scanner.Scan() {
		strMap = append(strMap, r.Replace(scanner.Text()))
	}

	segments := CreateSegments(strMap)
	for i := 0; i < len(segments); i++ {
		fmt.Println(segments[i].start, segments[i].end)
		// PrintMap(strMap, []Segment{segments[i]})
	}

	// Find all the segment that end, walk backwards to find the longest
	longest, seen := FindLongestPathLen(Point{1, 0}, 0, segments, strMap, []Point{})
	fmt.Println("longest", longest, seen)
	PrintFromSeen(strMap, segments, seen)
	// 6685 too high

}
