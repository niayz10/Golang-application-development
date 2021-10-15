package main

import (
	"fmt"
	"testing"
)

func TestValidityOfInput(t *testing.T) {
	input := 1.0

	actual, err := ValidityOfInput(input)
	
	if err !=nil {
		t.Errorf("expected number betwen 1 and 5, but got %v", err)
	}

	if actual > 0 || actual <6 {
		fmt.Printf("Ok")
	}
}
