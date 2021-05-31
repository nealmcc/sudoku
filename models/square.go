package models

// A Square represents the set of possible values for a given sudoku square
type Square interface {
	Values() []int
	IsSingle() bool
	Include(...int) Square
	Exclude(...int) Square
}

type square uint16

var _ Square = square(0)

const (
	one square = 1 << iota
	two
	three
	four
	five
	six
	seven
	eight
	nine
	any  square = (1 << iota) - 1
	none square = 0
)

var squares = [10]square{
	any, one, two, three, four, five, six, seven, eight, nine,
}

// New square returns a square that is initialised to be the given value
// The values 1-9 will return a square that can only be that value.
// If 0 is given, then the square could only any of 1-9
func NewSquare(n int) Square {
	return squares[n]
}

// Values returns all the potential values that this square could hold
func (sq square) Values() []int {
	vals := make([]int, 0, 9)
	for i := 0; i < 9; i++ {
		if sq&(1<<i) > 0 {
			vals = append(vals, i+1)
		}
	}
	return vals
}

// IsSingle is used to determine if a square has exactly one value.
// if sq.IsSingle() is true, then len(sq.Values) == 1 and vice versa.
func (sq square) IsSingle() bool {
	return sq != 0 && sq&(sq-1) == 0
}

// Include is used when adding possible values to a given square.
// for example, if a row currently has no 2,3,7 then we could add 2,3,7 to
// any empty squares in that row.
func (sq square) Include(others ...int) Square {
	for _, n := range others {
		sq = sq | squares[n]
	}
	return sq
}

// Exclude is used to rule out possible values for a given square.
// For example, if a square is on a row with 1,4,5 then we exclude 1,4,5 from
// this square's possible values.
func (sq square) Exclude(others ...int) Square {
	for _, n := range others {
		sq = sq & ^squares[n]
	}
	return sq
}
