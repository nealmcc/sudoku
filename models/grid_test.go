package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var casesNewGrid = []struct {
	name    string
	in      string
	want    [81]Square
	display string
}{
	{
		name: "a filled grid",
		in: `
				435 269 781
				682 571 493
				197 834 562

				826 195 347
				374 682 915
				951 743 628

				519 326 874
				248 957 136
				763 418 259
			`,
		want: [81]Square{
			four, three, five, two, six, nine, seven, eight, one,
			six, eight, two, five, seven, one, four, nine, three,
			one, nine, seven, eight, three, four, five, six, two,

			eight, two, six, one, nine, five, three, four, seven,
			three, seven, four, six, eight, two, nine, one, five,
			nine, five, one, seven, four, three, six, two, eight,

			five, one, nine, three, two, six, eight, seven, four,
			two, four, eight, nine, five, seven, one, three, six,
			seven, six, three, four, one, eight, two, five, nine,
		},
		display: `╔═══╤═══╤═══╗
║435│269│781║
║682│571│493║
║197│834│562║
╟───┼───┼───╢
║826│195│347║
║374│682│915║
║951│743│628║
╟───┼───┼───╢
║519│326│874║
║248│957│136║
║763│418│259║
╚═══╧═══╧═══╝`,
	},
}

func TestNewGrid(t *testing.T) {
	for _, tc := range casesNewGrid {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := NewGrid([]byte(tc.in))
			require.Exactly(t, tc.want, *(got.squares))
		})
	}
}

func BenchmarkNewGrid(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for _, tc := range casesNewGrid {
			NewGrid([]byte(tc.in))
		}
	}
}

func TestString(t *testing.T) {
	for _, tc := range casesNewGrid {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			grid := NewGrid([]byte(tc.in))
			require.Equal(t, tc.display, grid.String())
		})
	}
}

var casesRefineOne = []struct {
	name string
	n    int
	grid Grid
	want Square
}{
	{
		name: "refine by row (row 0)",
		n:    4,
		grid: NewGrid([]byte(`
			435 2.9 781
			... ... ...
			... ... ...

			... ... ...
			... ... ...
			... ... ...

			... ... ...
			... ... ...
			... ... ...
			`)),
		want: six,
	},
	{
		name: "refine by row and column (7, 7)",
		n:    70,
		grid: NewGrid([]byte(`
			... ... ...
			... ... .8.
			... ... ...

			... ... ...
			... ... .2.
			... ... ...

			... ... ...
			4.5 .6. 7.1
			... ... ...
			`)),
		want: three | nine,
	},
	{
		name: "refine by block (middle)",
		n:    40,
		grid: NewGrid([]byte(`
			... ... ...
			... ... ...
			... ... ...

			... 123 ...
			... 5.6 ...
			... 987 ...

			... ... ...
			... ... ...
			... ... ...
			`)),
		want: four,
	},
	{
		name: "bottom right",
		n:    80,
		grid: NewGrid([]byte(`
			... ... ..1
			... ... ..2
			... ... ..3

			... ... ...
			... ... ...
			... ... ...

			... ... .5.
			... ... 6..
			... 7.9 ...
			`)),
		want: four | eight,
	},
}

func TestRefineOne(t *testing.T) {
	for _, tc := range casesRefineOne {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := require.New(t)

			tc.grid.excludeSquare(tc.n)

			got := tc.grid.squares[tc.n]
			r.Equal(tc.want, got)
		})
	}
}

func BenchmarkRefineOne(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for _, tc := range casesRefineOne {
			tc.grid.excludeSquare(tc.n)
		}
	}
}

