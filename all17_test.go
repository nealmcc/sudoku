// +build slow

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"testing"

	"mcconachie.co/sudoku/models"
)

func TestAll17(t *testing.T) {
	infile := "./all_17_clue_sudokus.txt"
	f, err := os.Open(infile)
	if err != nil {
		log.Fatal(err)
	}

	s := bufio.NewScanner(f)
	s.Scan() // discard count of puzzles

	i := 0
	for s.Scan() {
		if err := s.Err(); err != nil {
			log.Fatal(fmt.Errorf("error reading input: %w", err))
		}
		sudoku := models.NewGrid(s.Bytes())
		solve(&sudoku)
		i++
	}

	if err := s.Err(); err != nil {
		log.Fatal(fmt.Errorf("error reading input: %w", err))
	}
}
