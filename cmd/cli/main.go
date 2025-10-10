package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"advent2024/pkg/solver"

	_ "advent2024/pkg/d1"
	_ "advent2024/pkg/d2"
	_ "advent2024/pkg/d3"
	_ "advent2024/pkg/d4"
	_ "advent2024/pkg/d5"
	_ "advent2024/pkg/d6"
	_ "advent2024/pkg/d7"
)

func main() {
	filename := flag.String("filename", "", "Specify filename with puzzle input")
	part := flag.Int("part", 1, "Specify which puzzle part to run")
	day := flag.String("day", "d1", "Specify which day to run")

	flag.Parse()

	solver, ok := solver.New(*day)

	if !ok {
		log.Fatal("Unable to find solver for day ", *day)
	}

	fh, err := os.Open(*filename)

	if err != nil {
		log.Fatal(err)
	}

	defer fh.Close()

	err = solver.Init(fh)

	if err != nil {
		log.Fatal(err)
	}

	result, err := solver.Solve(*part)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Result - Part %d: %s", *part, result)

	// bump test 3
}
