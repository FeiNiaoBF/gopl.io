// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 10.
//!+

// Dup2 prints the count and text of lines that appear more than once
// in the input.  It reads from stdin or from a list of named files.
package main

import (
	"bufio"
	"fmt"
	"os"
)

type Line struct {
	Count int
	Name  string
}

func main() {
	counts := make(map[string]Line)
	files := os.Args[1:]
	if len(files) == 0 {
		countLines(os.Stdin, counts)
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue
			}
			countLines(f, counts)
			f.Close()
		}
	}
	for val, n := range counts {
		if n.Count > 1 {
			fmt.Printf("%d\t%s\tIn file %v\n", n.Count, val, n.Name)
		}
	}
}

func countLines(f *os.File, counts map[string]Line) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		counts[input.Text()] = Line{counts[input.Text()].Count + 1, f.Name()}
	}
	// NOTE: ignoring potential errors from input.Err()
}

//!-
