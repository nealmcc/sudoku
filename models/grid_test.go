package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var newgridcases = []struct {
	name        string
	in          string
	want        [81]Square
	expectedErr error
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
		expectedErr: nil,
	},
}

func TestNewGrid(t *testing.T) {
	for _, tc := range newgridcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := require.New(t)

			got, err := NewGrid(tc.in)

			if tc.expectedErr != nil {
				r.Error(tc.expectedErr, err)
			} else {
				r.NoError(err)
			}

			r.Exactly(tc.want, *(got.squares))
		})
	}
}

func BenchmarkNewGrid(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for _, tc := range newgridcases {
			NewGrid(tc.in)
		}
	}
}
