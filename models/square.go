package models

// A Square represents the set of possible values for a given sudoku square
type Square uint16

const (
	one Square = 1 << iota
	two
	three
	four
	five
	six
	seven
	eight
	nine
	any  Square = (1 << iota) - 1
	none Square = 0
)

var squareEnum = [10]Square{
	any, one, two, three, four, five, six, seven, eight, nine,
}

// New square returns a square that is initialised to be the given value
// The values 1-9 will return a square that can only be that value.
// If 0 is given, then the square could be any of 1-9
func NewSquare(n int) Square {
	return squareEnum[n]
}

// Display returns the display character for this square.
// For squares which are defined, we return their digit.
// For squares which are uncertain, we return '.'
// for squares which are impossible, we return '!'
func (sq Square) Display() byte {
	if sq.IsDefined() {
		var exp byte
		for sq > 0 {
			sq >>= 1
			exp++
		}
		return '0' + exp
	}
	if sq == none {
		return '!'
	}
	return '.'
}

// Values returns all the potential values that this square could hold
func (sq Square) Values() []int {
	vals := make([]int, 0, 9)
	for i := 0; i < 9; i++ {
		if sq&(1<<i) > 0 {
			vals = append(vals, i+1)
		}
	}
	return vals
}

// IsDefined is used to determine if a square has exactly one value.
// if sq.IsDefined() is true, then len(sq.Values) == 1 and vice versa.
func (sq Square) IsDefined() bool {
	return sq != 0 && sq&(sq-1) == 0
}

// ExcludeDefined refines the value of the given square, by excluding
// any definite squares from the given group.
func (sq Square) ExcludeDefined(group []Square) Square {
	for _, other := range group {
		if sq.IsDefined() {
			break
		}
		if other.IsDefined() {
			sq = sq &^ other
		}
	}
	return sq
}

// Missing finds the set of values which are not possible in the given group.
func Missing(group ...Square) Square {
	notFound := none
	for i := 1; i <= 9; i++ {
		needle, found := squareEnum[i], false
		for _, sq := range group {
			if sq&needle > 0 {
				found = true
				break
			}
		}
		if !found {
			notFound |= needle
		}
	}
	return notFound
}
