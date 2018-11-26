package hash

import (
	"fmt"
	"testing"
)

func TestBKDRHash(t *testing.T) {
	s1 := "hello"
	s2 := "world"

	fmt.Println(BKDRHash(s1))
	fmt.Println(BKDRHash(s2))
}
