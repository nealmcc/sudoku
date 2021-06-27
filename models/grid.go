package models

import (
	"errors"
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
	copy(sq[:], g.squares[:])
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

// CanSet checks to see if it is legal to set square i to the given value.
// This function is only valid if the grid has been Reduced, and has not
// been added to since.
func (g Grid) CanSet(i, val int) bool {
	curr, next := g.squares[i], squareEnum[val]
	return curr&next > 0
}

// Set assigns the value k to square i
// Reduce should be called after Set, to maintain the integrity of the grid
func (g Grid) Set(i, k int) {
	g.squares[i] = squareEnum[k]
}

// Normalize applies logic to the grid, identifying possible and impossible
// values for each square without using trial and error.
// Returns an error if the grid is invalid.
func (g Grid) Normalize() error {
	for {
		delta := 0

		ids, err := g.reduce()
		if err != nil {
			return err
		}
		delta += len(ids)

		ids, err = g.deduce()
		if err != nil {
			return err
		}
		delta += len(ids)

		if delta == 0 {
			break
		}
	}
	return nil
}

// reduce performs one round of refinement based on the process of
// set subtraction.  That is, the process of excluding from each square the
// values that are definitely assigned within the same row, column or block.
// Returns a list of square indices that are now defined that weren't before.
// Returns an error if any square has been reduced to the point where it cannot
// be any possible value.
func (g Grid) reduce() ([]int, error) {
	newlyDefined := make([]int, 0, 8)
	for i := 0; i < 81; i++ {
		didUpdate, err := g.reduceSquare(i)
		if err != nil {
			return nil, err
		}
		if didUpdate {
			newlyDefined = append(newlyDefined, i)
		}
	}
	return newlyDefined, nil
}

// reduceSquare refines the nth square for this grid by excluding candidate
// values are already defined in the same row, column, or 3x3 block.
// returns true if the square is now defined and wasn't before.
func (g Grid) reduceSquare(n int) (bool, error) {
	if g.squares[n].IsDefined() {
		return false, nil
	}
	var (
		row   = g.getRow(n)
		col   = g.getCol(n)
		block = g.getBlock(n)
	)
	sq := g.squares[n]
	sq = excludeDefined(sq, row)
	sq = excludeDefined(sq, col)
	sq = excludeDefined(sq, block)

	if sq == none {
		return false, errors.New("no possible value for this square")
	}

	g.squares[n] = sq
	return sq.IsDefined(), nil
}

// excludeDefined refines the set of values of the given square
// by any strongly defined squares in the group
func excludeDefined(sq Square, others []Square) Square {
	for _, other := range others {
		if other.IsDefined() {
			sq = sq &^ other
		}
	}
	return sq
}

// deduce performs one round of refinement based on the process of deduction.
// That is, the process of setting a square if it is the ONLY
// square in its row/column/block which can have a particular value.
// Returns a list of square indices that are now defined that weren't before.
// If there is more than one value which a square *must* be, then
// we return an error
func (g Grid) deduce() ([]int, error) {
	newlyDefined := make([]int, 0, 8)
	for i := 0; i < 81; i++ {
		isFound, err := g.deduceSquare(i)
		if err != nil {
			return nil, err
		}
		if isFound {
			newlyDefined = append(newlyDefined, i)
		}
	}
	return newlyDefined, nil
}

// deduceSquare is used to detect cases where there is some value that
// no other square in the row / column / block can possibly be.
// Returns true if the square is now defined and wasn't before.
// Returns an error if this square would need to have more than one value
// to satisfy the row / column / block requirements.
func (g Grid) deduceSquare(n int) (bool, error) {
	if g.squares[n].IsDefined() {
		return false, nil
	}
	var (
		row   = g.getRow(n)
		col   = g.getCol(n)
		block = g.getBlock(n)
	)

	need, err := findMissing(row)
	if err != nil {
		return false, err
	}
	if need != none {
		g.squares[n] = need
		return true, nil
	}

	need, err = findMissing(col)
	if err != nil {
		return false, err
	}
	if need != none {
		g.squares[n] = need
		return true, nil
	}

	need, err = findMissing(block)
	if err != nil {
		return false, err
	}
	if need != none {
		g.squares[n] = need
		return true, nil
	}

	return false, nil
}

// findMissing looks for a single value which none of the given squares can be.
// Returns an error if there is more than one value missing.
func findMissing(group []Square) (Square, error) {
	exists := none
	for _, sq := range group {
		exists |= sq
	}
	notExists := any &^ exists

	if notExists != none && !notExists.IsDefined() {
		return notExists, errors.New("more than one value is missing")
	}

	return notExists, nil
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
