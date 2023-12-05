package main

import (
	"testing"
)

func TestMappedRangeIdentity(t *testing.T) {
	md := mappingData{
		ranges: [][3]int{
			{0, 15, 2},
		},
	}
	r := mappedRange{
		[2]int{0, 5},
		[2]int{7, 15},
	}
	mr := getMappedRanges(md, r)
	if len(mr) != 2 {
		t.Fatalf("Failed to leave non-interacting ranges intact, wanted len 1, got len %v", len(mr))
	}
}

func TestMappedRangeSourceStart(t *testing.T) {
	md := mappingData{
		ranges: [][3]int{
			{0, 10, 7},
		},
	}
	r := mappedRange{
		[2]int{0, 5},
		[2]int{7, 15},
	}
	mr := getMappedRanges(md, r)
	if len(mr) != 3 {
		t.Fatalf("Failed to find SourceStart in a range to split")
	}
}

func TestMappedRangeSourceEnd(t *testing.T) {
	md := mappingData{
		ranges: [][3]int{
			{0, 6, 7},
		},
	}
	r := mappedRange{
		[2]int{0, 5},
		[2]int{7, 15},
	}
	mr := getMappedRanges(md, r)
	if len(mr) != 3 {
		t.Fatalf("Failed to find SourceStart+length in a range to split")
	}
}

func TestMappedRangeFilledSource(t *testing.T) {
	md := mappingData{
		ranges: [][3]int{
			{0, 8, 2},
		},
	}
	r := mappedRange{
		[2]int{0, 5},
		[2]int{7, 15},
	}
	mr := getMappedRanges(md, r)
	if len(mr) != 4 {
		t.Fatalf("Failed to find Source fully in a range to split")
	}
}
