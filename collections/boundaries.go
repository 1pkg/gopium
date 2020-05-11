package collections

import "go/token"

// Boundary defines sorted pos pair type
type Boundary struct {
	First token.Pos
	Last  token.Pos
}

// Less checks if pos strictly less then boundary
func (b Boundary) Less(p token.Pos) bool {
	return b.Last < p
}

// Greater checks if pos strictly greater then boundary
func (b Boundary) Greater(p token.Pos) bool {
	return b.First > p
}

// Inside checks if pos inside boundary
func (b Boundary) Inside(p token.Pos) bool {
	return !b.Less(p) && !b.Greater(p)
}

// Boundaries defines ordered set of boundary
type Boundaries []Boundary

// Inside checks if pos inside boundaries
// by using binary search to check boundaries
func (bs Boundaries) Inside(p token.Pos) bool {
	// use binary search to check boundaries
	l, r := 0, len(bs)-1
	for l <= r {
		// calculate the index
		i := (l + r) / 2
		b := bs[i]
		// if pos is inside
		// we found the answer
		if b.Inside(p) {
			return true
		}
		// if pos is inside
		// left half search there
		if b.Greater(p) {
			r = i - 1
			continue
		}
		// if comment is inside
		// right half search there
		if b.Less(p) {
			l = i + 1
			continue
		}
	}
	// we found the answer
	return false
}
