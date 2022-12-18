package main

import (
	"reflect"

	"golang.org/x/exp/slices"
)

type Rock interface {
	InitPosition(x, y int)
	WillCollide(dir Direction, chamber [][]rune) bool
	Move(dir Direction, chamber [][]rune) bool
	PlaceRock(chamber *[][]rune)
}

// GetRockGenerator returns a function that will return the rock at
// the given index, and the index for the next rock
func GetRockGenerator() func(i int) (Rock, int) {
	var horizontal *HorizontalRock
	var plus *PlusRock
	var rightangle *RightAngleRock
	var vertical *VerticalRock
	var square *SquareRock

	rocks := map[int]Rock{
		0: horizontal,
		1: plus,
		2: rightangle,
		3: vertical,
		4: square,
	}

	return func(i int) (Rock, int) {
		rockType := reflect.TypeOf(rocks[i%len(rocks)]).Elem()
		rockPtr := reflect.New(rockType)
		rock := rockPtr.Interface().(Rock)
		return rock, (i + 1) % len(rocks)
	}
}

// HorizontalRock is a horizontal line shaped rock.
// Visually it looks like:
//
// ####
type HorizontalRock struct {
	X int // Indicates the left most piece's X position
	Y int // Indicates the left most piece's Y position
}

func (r *HorizontalRock) InitPosition(x, y int) {
	r.X = x
	r.Y = y
}

func (r *HorizontalRock) WillCollide(dir Direction, chamber [][]rune) bool {
	switch dir {
	case Left:
		return r.X <= 0 || chamber[r.Y][r.X-1] == RockRune
	case Right:
		return r.X+4 >= len(chamber[r.Y]) || chamber[r.Y][r.X+4] == RockRune
	case Down:
		return r.Y <= 0 || slices.Contains(chamber[r.Y-1][r.X:r.X+4], RockRune)
	default:
		panic("invalid direction")
	}
}

func (r *HorizontalRock) Move(dir Direction, chamber [][]rune) bool {
	if r.WillCollide(dir, chamber) {
		return false
	}

	switch dir {
	case Left:
		r.X--
	case Right:
		r.X++
	case Down:
		r.Y--
	default:
		panic("invalid direction")
	}
	return true
}

func (r *HorizontalRock) PlaceRock(chamber *[][]rune) {
	(*chamber)[r.Y][r.X] = RockRune
	(*chamber)[r.Y][r.X+1] = RockRune
	(*chamber)[r.Y][r.X+2] = RockRune
	(*chamber)[r.Y][r.X+3] = RockRune
}

// PlusRock is a plus shaped rock.
// Visually it looks like:
//
// .#.
// ###
// .#.
type PlusRock struct {
	X int // Indicates the bottom left empty space's X position
	Y int // Indicates the bottom left empty space's Y position
}

func (r *PlusRock) InitPosition(x, y int) {
	r.X = x
	r.Y = y
}

func (r *PlusRock) WillCollide(dir Direction, chamber [][]rune) bool {
	switch dir {
	case Left:
		return r.X <= 0 ||
			chamber[r.Y][r.X] == RockRune || // Next bottom middle
			chamber[r.Y+1][r.X-1] == RockRune || // Next middle left
			chamber[r.Y+2][r.X] == RockRune // Next top middle
	case Right:
		return r.X+3 >= len(chamber[r.Y]) ||
			chamber[r.Y][r.X+2] == RockRune || // Next bottom middle
			chamber[r.Y+1][r.X+3] == RockRune || // Next middle left
			chamber[r.Y+2][r.X+2] == RockRune // Next top middle
	case Down:
		return r.Y <= 0 ||
			chamber[r.Y][r.X] == RockRune || // Next middle left
			chamber[r.Y-1][r.X+1] == RockRune || // Next bottom middle
			chamber[r.Y][r.X+2] == RockRune // Next middle right
	default:
		panic("invalid direction")
	}
}

