package models

type Grid interface {
}

type grid struct {
	rows [9]line
}

type line [9]square

func NewGrid(start string) (Grid, error) {
	return &grid{}, nil
}
