package main

import (
	"embed"
	"fmt"

	"github.com/ShajeshJ/adventofcode_2022/common/logging"
	"github.com/ShajeshJ/adventofcode_2022/common/util"
)

var log = logging.GetLogger()

//go:embed input.txt
var files embed.FS

const (
	RIGHT = iota
	DOWN
	LEFT
	UP
)

type Facing int

const (
	WALL = '#'
	OPEN = '.'
)

type TileType rune

type Tile struct {
	Up    *Tile
	Down  *Tile
	Left  *Tile
	Right *Tile
	Type  TileType
	Row   int
	Col   int
}

func (t *Tile) Walk(dir Facing) (*Tile, bool) {
	switch dir {
	case RIGHT:
		if t.Right.Type == WALL {
			return t, false
		}
		return t.Right, true
	case DOWN:
		if t.Down.Type == WALL {
			return t, false
		}
		return t.Down, true
	case LEFT:
		if t.Left.Type == WALL {
			return t, false
		}
		return t.Left, true
	case UP:
		if t.Up.Type == WALL {
			return t, false
		}
		return t.Up, true
	}
	panic("Invalid direction")
}

const (
	CW  = 'R'
	CCW = 'L'
)

type Rotation rune

type Instruction struct {
	IsTurn bool
	Turn   Rotation
	Steps  int
}

func (i *Instruction) TurnDir(curDir Facing) Facing {
	switch i.Turn {
	case CW:
		return (curDir + 1) % 4
	case CCW:
		// go's (x-1 % 4) behaves differently than in python
		// https://github.com/golang/go/issues/448#issuecomment-66049769
		return (curDir + 3) % 4
	}
	panic("Invalid turn")
}

func buildInstructions(instructionStr string) []Instruction {
	instructions := []Instruction{}

	numBuilder := ""
	for _, c := range instructionStr {
		if c != CW && c != CCW {
			numBuilder += string(c)
			continue
		}

		// We have a turn, so add the steps and the turn
		if numBuilder != "" {
			instructions = append(instructions, Instruction{Steps: util.AtoiNoError(numBuilder)})
			numBuilder = ""
		}
		instructions = append(instructions, Instruction{IsTurn: true, Turn: Rotation(c)})
	}

	if numBuilder != "" {
		// Add the last steps, if any
		instructions = append(instructions, Instruction{Steps: util.AtoiNoError(numBuilder)})
		numBuilder = ""
	}

	return instructions
}

func getPartOneData() ([][]*Tile, []Instruction) {
	data := util.ReadProblemInput(files, 1)
	instructions := buildInstructions(data[len(data)-1])
	data = data[:len(data)-2] // Remove non-map data

	height := len(data)
	width := len(data[0])
	for _, line := range data {
		if len(line) > width {
			width = len(line)
		}
	}

	board := make([][]*Tile, height) // Temporary board to help build tile connections

	for i, line := range data {
		board[i] = make([]*Tile, width)
		for j, c := range line {
			if c != WALL && c != OPEN {
				continue // Empty spaces that are not part of the map
			}
			board[i][j] = &Tile{Type: TileType(c), Row: i + 1, Col: j + 1}
		}
	}

	// Build vertical connections across each column
	for j := 0; j < width; j++ {
		var topMost, prev *Tile
		for i := 0; i < height; i++ {
			if board[i][j] == nil {
				continue
			}
			if topMost == nil {
				topMost = board[i][j]
				prev = board[i][j]
				continue
			}
			prev.Down = board[i][j]
			board[i][j].Up = prev
			prev = board[i][j]
		}
		// Connect bottom most to the top most
		prev.Down = topMost
		topMost.Up = prev
	}

	// Build horizontal connections across each row
	for i := 0; i < height; i++ {
		var leftMost, prev *Tile
		for j := 0; j < width; j++ {
			if board[i][j] == nil {
				continue
			}
			if leftMost == nil {
				leftMost = board[i][j]
				prev = board[i][j]
				continue
			}
			prev.Right = board[i][j]
			board[i][j].Left = prev
			prev = board[i][j]
		}
		// Connect right most to the left most
		prev.Right = leftMost
		leftMost.Left = prev
	}

	return board, instructions
}

func getStartingTile(board [][]*Tile) *Tile {
	for i := 0; i < len(board[0]); i++ {
		if board[0][i] != nil && board[0][i].Type == OPEN {
			return board[0][i]
		}
	}

	panic("No starting tile found")
}

