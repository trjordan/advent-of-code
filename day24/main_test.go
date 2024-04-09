package main

import "testing"

func TestWillCollide(t *testing.T) {
	// Two hailstones moving apart along both x and y axes
	h1 := Hailstone{x: 0, y: 0, vx: -1, vy: -1} // Hailstone 1
	h2 := Hailstone{x: 2, y: 2, vx: 1, vy: 1}   // Hailstone 2
	if h1.WillCollide(h2) {
		t.Error("Expected false, got true")
	}

	// Two hailstones moving towards each other along both x and y axes
	h3 := Hailstone{x: 0, y: 0, vx: 1, vy: 1}  // Hailstone 1
	h4 := Hailstone{x: 0, y: 2, vx: 1, vy: -1} // Hailstone 2
	if !h3.WillCollide(h4) {
		t.Error("Expected true, got false")
	}

	// Two hailstones moving parallel but not colliding
	h5 := Hailstone{x: 0, y: 0, vx: 1, vy: 1} // Hailstone 1
	h6 := Hailstone{x: 2, y: 0, vx: 2, vy: 2} // Hailstone 2
	if h5.WillCollide(h6) {
		t.Error("Expected false, got true")
	}

	// Two hailstones moving parallel but not colliding, one stationary
	h7 := Hailstone{x: 0, y: 0, vx: 0, vy: 0}  // Hailstone 1 (stationary)
	h8 := Hailstone{x: 2, y: 0, vx: -1, vy: 1} // Hailstone 2
	if h7.WillCollide(h8) {
		t.Error("Expected false, got true")
	}

	// Two hailstones moving together with negative velocities
	h11 := Hailstone{x: 5, y: 0, vx: -1, vy: 1} // Hailstone 1
	h12 := Hailstone{x: 0, y: 2, vx: -2, vy: 1} // Hailstone 2
	if !h11.WillCollide(h12) {
		t.Error("Expected true, got false")
	}
}
