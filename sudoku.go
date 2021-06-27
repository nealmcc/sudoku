package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"mcconachie.co/sudoku/models"
)

func main() {
	start := time.Now()

	infile := os.Args[1]
	f, err := os.Open(infile)
	if err != nil {
		log.Fatal(err)
	}

	s := bufio.NewScanner(f)

	// discard the first line (number of puzzles)
	s.Scan()

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

	duration := time.Since(start)
	fmt.Printf("solved %d sudokus in %s\n", i, duration)
}

// solve recursively solves a sudoku grid, returning true when it is solved,
// along with the number of times we had to backtrack.
func solve(g *models.Grid) (bool, int) {
	if err := g.Normalize(); err != nil {
		return false, 0
	}
	ix, done := findNextEmptyCell(g)
	if done {
		return true, 0
	}

	backtracks := 0
	for k := 1; k <= 9; k++ {
		if !g.CanSet(ix, k) {
			continue
		}
		snapshot := g.Clone()
		g.Set(ix, k)
		done, b := solve(g)
		backtracks += b
		if done {
			return true, backtracks
		}
		*g = *snapshot
		backtracks++
	}

	return false, backtracks
}

func findNextEmptyCell(g *models.Grid) (int, bool) {
	ix := 0
	for ix < 81 {
		sq := g.Get(ix)
		if !sq.IsDefined() {
			return ix, false
		}
		ix++
	}
	return 0, true
}
