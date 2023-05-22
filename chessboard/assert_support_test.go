package ai

import (
	"fmt"
	"testing"
)

func expectInt(t *testing.T, k string, s int, s2 int) {
	// fmt.Printf("Compare:%s:%s\n", s , s2)
	if s != s2 {
		t.Errorf(fmt.Sprintf("Expected %d to be equal to %d   key:%s", s, s2, k))
	}
}

func expectBool(t *testing.T, k string, s bool, s2 bool) {
	// fmt.Printf("Compare:%s:%s\n", s , s2)
	if s != s2 {
		t.Errorf(fmt.Sprintf("Expected %t to be equal to %t   key:%s", s, s2, k))
	}
}

func expect(t *testing.T, s string, s2 string) {
	fmt.Printf("Compare:%s:%s\n", s, s2)
	if s != s2 {
		t.Errorf(fmt.Sprintf("Expected %s to be equal to %s", s, s2))
	}
}
