package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type ModuleMap map[string]*Module

type Module struct {
	// Declared stuff
	op    string
	name  string
	dests []string

	// Internal state
	isOn    bool
	conjMap map[string]bool // predecessor label -> wasLow

	// For counting
	cycleLength int
	offset      int
}

type Pulse struct {
	source string
	dest   string
	isHigh bool
}

func (m Module) sendPulses(pulseIsHigh bool, modules ModuleMap) []Pulse {
	pulses := []Pulse{}
	for _, d := range m.dests {
		// fmt.Printf("Creating pulse:  %v %v -> %v\n", m.name, pulseIsHigh, d)
		pulses = append(pulses, Pulse{source: m.name, dest: d, isHigh: pulseIsHigh})
	}
	return pulses
}

func process(p Pulse, modules ModuleMap) []Pulse {
	m, exists := modules[p.dest]
	pulseIsHigh := p.isHigh
	pulseSource := p.source

	var pulses []Pulse
	if !exists {
		if p.dest == "rx" && !p.isHigh {
			fmt.Println("got a low on rx")
			os.Exit(0)
		}
	} else if m.op == "%" {
		if pulseIsHigh {
			// Ignored
		} else {
			m.isOn = !m.isOn
			if m.isOn {
				// If it was off, it turns on and sends a high pulse.
				// fmt.Println("% is on, sending pulses")
				pulses = m.sendPulses(true, modules)
			} else {
				// fmt.Println("% is off, sending pulses")
				// If it was on, it turns off and sends a low pulse.
				pulses = m.sendPulses(false, modules)
			}
		}
	} else if m.op == "&" {
		m.conjMap[pulseSource] = pulseIsHigh
		allHigh := true
		for _, isHigh := range m.conjMap {
			if !isHigh {
				allHigh = false
				break
			}
		}
		// fmt.Println("& sending pulses, are all high?", allHigh)
		pulses = m.sendPulses(!allHigh, modules)
	} else if m.op == "b" {
		// fmt.Println("broadcasting")
		pulses = m.sendPulses(false, modules)
	}
	return pulses
}

func printModuleTree(root *Module, modules ModuleMap, depth int, seen map[string]bool) {
	fmt.Printf("%v%v%v\n", strings.Repeat("  ", depth), root.op, root.name)
	seen[root.name] = true
	for _, childName := range root.dests {
		child, err := modules[childName]
		if !err {
			fmt.Printf("%v%v - dummy\n", strings.Repeat("  ", depth+1), childName)
			continue
		}
		if !seen[child.name] {
			printModuleTree(child, modules, depth+1, seen)
		} else {
			fmt.Printf("%v%v - recurse\n", strings.Repeat("  ", depth+1), childName)

		}
	}
}

func main() {
	f, _ := os.Open("./input.txt")

	scanner := bufio.NewScanner(f)

	// Parse the modules
	modules := ModuleMap{}
	conjMaps := map[string][]string{}
	for scanner.Scan() {
		line := scanner.Text()
		matcher := regexp.MustCompile(`(.)(.*) -> (.+)`)
		fields := matcher.FindStringSubmatch(line)
		name := fields[2]
		dests := strings.Split(fields[3], ", ")
		m := Module{op: fields[1], name: name, dests: dests, conjMap: map[string]bool{}}
		// Capture who points where
		for _, d := range dests {
			conjMaps[d] = append(conjMaps[d], name)
		}
		modules[m.name] = &m
	}

	// Fix up all the conj modules
	for _, m := range modules {
		cm := map[string]bool{}
		for _, n := range conjMaps[m.name] {
			cm[n] = false
		}
		m.conjMap = cm
	}

	for _, m := range modules {
		fmt.Println(*m)
	}

	// Print so I can dump it into graphviz
	// https://dreampuf.github.io/GraphvizOnline/
	//
	// Replace their list of nodes (line 17) with this output to see your graph!
	for _, m := range modules {
		for _, dName := range m.dests {
			d := modules[dName]
			var mOp, dOp string

			if m.op == "&" {
				mOp = "CON"
			} else {
				mOp = "FLF"
			}
			if dName == "rx" {
				fmt.Printf(`%v%v -> rx`, mOp, m.name)
				fmt.Println()
				continue
			}
			if d.op == "&" {
				dOp = "CON"
			} else {
				dOp = "FLF"
			}
			fmt.Printf(`%v%v -> %v%v`, mOp, m.name, dOp, d.name)
			fmt.Println()
		}
	}

	// OK so there's only a handful of nodes that feed the single node that
	// feeds rx rx predecessors: pg, sp, sv, qs
	//
	// Let those run, figure out their cycle length and offset. My input had 4
	// and the cycles all emitted a low on the last state (presumably all inputs
	// were like this). Here we do a bunch of processing in order to grab those.
	preOffsets := map[string]int{}
	preCycles := map[string]int{}
	for i := 1; i < 10000; i++ {

		// Send a broadcast pulse!!
		// fmt.Println("nth push", i)
		pulseQueue := []Pulse{{source: "button", dest: "roadcaster", isHigh: false}}
		for len(pulseQueue) > 0 {
			p := pulseQueue[0]
			pulseQueue = pulseQueue[1:]
			// fmt.Printf("Processing pulse:  %v %v -> %v\n", p.source, p.isHigh, p.dest)
			newPulses := process(p, modules)
			pulseQueue = append(pulseQueue, newPulses...)

			// Let's capture their cycle length!
			if (p.dest == "pg" || p.dest == "sp" || p.dest == "sv" || p.dest == "qs") && !p.isHigh {
				fmt.Println("YAY", p.dest, i)
				if preOffsets[p.dest] == 0 {
					preOffsets[p.dest] = i
				} else if preCycles[p.dest] == 0 {
					preCycles[p.dest] = i - preOffsets[p.dest]
				}
			}
		}
	}

	// What's an LCM algorithm, I plugged the 4 offsets into a calculator
	// https://www.calculatorsoup.com/calculators/math/lcm.php
	//
	// You'd have to do something different if the offsets weren't the same as
	// the cycle, but I presume that the "low is the last thing in the cycle" is
	// part of the puzzle's design
	fmt.Println(preOffsets)
	fmt.Println(preCycles)

}
