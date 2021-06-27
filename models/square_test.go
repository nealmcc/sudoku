package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSquareEqual(t *testing.T) {
	tt := []struct {
		name string
		a, b Square
	}{
		{"One", one, 0x01},
		{"Two", two, 0x02},
		{"Three", three, 0x04},
		{"Four", four, 0x08},
		{"Five", five, 0x10},
		{"Six", six, 0x20},
		{"Seven", seven, 0x40},
		{"Eight", eight, 0x80},
		{"Nine", nine, 0x0100},
		{"Any", any, 0x01FF},
		{"None", none, 0},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := require.New(t)

			r.Exactlyf(tc.b, tc.a, "%09b should equal %09b", tc.a, tc.b)
		})
	}
}

func TestSquareValues(t *testing.T) {
	tt := []struct {
		name string
		in   Square
		want []int
	}{
		{
			name: "one,two,three",
			in:   one | two | three,
			want: []int{1, 2, 3},
		}, {
			name: "four,five,six",
			in:   four | five | six,
			want: []int{4, 5, 6},
		}, {
			name: "seven,eight,nine",
			in:   seven | eight | nine,
			want: []int{7, 8, 9},
		}, {
			name: "any",
			in:   any,
			want: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		}, {
			name: "none",
			in:   none,
			want: []int{},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := require.New(t)

			got := tc.in.Values()

			r.Equal(tc.want, got)
		})
	}
}

func TestSquareIsSingle(t *testing.T) {
	tt := []struct {
		name string
		in   Square
		want bool
	}{
		{"one", one, true},
		{"two", two, true},
		{"three", three, true},
		{"four", four, true},
		{"five", five, true},
		{"six", six, true},
		{"seven", seven, true},
		{"eight", eight, true},
		{"nine", nine, true},
		{"any", any, false},
		{"none", none, false},
		{"six + four", six | four, false},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := require.New(t)

			got := tc.in.IsDefined()

			r.Exactly(tc.want, got)
		})
	}
}
