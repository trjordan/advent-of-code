package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type tile struct {
	row int
	col int
}

type node struct {
	mirror    string
	tile      tile
	dir       string  // The way the light travels, u/d/l/r
	energized []tile  // All the subsequent non-interacting tiles
	next      []*node // All the immediate interacting tiles (could)
}

func isValidTile(t tile, grid []string) bool {
	// isInvalid just seems easier to read, let me live
	isInvalid := t.row < 0 || t.row >= len(grid) || t.col < 0 || t.col >= len(grid[0])
	return !isInvalid
}

func isValidNode(n node, grid []string) bool {
	// Check if the tile in this node actually exists in the grid
	return isValidTile(n.tile, grid)
}

// Unsafe, will return tiles out of range
func incrUnsafe(t tile, dir string) tile {
	nextTile := tile{row: t.row, col: t.col}
	if dir == "u" {
		nextTile.row -= 1
	} else if dir == "d" {
		nextTile.row += 1
	} else if dir == "r" {
		nextTile.col += 1
	} else if dir == "l" {
		nextTile.col -= 1
	}
	return nextTile
}

// Bounce off a mirror, but will return nodes off the grid
func bounceUnsafe(t tile, nextChr string, dir string) []node {
	var nextNodes []node
	if nextChr == "|" && strings.Contains("lr", dir) {
		fmt.Println("but uhhhh")
		nextNodes = []node{
			{mirror: "|", tile: tile{row: t.row - 1, col: t.col}, dir: "u"},
			{mirror: "|", tile: tile{row: t.row + 1, col: t.col}, dir: "d"},
		}
	} else if nextChr == "-" && strings.Contains("ud", dir) {
		nextNodes = []node{
			{mirror: "-", tile: tile{row: t.row, col: t.col + 1}, dir: "r"},
			{mirror: "-", tile: tile{row: t.row, col: t.col - 1}, dir: "l"},
		}
	} else if nextChr == `\` {
		if dir == "r" {
			nextNodes = []node{{mirror: `\`, tile: tile{row: t.row + 1, col: t.col}, dir: "d"}}
		} else if dir == "l" {
			nextNodes = []node{{mirror: `\`, tile: tile{row: t.row - 1, col: t.col}, dir: "u"}}
		} else if dir == "u" {
			nextNodes = []node{{mirror: `\`, tile: tile{row: t.row, col: t.col - 1}, dir: "l"}}
		} else if dir == "d" {
			nextNodes = []node{{mirror: `\`, tile: tile{row: t.row, col: t.col + 1}, dir: "r"}}
		}
	} else if nextChr == `/` {
		if dir == "r" {
			nextNodes = []node{{mirror: `/`, tile: tile{row: t.row - 1, col: t.col}, dir: "u"}}
		} else if dir == "l" {
			nextNodes = []node{{mirror: `/`, tile: tile{row: t.row + 1, col: t.col}, dir: "d"}}
		} else if dir == "u" {
			nextNodes = []node{{mirror: `/`, tile: tile{row: t.row, col: t.col + 1}, dir: "r"}}
		} else if dir == "d" {
			nextNodes = []node{{mirror: `/`, tile: tile{row: t.row, col: t.col - 1}, dir: "l"}}
		}
	}
	return nextNodes
}

func updateNextNode(n *node, grid []string) {
	fmt.Println("starting on", n.tile, n.dir)

	for nextTile := n.tile; ; nextTile = incrUnsafe(nextTile, n.dir) {
		fmt.Println("working on", nextTile)
		if !isValidTile(nextTile, grid) {
			// Whoops, off the edge! Good enough
			fmt.Println("bye bye")
			break
		}
		n.energized = append(n.energized, nextTile)
		nextChr := string(grid[nextTile.row][nextTile.col])

		if nextChr != "." &&
			!(nextChr == "-" && strings.Contains("lr", n.dir)) &&
			!(nextChr == "|" && strings.Contains("ud", n.dir)) {
			// Found a reason to move on!
			fmt.Println("bouncing at", nextTile)
			nextNodes := bounceUnsafe(nextTile, nextChr, n.dir)
			fmt.Println("found", len(nextNodes))
			for i := 0; i < len(nextNodes); i++ {
				if isValidNode(nextNodes[i], grid) {
					fmt.Println("Adding", &nextNodes[i].tile)
					n.next = append(n.next, &nextNodes[i])
				}
			}
			break
		}

	}
	fmt.Println("done updating", n)
}

func getKey(n node) string {
	return fmt.Sprintf("%v,%v,%v", n.tile.row, n.tile.col, n.dir)
}

func main() {
	f, _ := os.Open("./input.txt")

	scanner := bufio.NewScanner(f)
	grid := []string{}
	for scanner.Scan() {
		grid = append(grid, scanner.Text())
	}

	// Let's make a graph!
	allNodes := []node{{
		mirror: "",
		dir:    "r",
		tile:   tile{row: 0, col: 0},
	}}
	seen := map[string]bool{}      // key is row,col,dir
	energized := map[string]bool{} // key is row,col. ALso, TODO, make a proper graph traversal!!
	for len(allNodes) > 0 {
		n := allNodes[0]
		allNodes = allNodes[1:]
		// Fill out energized and next fields
		updateNextNode(&n, grid)
		// Mark it as seen to prevent cycles
		seen[getKey(n)] = true
		// Queue the next nodes (breadth-first)
		for i := 0; i < len(n.next); i++ {
			fmt.Println("considering adding", n.next[i], getKey(*n.next[i]))
			if !seen[getKey(*n.next[i])] {
				fmt.Println("truly added", n.next[i], getKey(*n.next[i]))
				allNodes = append(allNodes, *n.next[i])
			}
		}
		fmt.Println(n)
		fmt.Println("remaining nodes", len(allNodes))
		fmt.Println("--------------main loop done---------------")
		fmt.Println("")

		// Lazycakes
		for _, e := range n.energized {
			energized[fmt.Sprintf("%v,%v", e.row, e.col)] = true
		}
	}

	// pretty print
	for i := 0; i < len(grid); i++ {
		fmt.Printf("%03d ", i)
		for j := 0; j < len(grid[0]); j++ {
			if energized[fmt.Sprintf("%v,%v", i, j)] {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Print("\n")
	}
	fmt.Println("total", len(energized))

}