func (r *PlusRock) Move(dir Direction, chamber [][]rune) bool {
	if r.WillCollide(dir, chamber) {
		return false
	}

	switch dir {
	case Left:
		r.X--
	case Right:
		r.X++
	case Down:
		r.Y--
	default:
		panic("invalid direction")
	}
	return true
}

func (r *PlusRock) PlaceRock(chamber *[][]rune) {
	(*chamber)[r.Y][r.X+1] = RockRune   // Bottom Middle
	(*chamber)[r.Y+1][r.X] = RockRune   // Middle left
	(*chamber)[r.Y+1][r.X+1] = RockRune // Middle middle
	(*chamber)[r.Y+1][r.X+2] = RockRune // Middle right
	(*chamber)[r.Y+2][r.X+1] = RockRune // Top Middle
}

// RightAngleRock is a backwards L shaped rock.
// Visually it looks like:
//
// ..#
// ..#
// ###
type RightAngleRock struct {
	X int // Indicates the bottom-left most piece's X position
	Y int // Indicates the bottom-left most piece's Y position
}

func (r *RightAngleRock) InitPosition(x, y int) {
	r.X = x
	r.Y = y
}

func (r *RightAngleRock) WillCollide(dir Direction, chamber [][]rune) bool {
	switch dir {
	case Left:
		return r.X <= 0 ||
			chamber[r.Y][r.X-1] == RockRune || // Next bottom left
			chamber[r.Y+1][r.X+1] == RockRune || // Next middle right
			chamber[r.Y+2][r.X+1] == RockRune // Next top right
	case Right:
		return r.X+3 >= len(chamber[r.Y]) ||
			chamber[r.Y][r.X+3] == RockRune || // Next bottom right
			chamber[r.Y+1][r.X+3] == RockRune || // Next middle right
			chamber[r.Y+2][r.X+3] == RockRune // Next top right
	case Down:
		return r.Y <= 0 || slices.Contains(chamber[r.Y-1][r.X:r.X+3], RockRune)
	default:
		panic("invalid direction")
	}
}

func (r *RightAngleRock) Move(dir Direction, chamber [][]rune) bool {
	if r.WillCollide(dir, chamber) {
		return false
	}

	switch dir {
	case Left:
		r.X--
	case Right:
		r.X++
	case Down:
		r.Y--
	default:
		panic("invalid direction")
	}
	return true
}

func (r *RightAngleRock) PlaceRock(chamber *[][]rune) {
	// (*chamber)[r.Y][r.X] = RockRune
	(*chamber)[r.Y][r.X] = RockRune     // Bottom left
	(*chamber)[r.Y][r.X+1] = RockRune   // Bottom middle
	(*chamber)[r.Y][r.X+2] = RockRune   // Bottom right
	(*chamber)[r.Y+1][r.X+2] = RockRune // Middle right
	(*chamber)[r.Y+2][r.X+2] = RockRune // Top right

}

// VerticalRock is vertical line shaped rock.
// Visually it looks like:
//
// #
// #
// #
// #
type VerticalRock struct {
	X int // Indicates the bottom most piece's X position
	Y int // Indicates the bottom most piece's Y position
}

func (r *VerticalRock) InitPosition(x, y int) {
	r.X = x
	r.Y = y
}

func (r *VerticalRock) WillCollide(dir Direction, chamber [][]rune) bool {
	switch dir {
	case Left:
		return r.X <= 0 ||
			chamber[r.Y][r.X-1] == RockRune || // Bottom
			chamber[r.Y+1][r.X-1] == RockRune || // 2nd from bottom
			chamber[r.Y+2][r.X-1] == RockRune || // 2nd from top
			chamber[r.Y+3][r.X-1] == RockRune // Top
	case Right:
		return r.X+1 >= len(chamber[r.Y]) ||
			chamber[r.Y][r.X+1] == RockRune || // Bottom
			chamber[r.Y+1][r.X+1] == RockRune || // 2nd from bottom
			chamber[r.Y+2][r.X+1] == RockRune || // 2nd from top
			chamber[r.Y+3][r.X+1] == RockRune // Top
	case Down:
		return r.Y <= 0 || chamber[r.Y-1][r.X] == RockRune // Bottom
	default:
		panic("invalid direction")
	}
}

