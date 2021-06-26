package models

import (
	"fmt"
	"strings"
)

type Grid struct {
	squares *[81]Square
}

// NewGrid initializes a sudoku grid using the given string.
// Each square should be given as a digit (if the square is defined).
// If the square is undefined, it should be given as either '0' or '.'
// All other characters are ignored.
func NewGrid(start string) Grid {
	g := Grid{
		squares: new([81]Square),
	}
	i := 0
	for _, ch := range []byte(start) {
		switch ch {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			g.squares[i] = NewSquare(int(ch - '0'))
			i++
		case '.':
			g.squares[i] = NewSquare(0)
			i++
		default:
			continue
		}
	}

	return g
}

func (g Grid) Squares() string {
	var b strings.Builder
	for i := 0; i < 81; i += 9 {
		row := g.squares[i : i+9]
		for _, sq := range row {
			fmt.Fprintf(&b, "%v\t", sq.Values())
		}
		b.WriteByte('\n')
	}
	return b.String()
}

// String implements the fmt.Stringer interface
func (g Grid) String() string {
	var b strings.Builder
	b.WriteString("╔═══╤═══╤═══╗\n")
	i := 0
	for x := 0; x < 3; x++ {
		writeRow(&b, g.squares[i:i+9])
		b.WriteByte('\n')
		i += 9
	}
	b.WriteString("╟───┼───┼───╢\n")
	for x := 0; x < 3; x++ {
		writeRow(&b, g.squares[i:i+9])
		b.WriteByte('\n')
		i += 9
	}
	b.WriteString("╟───┼───┼───╢\n")
	for x := 0; x < 3; x++ {
		writeRow(&b, g.squares[i:i+9])
		b.WriteByte('\n')
		i += 9
	}
	b.WriteString("╚═══╧═══╧═══╝")
	return b.String()
}

func writeRow(b *strings.Builder, row []Square) {
	b.WriteRune('║')
	for i := 0; i < 3; i++ {
		b.WriteByte(row[i].Display())
	}
	b.WriteRune('│')
	for i := 3; i < 6; i++ {
		b.WriteByte(row[i].Display())
	}
	b.WriteRune('│')
	for i := 6; i < 9; i++ {
		b.WriteByte(row[i].Display())
	}
	b.WriteRune('║')
}

// RefineByEclusion performs one round of refinement based on the process of
// exclusion.  That is, the process of excluding from each square the values
// that are definitely assigned within the same row, column or block.
// Returns a list of squares that are now defined that weren't before.
func (g Grid) RefineByExclusion() []int {
	newlyDefined := make([]int, 0, 8)
	for i := 0; i < 81; i++ {
		if g.squares[i].IsDefined() {
			continue
		}
		g.refineOne(i)
		if g.squares[i].IsDefined() {
			newlyDefined = append(newlyDefined, i)
		}
	}
	return newlyDefined
}

// refineOne the nth square for this grid, by excluding
// values which would conflict with others in the same
// row, column, or 3x3 block.
// returns true if the square has been updated
func (g Grid) refineOne(n int) bool {
	var (
		start = g.squares[n]
		row   = g.getRowSquares(n)
		col   = g.getColSquares(n)
		block = g.getBlockSquares(n)
	)
	sq := start.
		ExcludeDefined(row).
		ExcludeDefined(col).
		ExcludeDefined(block)
	g.squares[n] = sq
	return sq != start
}

// getOtherRowSquares is used to get the other squares from the same row
// as the square at n. Excludes square n from the list.
func (g Grid) getRowSquares(n int) []Square {
	row := make([]Square, 8)
	j := 0
	start := (n / 9) * 9
	for i := start; i < start+9; i++ {
		if i == n {
			continue
		}
		row[j] = g.squares[i]
		j++
	}
	return row
}

// getOtherColSquares is used to get the other squares from the same column
// as the square at n. Excludes square n from the list.
func (g Grid) getColSquares(n int) []Square {
	col := make([]Square, 8)
	j := 0
	start := n % 9
	for i := start; i < 81; i += 9 {
		if i == n {
			continue
		}
		col[j] = g.squares[i]
		j++
	}
	return col
}

// getOtherBlockSquares is used to get the other squares from the same
// 3x3 block as the square at n. Excludes square n from the list.
func (g Grid) getBlockSquares(n int) []Square {
	// r0, c0 is the top left corner of this 3x3 block:
	r0, c0 := (n/27)*27, ((n%9)/3)*3
	block := make([]Square, 8)
	b := 0
	for c := 0; c < 3; c++ {
		for r := 0; r < 3; r++ {
			ix := r0 + r*9 + c0 + c
			if ix == n {
				continue
			}
			block[b] = g.squares[ix]
			b++
		}
	}
	return block
}
