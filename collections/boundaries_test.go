package collections

import (
	"go/token"
	"reflect"
	"testing"
)

func TestBoundaries(t *testing.T) {
	// prepare
	table := map[string]struct {
		b  Boundaries
		p  token.Pos
		in bool
	}{
		"nil boundaries shouldn't containt the pos": {
			b: nil,
			p: token.Pos(1),
		},
		"empty boundaries shouldn't containt the pos": {
			b: Boundaries{},
			p: token.Pos(1),
		},
		"bigger boundaries shouldn't containt the pos": {
			b: Boundaries{
				Boundary{20, 40},
				Boundary{50, 80},
				Boundary{90, 100},
			},
			p: token.Pos(10),
		},
		"lowwer boundaries shouldn't containt the pos": {
			b: Boundaries{
				Boundary{20, 40},
				Boundary{50, 80},
				Boundary{90, 100},
			},
			p: token.Pos(1000),
		},
		"non overlapped boundaries shouldn't containt the pos": {
			b: Boundaries{
				Boundary{1, 15},
				Boundary{20, 40},
				Boundary{50, 80},
				Boundary{90, 100},
				Boundary{120, 200},
				Boundary{250, 512},
				Boundary{600, 999},
			},
			p: token.Pos(88),
		},
		"overlapped boundaries should containt the pos": {
			b: Boundaries{
				Boundary{1, 15},
				Boundary{20, 40},
				Boundary{50, 80},
				Boundary{90, 100},
				Boundary{120, 200},
				Boundary{250, 512},
				Boundary{600, 999},
			},
			p:  token.Pos(510),
			in: true,
		},
	}
	for name, tcase := range table {
		t.Run(name, func(t *testing.T) {
			// exec
			in := tcase.b.Inside(tcase.p)
			// check
			if !reflect.DeepEqual(in, tcase.in) {
				t.Errorf("actual %v doesn't equal to %v", in, tcase.in)
			}
		})
	}
}
