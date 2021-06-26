package models

import (
	"errors"
	"strings"
)

type Grid struct {
	squares *[81]Square
}

// NewGrid initializes a sudoku grid using the given string.
// Each square should either be given as a digit (if the square is defined)
// or as a '0' or '.' if the square could be anything.
// All other characters are ignored.
func NewGrid(start string) (Grid, error) {
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

	if i != 81 {
		return Grid{}, errors.New("the grid must have 81 squares")
	}
	return g, nil
}

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
