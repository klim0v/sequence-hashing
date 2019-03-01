package entity

import (
	"testing"
)

const startNumber uint64 = 111111

func TestMakeResult(t *testing.T) {
	unique = map[uint64]struct{}{}
	valid := true
	for i := 0; i < MaxCount; i++ {
		result := MakeResult(startNumber)
		if result == nil {
			valid = false
			break
		}
	}
	if !valid {
		t.Error("repeated")
	}
	unique = map[uint64]struct{}{}
}

func TestMakeResultMoreMaxVariants(t *testing.T) {
	unique = map[uint64]struct{}{}
	valid := true
	for i := 0; i < MaxCount+1; i++ {
		result := MakeResult(startNumber)
		if result == nil {
			valid = false
			break
		}
	}
	if valid {
		t.Error("not repeated")
	}
	unique = map[uint64]struct{}{}
}

func TestMakeResultBinaries(t *testing.T) {
	unique = map[uint64]struct{}{}
	result := MakeResult(startNumber)
	bytes, _ := result.MarshalBinary()
	unmResult := new(Result)
	_ = unmResult.UnmarshalBinary(bytes)

	if result.Number != unmResult.Number {
		t.Error("not equal")
	}
	unique = map[uint64]struct{}{}
}
