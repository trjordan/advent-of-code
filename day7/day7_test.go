package main

import (
	"testing"
)

func TestMatchingThree(t *testing.T) {
	low := "T5595"
	high := "QQQ9A"
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
	low := "A555A"
	high := "QQQQ2"
	ret := compareHandStrings(low, high)
	if ret >= 0 {
		t.Fatalf("%v is < %v, but compareHand returned %v", low, high, ret)
	}
}

func TestMatchingThreeJoker(t *testing.T) {
	low := "J5592"
	high := "77J9A"
	ret := compareHandStrings(low, high)
	if ret >= 0 {
		t.Fatalf("%v is < %v, but compareHand returned %v", low, high, ret)
	}
}

func TestFiveToFourJoker(t *testing.T) {
	low := "A5555"
	high := "QQQJQ"
	ret := compareHandStrings(low, high)
	if ret >= 0 {
		t.Fatalf("%v is < %v, but compareHand returned %v", low, high, ret)
	}
}

func TestFourToFullHouseJoker(t *testing.T) {
	low := "22QQQ"
	high := "AQQJJ"
	ret := compareHandStrings(low, high)
	if ret >= 0 {
		t.Fatalf("%v is < %v, but compareHand returned %v", low, high, ret)
	}
}

func TestValueFiveWithJoker(t *testing.T) {
	s := "J5555"
	ret := findTypeValue(s)
	if ret != 6 {
		t.Fatalf("%v is type 6 but compareHand returned %v", s, ret)
	}
}

func TestValueFourWithJoker(t *testing.T) {
	s := "J55A5"
	ret := findTypeValue(s)
	if ret != 5 {
		t.Fatalf("%v is type 5 but compareHand returned %v", s, ret)
	}
}

func TestValueFourWithJokers(t *testing.T) {
	s := "J55JA"
	ret := findTypeValue(s)
	if ret != 5 {
		t.Fatalf("%v is type 5 but compareHand returned %v", s, ret)
	}
}

func TestValueFullHouseWithJoker(t *testing.T) {
	s := "J5445"
	ret := findTypeValue(s)
	if ret != 4 {
		t.Fatalf("%v is type 4 but compareHand returned %v", s, ret)
	}
}

func TestValueThreeWithJoker(t *testing.T) {
	s := "1233J"
	ret := findTypeValue(s)
	if ret != 3 {
		t.Fatalf("%v is type 3 but compareHand returned %v", s, ret)
	}
}

func TestValuePairWithJoker(t *testing.T) {
	s := "1234J"
	ret := findTypeValue(s)
	if ret != 1 {
		t.Fatalf("%v is type 1 but compareHand returned %v", s, ret)
	}
}
