package main

import (
	"testing"
)

func TestMatchingThreePair(t *testing.T) {
	low := "T55J5"
	high := "QQQJA"
	ret := compareHandStrings(low, high)
	if ret >= 0 {
		t.Fatalf("%v is < %v, but compareHand returned %v", low, high, ret)
	}
}

func TestFiveToFour(t *testing.T) {
	low := "T5555"
	high := "QQQQQ"
	ret := compareHandStrings(low, high)
	if ret >= 0 {
		t.Fatalf("%v is < %v, but compareHand returned %v", low, high, ret)
	}
}

func TestFourToFullHouse(t *testing.T) {
	low := "A5555"
	high := "QQQQ2"
	ret := compareHandStrings(low, high)
	if ret >= 0 {
		t.Fatalf("%v is < %v, but compareHand returned %v", low, high, ret)
	}
}
