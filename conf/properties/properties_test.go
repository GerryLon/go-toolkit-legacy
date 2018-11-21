package properties

import (
	"fmt"
	"testing"
)

func TestProperties_Value(t *testing.T) {
	p := Properties{
		Filename: "demo.properties",
	}

	v := p.MustValue("a")

	if v != "1" {
		t.Errorf(`p.MustValue("a"), err`)
	}

	v, ok := p.Value("b")
	fmt.Println(v, ok)

	fmt.Printf("p.All: %v", p.MustAll())
}
