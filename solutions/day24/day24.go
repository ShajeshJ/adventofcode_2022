package main

import (
	"embed"
	"fmt"

	ds "github.com/ShajeshJ/adventofcode_2022/common/datastructures"
	"github.com/ShajeshJ/adventofcode_2022/common/logging"
	"github.com/ShajeshJ/adventofcode_2022/common/util"
	"golang.org/x/exp/slices"
)

var log = logging.GetLogger()

//go:embed input.txt
var files embed.FS

type Direction rune

const (
	NORTH = Direction('^')
	WEST  = Direction('<')
	EAST  = Direction('>')
	SOUTH = Direction('v')
)

type Position struct {
	X, Y int
}

type Blizzard struct {
	Position
	Dir Direction
}

type SimState struct {
	StartPos  Position
	EndPos    Position
	blizzards []Blizzard
	h, w      int
	maxCycles int
}

func (ss *SimState) UpdateBlizzards() {
	for i := 0; i < len(ss.blizzards); i++ {
		switch ss.blizzards[i].Dir {
		case NORTH:
			ss.blizzards[i].Y = (ss.blizzards[i].Y - 1 + ss.h) % ss.h
		case SOUTH:
			ss.blizzards[i].Y = (ss.blizzards[i].Y + 1) % ss.h
		case WEST:
			ss.blizzards[i].X = (ss.blizzards[i].X - 1 + ss.w) % ss.w
		case EAST:
			ss.blizzards[i].X = (ss.blizzards[i].X + 1) % ss.w
		}
	}
}

func (ss *SimState) GetValidMoves(p Position) []Position {
	allMoves := []Position{
		{p.X, p.Y + 1}, // South
		{p.X + 1, p.Y}, // East
		p,              // Wait
		{p.X, p.Y - 1}, // North
		{p.X - 1, p.Y}, // West
	}

	validMoves := []Position{}

	for _, m := range allMoves {
		if m.X >= 0 && m.X < ss.w && m.Y >= 0 && m.Y < ss.h {
			validMoves = append(validMoves, m)
		}
		if m == ss.StartPos || m == ss.EndPos {
			validMoves = append(validMoves, m)
		}
	}
	return validMoves
}

func getPartOneData() SimState {
	data := util.ReadProblemInput(files)
	h, w := len(data)-2, len(data[0])-2
	blizzards := make([]Blizzard, 0)
	start, end := Position{}, Position{}

	for i, line := range data {
		for j, c := range line {
			switch c {
			case '#':
			case '.':
				if i == 0 {
					start = Position{j - 1, i - 1}
				}
				if i == len(data)-1 {
					end = Position{j - 1, i - 1}
				}
			default:
				blizzards = append(blizzards, Blizzard{Position{j - 1, i - 1}, Direction(c)})
			}
		}
	}

	return SimState{
		StartPos:  start,
		EndPos:    end,
		blizzards: blizzards,
		h:         h,
		w:         w,
		maxCycles: util.LCM(h, w),
	}
}

type RepeatState struct {
	x, y, cycle int
}

func FindMinTravelTime(ss *SimState) int {
	q := ds.Queue[Position]{ss.StartPos}
	timeTaken := 0
	curCycle := 0
	visited := make(map[RepeatState]bool)

	for {
		nextQ := ds.Queue[Position]{}

		timeTaken++
		curCycle = (curCycle + 1) % ss.maxCycles
		ss.UpdateBlizzards()

		for !q.IsEmpty() {
			p, _ := q.Dequeue()
			if visited[RepeatState{p.X, p.Y, curCycle}] {
				continue
			}
			for _, move := range ss.GetValidMoves(p) {
				if move == ss.EndPos {
					return timeTaken
				}
				if blizzard := slices.IndexFunc(
					ss.blizzards,
					func(b Blizzard) bool { return b.Position == move },
				); blizzard != -1 {
					// Blizzard in the way, can't move here
					continue
				}

				nextQ.Enqueue(move)
			}
			visited[RepeatState{p.X, p.Y, curCycle}] = true
		}

		q = nextQ
		// log.Info("done minute: ", timeTaken)
	}
}

func PartOne() any {
	simState := getPartOneData()
	return FindMinTravelTime(&simState)
}

func PartTwo() any {
	simState := getPartOneData()
	leg1 := FindMinTravelTime(&simState) // Start -> End
	// log.Info("minutes to travel from start to end: ", leg1)

	simState.StartPos, simState.EndPos = simState.EndPos, simState.StartPos
	leg2 := FindMinTravelTime(&simState) // End -> Start
	// log.Info("minutes to travel from end back to the start: ", leg2)

	simState.StartPos, simState.EndPos = simState.EndPos, simState.StartPos
	leg3 := FindMinTravelTime(&simState) // Start -> End
	// log.Info("minutes to travel from start back to the end: ", leg3)

	return leg1 + leg2 + leg3
}

func main() {
	log.Infow(fmt.Sprintf("Answer: %v", PartOne()), "part", 1)
	log.Infow(fmt.Sprintf("Answer: %v", PartTwo()), "part", 2)
}