func PartOne() any {
	board, instructions := getPartOneData()
	curTile := getStartingTile(board)
	var curDir Facing = RIGHT

	for _, instruction := range instructions {
		if instruction.IsTurn {
			curDir = instruction.TurnDir(curDir)
			continue
		}
		for i := 0; i < instruction.Steps; i++ {
			var walked bool
			curTile, walked = curTile.Walk(curDir)
			if !walked {
				break
			}
		}
	}

	return curTile.Row*1000 + curTile.Col*4 + int(curDir)
}

type Remap struct {
	t   *Tile
	dir Facing
}

// getCubeRemapping manually remaps the direction from a tile to
// a new tile/direction, such that the map connects like a cube net.
// Assumes the following cube net, where each # is a cube face of 50x50:
//
//	.12
//	.3.
//	45.
//	6..
func getCubeRemapping(board [][]*Tile) [][]map[Facing]Remap {
	remaps := make([][]map[Facing]Remap, len(board))
	for i := range remaps {
		remaps[i] = make([]map[Facing]Remap, len(board[i]))
		for j := range remaps[i] {
			remaps[i][j] = make(map[Facing]Remap)
		}
	}

	for i := 0; i < 50; i++ {
		// 1 <-> 4 (4 rotated 180 degrees left of 1)
		remaps[i][50][LEFT] = Remap{board[149-i][0], RIGHT}
		remaps[149-i][0][LEFT] = Remap{board[i][50], RIGHT}

		// 1 <-> 6 (6 rotated 90 degrees CCW, and above 1)
		remaps[0][50+i][UP] = Remap{board[150+i][0], RIGHT}
		remaps[150+i][0][LEFT] = Remap{board[0][50+i], DOWN}

		// 2 <-> 6 (6 above 2 directly)
		remaps[0][100+i][UP] = Remap{board[199][i], UP}
		remaps[199][i][DOWN] = Remap{board[0][100+i], DOWN}

		// 2 <-> 3 (3 rotated 90 degrees CCW, below 2)
		remaps[49][100+i][DOWN] = Remap{board[50+i][99], LEFT}
		remaps[50+i][99][RIGHT] = Remap{board[49][100+i], UP}

		// 2 <-> 5 (5 rotated 180 degrees, right of 2)
		remaps[i][149][RIGHT] = Remap{board[149-i][99], LEFT}
		remaps[149-i][99][RIGHT] = Remap{board[i][149], LEFT}

		// 3 <-> 4 (4 rotated 90 degrees CW, left of 3)
		remaps[50+i][50][LEFT] = Remap{board[100][i], DOWN}
		remaps[100][i][UP] = Remap{board[50+i][50], RIGHT}

		// 5 <-> 6 (6 rotated 90 degrees CCW, below 5)
		remaps[149][50+i][DOWN] = Remap{board[150+i][49], LEFT}
		remaps[150+i][49][RIGHT] = Remap{board[149][50+i], UP}
	}

	return remaps
}

func walkWithRemap(curTile *Tile, curDir Facing, remaps [][]map[Facing]Remap) (*Tile, bool, Facing) {
	defaultWalk := func() (*Tile, bool, Facing) {
		newTile, walked := curTile.Walk(curDir)
		return newTile, walked, curDir
	}

	remap := remaps[curTile.Row-1][curTile.Col-1]

	if len(remap) == 0 {
		return defaultWalk()
	}

	next, noRemapForDir := remap[curDir]
	if !noRemapForDir {
		return defaultWalk()
	}

	if next.t.Type == WALL {
		return curTile, false, curDir
	}

	return next.t, true, next.dir
}

func PartTwo() any {
	board, instructions := getPartOneData()
	curTile := getStartingTile(board)
	remaps := getCubeRemapping(board)
	var curDir Facing = RIGHT

	for _, instruction := range instructions {
		if instruction.IsTurn {
			curDir = instruction.TurnDir(curDir)
			continue
		}
		for i := 0; i < instruction.Steps; i++ {
			var walked bool
			curTile, walked, curDir = walkWithRemap(curTile, curDir, remaps)
			if !walked {
				break
			}
		}
	}

	return curTile.Row*1000 + curTile.Col*4 + int(curDir)
}

func main() {
	log.Infow(fmt.Sprintf("Answer: %v", PartOne()), "part", 1)
	log.Infow(fmt.Sprintf("Answer: %v", PartTwo()), "part", 2)
}