var casesRefine = []struct {
	name string
	in   string
	want [81][]int
}{
	{
		"easy sudoku (google)",
		`
			.6. 3.. 8.4
			537 .9. ...
			.4. ..6 3.7

			.9. .51 238
			... ... ...
			713 62. .4.

			3.6 4.. .1.
			... .6. 523
			1.2 ..9 .8.`,
		[81][]int{
			{2}, {6}, {1}, {3}, {7}, {5}, {8}, {9}, {4},
			{5}, {3}, {7}, {8}, {9}, {4}, {1}, {6}, {2},
			{9}, {4}, {8}, {2}, {1}, {6}, {3}, {5}, {7},

			{6}, {9}, {4}, {7}, {5}, {1}, {2}, {3}, {8},
			{8}, {2}, {5}, {9}, {4}, {3}, {6}, {7}, {1},
			{7}, {1}, {3}, {6}, {2}, {8}, {9}, {4}, {5},

			{3}, {5}, {6}, {4}, {8}, {2}, {7}, {1}, {9},
			{4}, {8}, {9}, {1}, {6}, {7}, {5}, {2}, {3},
			{1}, {7}, {2}, {5}, {3}, {9}, {4}, {8}, {6},
		},
	}, {
		"mit courseware example",
		`
			... 1.4 ...
			..1 ... 8..
			.8. 7.3 .6.

			9.7 ... 1.6
			... ... ...
			3.4 ... 5.8

			.5. 2.6 .3.
			..9 ... 6..
			... 8.5 ...
		`,
		[81][]int{
			{5, 6, 7}, {3, 6, 7, 9}, {3, 6}, {1}, {8}, {4}, {2, 3, 7, 9}, {2, 5, 7, 9}, {2, 3, 7, 9},
			{4, 5, 6, 7}, {3, 4, 6, 7, 9}, {1}, {5, 6, 9}, {2, 5, 6, 9}, {2, 9}, {8}, {2, 5, 7, 9}, {2, 3, 4, 7, 9},
			{4, 5}, {8}, {2}, {7}, {5, 9}, {3}, {4, 9}, {6}, {1},
			{9}, {2}, {7}, {3, 5}, {3, 5}, {8}, {1}, {4}, {6},
			{8}, {1, 6}, {5}, {3, 4, 6, 9}, {1, 2, 3, 4, 6, 7, 9}, {1, 2, 7, 9}, {2, 3, 7, 9}, {2, 7, 9}, {2, 3, 7, 9},
			{3}, {1, 6}, {4}, {6, 9}, {1, 2, 6, 7, 9}, {1, 2, 7, 9}, {5}, {2, 7, 9}, {8},
			{1}, {5}, {8}, {2}, {4, 7, 9}, {6}, {4, 7, 9}, {3}, {4, 7, 9},
			{2}, {3, 4, 7}, {9}, {3, 4}, {1, 3, 4, 7}, {1, 7}, {6}, {8}, {5},
			{4, 6, 7}, {3, 4, 6, 7}, {3, 6}, {8}, {3, 4, 7, 9}, {5}, {2, 4, 7, 9}, {1}, {2, 4, 7, 9},
		},
	},
}

func TestRefine(t *testing.T) {
	for _, tc := range casesRefine {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			grid := NewGrid([]byte(tc.in))
			numLoops := 0
			for {
				numChanges := 0
				changes := grid.reduceByExclusion()
				numChanges += len(changes)
				changes = grid.reduceByDeduction()
				numChanges += len(changes)
				if numChanges == 0 {
					break
				}
				numLoops++
			}
			t.Logf("after %d loops:\n%s\n%s\n", numLoops, grid, grid.genSquaresData())
			want := rebuildSquares(tc.want)
			require.Equal(t, *want, *grid.squares)
		})
	}
}

var casesDeduceOne = []struct {
	name      string
	grid      [81][]int
	n         int
	didChange bool
}{
	{
		"combine deductions when one group has none missing",
		[81][]int{
			{2, 5, 6, 7}, {3, 6, 7, 9}, {2, 3, 5, 6}, {1}, {8}, {4}, {2, 3, 7, 9}, {2, 5, 7, 9}, {2, 3, 5, 7, 9},
			{2, 4, 5, 6, 7}, {3, 4, 6, 7, 9}, {1}, {5, 6, 9}, {2, 5, 6, 9}, {2, 9}, {8}, {2, 5, 7, 9}, {2, 3, 4, 5, 7, 9},
			{2, 4, 5}, {8}, {2, 5}, {7}, {2, 5, 9}, {3}, {2, 4, 9}, {6}, {1, 2, 4, 5, 9},
			{9}, {2}, {7}, {3, 5}, {3, 5}, {8}, {1}, {4}, {6},
			{8}, {1, 6}, {5, 6}, {3, 4, 5, 6, 9}, {1, 2, 3, 4, 5, 6, 7, 9}, {1, 2, 7, 9}, {2, 3, 7, 9}, {2, 7, 9}, {2, 3, 7, 9},
			{3}, {1, 6}, {4}, {6, 9}, {1, 2, 6, 7, 9}, {1, 2, 7, 9}, {5}, {2, 7, 9}, {8},
			{1, 4, 7}, {5}, {8}, {2}, {1, 4, 7, 9}, {6}, {4, 7, 9}, {3}, {1, 4, 7, 9},
			{1, 2, 4, 7}, {1, 3, 4, 7}, {9}, {3, 4}, {1, 3, 4, 7}, {1, 7}, {6}, {8}, {1, 2, 4, 5, 7},
			{1, 2, 4, 6, 7}, {1, 3, 4, 6, 7}, {2, 3, 6}, {8}, {1, 3, 4, 7, 9}, {5}, {2, 4, 7, 9}, {1, 2, 7, 9}, {1, 2, 4, 7, 9},
		},
		26, true,
	},
}

func TestDeduceOne(t *testing.T) {
	for _, tc := range casesDeduceOne {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			grid := Grid{rebuildSquares(tc.grid)}
			got := grid.deduceSquare(tc.n)
			assert.Equal(t, tc.didChange, got)
		})
	}
}

func BenchmarkDeduceOne(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for _, tc := range casesDeduceOne {
			grid := Grid{rebuildSquares(tc.grid)}
			grid.deduceSquare(tc.n)
		}
	}
}

func rebuildSquares(want [81][]int) *[81]Square {
	w := [81]Square{}
	for i, vals := range want {
		w[i] = none.Include(vals...)
	}
	return &w
}
