package main

import (
	"fmt"
	"github.com/GerryLon/go-toolkit/argparser"
)

func main() {
	f := argparser.Option("f", "force").Default(false).Required().Bool()

	conf := argparser.Option("conf", "config file").Required().String()
	argparser.Parse()

	fmt.Println(f.Value(), conf.Value())

	// g := flag.String("g", "", "g g")
	// flag.Parse()
	// fmt.Printf("g=%s\n", *g)
}
