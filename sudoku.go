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
		// unsolved := sudoku.String()
		// _, guesses, backtracks := solve(&sudoku, 0)
		solve(&sudoku, 0)
		// fmt.Printf("solved puzzle %d with %d guesses and %d backtracks:\n", i, guesses, backtracks)
		// fmt.Println(unsolved)
		// fmt.Println(sudoku)
		// fmt.Println()
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
func solve(g *models.Grid, startAt int) (bool, int, int) {
	g.Reduce()
	ix, done := findNextEmptyCell(g, startAt)
	if done {
		return true, 0, 0
	}

	var (
		guesses    int
		backtracks int
		snapshots  stack = stack{}
	)

	for k := 1; k < 10; k++ {
		if g.CanSet(ix, k) {
			guesses++
			snapshots.push(g.Clone())
			g.Set(ix, k)
			done, a, b := solve(g, ix)
			guesses += a
			backtracks += b
			if done {
				return true, guesses, backtracks
			}
			g, _ = snapshots.pop()
			backtracks++
		}
	}

	return false, guesses, backtracks
}

func findNextEmptyCell(g *models.Grid, startAt int) (int, bool) {
	ix := startAt
	for ix < 81 {
		sq := g.Get(ix)
		if !sq.IsDefined() {
			return ix, false
		}
		ix++
	}
	return 0, true
}

type stack struct {
	v []*models.Grid
}

func (s stack) push(g *models.Grid) {
	s.v = append(s.v, g)
}

func (s stack) pop() (*models.Grid, bool) {
	if len(s.v) == 0 {
		return nil, false
	}
	last := len(s.v) - 1
	top := s.v[last]
	s.v = s.v[:last]
	return top, true
}
