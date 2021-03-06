package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"mcconachie.co/sudoku/models"
)

var cases_solve = []struct {
	name, in, want string
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
		want: `435 269 781
682 571 493
197 834 562

826 195 347
374 682 915
951 743 628

519 326 874
248 957 136
763 418 259
`,
	},
	{
		name: "easy sudoku (googled)",
		in: `
			.6. 3.. 8.4
			537 .9. ...
			.4. ..6 3.7

			.9. .51 238
			... ... ...
			713 62. .4.

			3.6 4.. .1.
			... .6. 523
			1.2 ..9 .8.`,
		want: `261 375 894
537 894 162
948 216 357

694 751 238
825 943 671
713 628 945

356 482 719
489 167 523
172 539 486
`,
	},
	{
		name: "mit courseware example",
		in: `
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
		want: `596 184 273
731 629 854
482 753 961

927 538 146
815 462 397
364 917 528

158 246 739
249 371 685
673 895 412
`,
	},
	{
		name: "17 clues",
		in: `
			... 8.1 ...
			... ... .43
			5.. ... ...

			... .7. 8..
			... ... 1..
			.2. .3. ...

			6.. ... .75
			..3 4.. ...
			... 2.. 6..
		`,
		want: `237 841 569
186 795 243
594 326 718

315 674 892
469 582 137
728 139 456

642 918 375
853 467 921
971 253 684
`,
	},
	{
		name: "requires backtracking",
		in: `
			1.. 9.7 ..3
			.8. ... .7.
			..9 ... 6..

			..7 2.9 4..
			41. ... .95
			..8 5.4 3..

			..3 ... 7..
			.5. ... .4.
			2.. 8.6 ..9`,
		want: `164 957 283
385 621 974
729 438 651

537 289 416
412 763 895
698 514 327

843 195 762
956 372 148
271 846 539
`,
	},
}

func TestSolve(t *testing.T) {
	for _, tc := range cases_solve {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// t.Parallel()
			grid := models.NewGrid([]byte(tc.in))
			done, n := solve(&grid)
			assert.Equal(t, true, done)
			t.Log(n, "backtracks")
			assert.Equal(t, tc.want, grid.String())
		})
	}
}

func BenchmarkSolve(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for _, tc := range cases_solve {
			g := models.NewGrid([]byte(tc.in))
			solve(&g)
		}
	}
}
