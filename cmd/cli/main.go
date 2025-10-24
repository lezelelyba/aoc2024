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
	_ "advent2024/pkg/d8"
	_ "advent2024/pkg/d9"
)

var Version string = "dev"

func main() {

	filename := flag.String("filename", "", "Specify filename with puzzle input")
	part := flag.Int("part", 1, "Specify which puzzle part to run")
	day := flag.String("day", "d1", "Specify which day to run")
	version := flag.Bool("version", false, "List version")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Version %s\n\n", Version)
		flag.PrintDefaults()
	}
	flag.Parse()

	if *version == true {
		flag.Usage()
		os.Exit(0)
	}

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

	//bumptest6
}
