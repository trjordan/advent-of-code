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
		fmt.Println("dummy node")
	} else if m.op == "%" {
		if pulseIsHigh {
			// Ignored
		} else {
			m.isOn = !m.isOn
			if m.isOn {
				// If it was off, it turns on and sends a high pulse.
				fmt.Println("% is on, sending pulses")
				pulses = m.sendPulses(true, modules)
			} else {
				fmt.Println("% is off, sending pulses")
				// If it was on, it turns off and sends a low pulse.
				pulses = m.sendPulses(false, modules)
			}
		}
	} else if m.op == "&" {
		m.conjMap[pulseSource] = pulseIsHigh
		fmt.Println("conjmap", m.conjMap)
		allHigh := true
		for _, isHigh := range m.conjMap {
			if !isHigh {
				allHigh = false
				break
			}
		}
		fmt.Println("& sending pulses, are all high?", allHigh)
		pulses = m.sendPulses(!allHigh, modules)
	} else if m.op == "b" {
		fmt.Println("broadcasting")
		pulses = m.sendPulses(false, modules)
	}
	return pulses
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

	var low, high int
	for i := 0; i < 1000; i++ {

		// Send a broadcast pulse!!
		pulseQueue := []Pulse{{source: "button", dest: "roadcaster", isHigh: false}}
		for len(pulseQueue) > 0 {
			p := pulseQueue[0]
			pulseQueue = pulseQueue[1:]
			fmt.Printf("Processing pulse:  %v %v -> %v\n", p.source, p.isHigh, p.dest)
			newPulses := process(p, modules)
			pulseQueue = append(pulseQueue, newPulses...)
			if p.isHigh {
				high += 1
			} else {
				low += 1
			}
		}
	}
	fmt.Println("total", low, high, low*high)
}