func (r *VerticalRock) Move(dir Direction, chamber [][]rune) bool {
	if r.WillCollide(dir, chamber) {
		return false
	}

	switch dir {
	case Left:
		r.X--
	case Right:
		r.X++
	case Down:
		r.Y--
	default:
		panic("invalid direction")
	}
	return true
}

func (r *VerticalRock) PlaceRock(chamber *[][]rune) {
	(*chamber)[r.Y][r.X] = RockRune   // Bottom
	(*chamber)[r.Y+1][r.X] = RockRune // 2nd from bottom
	(*chamber)[r.Y+2][r.X] = RockRune // 2nd from top
	(*chamber)[r.Y+3][r.X] = RockRune // Top
}

// SquareRock is a square shaped rock.
// Visually it looks like:
//
// ##
// ##
type SquareRock struct {
	X int // Indicates the bottom-left most piece's X position
	Y int // Indicates the bottom-left most piece's Y position
}

func (r *SquareRock) InitPosition(x, y int) {
	r.X = x
	r.Y = y
}

func (r *SquareRock) WillCollide(dir Direction, chamber [][]rune) bool {
	switch dir {
	case Left:
		return r.X <= 0 ||
			chamber[r.Y][r.X-1] == RockRune || // Bottom left
			chamber[r.Y+1][r.X-1] == RockRune // Top left
	case Right:
		return r.X+2 >= len(chamber[r.Y]) ||
			chamber[r.Y][r.X+2] == RockRune || // Bottom right
			chamber[r.Y+1][r.X+2] == RockRune // Top right
	case Down:
		return r.Y <= 0 ||
			chamber[r.Y-1][r.X] == RockRune || // Bottom left
			chamber[r.Y-1][r.X+1] == RockRune // Bottom right
	default:
		panic("invalid direction")
	}
}

func (r *SquareRock) Move(dir Direction, chamber [][]rune) bool {
	if r.WillCollide(dir, chamber) {
		return false
	}

	switch dir {
	case Left:
		r.X--
	case Right:
		r.X++
	case Down:
		r.Y--
	default:
		panic("invalid direction")
	}
	return true
}

func (r *SquareRock) PlaceRock(chamber *[][]rune) {
	(*chamber)[r.Y][r.X] = RockRune     // Bottom left
	(*chamber)[r.Y][r.X+1] = RockRune   // Bottom right
	(*chamber)[r.Y+1][r.X] = RockRune   // Top left
	(*chamber)[r.Y+1][r.X+1] = RockRune // Top right
}

// Template-Rock is _ shaped rock.
// Visually it looks like:
//
// #
// type Template-Rock struct {
// 	X int // Indicates the bottom-left most piece's X position
// 	Y int // Indicates the bottom-left most piece's Y position
// }

// func (r* Template-Rock) InitPosition(x, y int) {
// 	r.X = x
// 	r.Y = y
// }

// func (r *Template-Rock) WillCollide(dir Direction, chamber [][]rune) bool {
// 	switch dir {
// 	case Left:
// 		// return
// 	case Right:
// 		// return
// 	case Down:
// 		// return
// 	default:
// 		panic("invalid direction")
// 	}
// }

// func (r *Template-Rock) Move(dir Direction, chamber [][]rune) bool {
// 	if r.WillCollide(dir, chamber) {
// 		return false
// 	}

// 	switch dir {
// 	case Left:
// 		r.X--
// 	case Right:
// 		r.X++
// 	case Down:
// 		r.Y--
// 	default:
// 		panic("invalid direction")
// 	}
// 	return true
// }

// func (r *Template-Rock) PlaceRock(chamber *[][]rune) {
// 	// (*chamber)[r.Y][r.X] = RockRune
// }
