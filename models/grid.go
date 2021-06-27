package models

import (
	"strings"
)

type Grid struct {
	squares *[81]Square
}

// NewGrid initializes a sudoku grid using the given input.
// Each square should be given as a digit (if the square is defined).
// If the square is undefined, it should be given as either '0' or '.'
// All other characters are ignored.
func NewGrid(in []byte) Grid {
	g := Grid{
		squares: new([81]Square),
	}
	i := 0
	for _, ch := range in {
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

// Clone creates a deep copy of this Grid
func (g Grid) Clone() *Grid {
	sq := new([81]Square)
	for i := 0; i < 81; i++ {
		sq[i] = g.squares[i]
	}
	return &Grid{sq}
}

// String implements the fmt.Stringer interface
func (g Grid) String() string {
	var b strings.Builder

	i := 0
	for x := 0; x < 3; x++ {
		writeRow(&b, g.squares[i:i+9])
		b.WriteByte('\n')
		i += 9
	}
	b.WriteByte('\n')
	for x := 0; x < 3; x++ {
		writeRow(&b, g.squares[i:i+9])
		b.WriteByte('\n')
		i += 9
	}
	b.WriteByte('\n')
	for x := 0; x < 3; x++ {
		writeRow(&b, g.squares[i:i+9])
		b.WriteByte('\n')
		i += 9
	}

	return b.String()
}

func writeRow(b *strings.Builder, row []Square) {
	for i := 0; i < 3; i++ {
		b.WriteByte(row[i].Display())
	}
	b.WriteRune(' ')
	for i := 3; i < 6; i++ {
		b.WriteByte(row[i].Display())
	}
	b.WriteRune(' ')
	for i := 6; i < 9; i++ {
		b.WriteByte(row[i].Display())
	}
}

// Get retrieves the current value of the square at index i
func (g Grid) Get(i int) Square {
	return g.squares[i]
}

// CanSet checks to see if it is legal to set square i to the value k
// This function is only valid if the grid has been Reduced, and has not
// been added to since.
func (g Grid) CanSet(i, k int) bool {
	curr, next := g.squares[i], squareEnum[k]
	return curr&next > 0
}

// Set assigns the value k to square i
// Reduce should be called after Set, to maintain the integrity of the grid
func (g Grid) Set(i, k int) {
	g.squares[i] = squareEnum[k]
}

// Reduce repeatedly applies  to the grid, identifying possible and impossible
// values for each square.
func (g Grid) Reduce() {
	for {
		delta := 0
		ids := g.reduceByNegation()
		delta += len(ids)

		ids = g.reduceByDeduction()
		delta += len(ids)

		if delta == 0 {
			break
		}
	}
}

// reduceByNegation performs one round of refinement based on the process of
// set subtraction.  That is, the process of excluding from each square the
// values that are definitely assigned within the same row, column or block.
// Returns a list of squares that are now defined that weren't before.
func (g Grid) reduceByNegation() []int {
	newlyDefined := make([]int, 0, 8)
	for i := 0; i < 81; i++ {
		if g.squares[i].IsDefined() {
			continue
		}
		if g.negateOthers(i) {
			newlyDefined = append(newlyDefined, i)
		}
	}
	return newlyDefined
}

// negateOthers refines the the nth square for this grid by excluding values
// which would conflict with others in the same row, column, or 3x3 block.
// returns true if the is now defined and wasn't before.
func (g Grid) negateOthers(n int) bool {
	var (
		row   = g.getRow(n)
		col   = g.getCol(n)
		block = g.getBlock(n)
	)
	sq := g.squares[n].
		ExcludeDefined(row).
		ExcludeDefined(col).
		ExcludeDefined(block)
	g.squares[n] = sq
	return sq.IsDefined()
}

// reduceByDeduction performs one round of refinement based on the process of
// deduction.  That is, the process of setting a square if it is the ONLY
// square in its row/column/block which can have a particular value.
// Returns a list of squares that are now defined that weren't before.
func (g Grid) reduceByDeduction() []int {
	newlyDefined := make([]int, 0, 8)
	for i := 0; i < 81; i++ {
		isFound := g.deduceSquare(i)
		if isFound {
			newlyDefined = append(newlyDefined, i)
		}
	}
	return newlyDefined
}

// deduceSquare is used to detect cases where there is some value that
// no other square in the row / column / block can possibly be.
// Returns true if the square is now defined and wasn't before.
func (g Grid) deduceSquare(n int) bool {
	if g.squares[n].IsDefined() {
		return false
	}
	var (
		row   = g.getRow(n)
		col   = g.getCol(n)
		block = g.getBlock(n)
	)

	need := Missing(row...)
	if need.IsDefined() {
		g.squares[n] = need
		return true
	}

	need = Missing(col...)
	if need.IsDefined() {
		g.squares[n] = need
		return true
	}

	need = Missing(block...)
	if need.IsDefined() {
		g.squares[n] = need
		return true
	}

	return false
}

// getRow gets the other squares from the same row as the square at index n.
// Excludes square n from the list.
func (g Grid) getRow(n int) []Square {
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

// getCol gets the other squares from the same column as the square at index n.
// Excludes square n from the list.
func (g Grid) getCol(n int) []Square {
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

// getBlock gets the other squares from the same 3x3 block as the square at n.
// Excludes square n from the list.
func (g Grid) getBlock(n int) []Square {
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
