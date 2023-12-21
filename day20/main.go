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
		// fmt.Printf("Sending pulse:  %v %v -> %v\n", d, pulseIsHigh, d)
		pulses = append(pulses, Pulse{source: m.name, dest: d, isHigh: pulseIsHigh})
	}
	return pulses
}

func process(p Pulse, modules ModuleMap) []Pulse {
	m := modules[p.dest]
	pulseIsHigh := p.isHigh
	pulseSource := p.dest

	var pulses []Pulse
	if m.op == "%" {
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
	f, _ := os.Open("./baby-input.txt")

	scanner := bufio.NewScanner(f)

	// Parse the modules
	modules := ModuleMap{}
	for scanner.Scan() {
		line := scanner.Text()
		matcher := regexp.MustCompile(`(.)(.*) -> (.+)`)
		fields := matcher.FindStringSubmatch(line)
		m := Module{op: fields[1], name: fields[2], dests: strings.Split(fields[3], ", "), conjMap: map[string]bool{}}
		modules[m.name] = &m
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
