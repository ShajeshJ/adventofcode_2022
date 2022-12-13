package main

import (
	"embed"
	"fmt"

	"github.com/ShajeshJ/adventofcode_2022/common/logging"
	"github.com/ShajeshJ/adventofcode_2022/common/util"
)

var log = logging.GetLogger()

//go:embed part1.txt
var files embed.FS

type Coordinates [2]int

func (c *Coordinates) Row() int {
	return (*c)[0]
}

func (c *Coordinates) Col() int {
	return (*c)[1]
}

type Square struct {
	IsStart   bool
	IsEnd     bool
	Elevation int
	Coords    Coordinates
}

type TraversedSquare struct {
	Self    Square
	Parent  Square
	F, G, H int
}

func getPartOneData() (map[Coordinates]Square, Square, Square) {
	heightmap := map[Coordinates]Square{}

	var start, end Square

	for row, line := range util.ReadProblemInput(files, 1) {
		for col, r := range line {
			p := Square{Coords: Coordinates{row, col}}

			if r >= 'a' && r <= 'z' {
				p.Elevation = int(r - 'a')
			} else if r == 'S' {
				p.IsStart = true
				p.Elevation = 0
				start = p
			} else if r == 'E' {
				p.IsEnd = true
				p.Elevation = 25
				end = p
			} else {
				panic("unexpected character")
			}

			heightmap[p.Coords] = p
		}
	}

	return heightmap, start, end
}

func ManhattanDist(start Square, dest Square) int {
	return (util.Abs(start.Coords.Row()-dest.Coords.Row()) +
		util.Abs(start.Coords.Col()-dest.Coords.Col()))
}

func FindShortestPath(hmap map[Coordinates]Square, start Square, end Square) (TraversedSquare, bool) {
	openList, closedList := map[Coordinates]TraversedSquare{}, map[Coordinates]TraversedSquare{}
	openList[start.Coords] = TraversedSquare{Self: start, F: 0, G: 0, H: 0}

	var found TraversedSquare

	for len(openList) > 0 {
		var q TraversedSquare

		for _, s := range openList {
			if q == (TraversedSquare{}) || s.F < q.F {
				q = s
			}
		}

		delete(openList, q.Self.Coords)

		candidateSuccessorCoords := []Coordinates{
			{q.Self.Coords.Row() - 1, q.Self.Coords.Col()},
			{q.Self.Coords.Row() + 1, q.Self.Coords.Col()},
			{q.Self.Coords.Row(), q.Self.Coords.Col() - 1},
			{q.Self.Coords.Row(), q.Self.Coords.Col() + 1},
		}

		for _, c := range candidateSuccessorCoords {
			s, ok := hmap[c]

			if !ok {
				continue // Invalid coord
			}

			if s.Elevation > q.Self.Elevation+1 {
				continue // Too high to reach
			}

			successor := TraversedSquare{Self: s, Parent: q.Self}
			successor.G = q.G + 1
			successor.H = ManhattanDist(successor.Self, end) + (q.Self.Elevation - s.Elevation)
			successor.F = successor.G + successor.H

			if successor.Self.IsEnd {
				found = successor
				break
			}

			if opened, exists := openList[successor.Self.Coords]; exists && opened.F < successor.F {
				continue // There's a to-be-processed path that's better than this successor
			}

			if closed, exists := closedList[successor.Self.Coords]; exists && closed.F < successor.F {
				continue // There's an already processed path that was better than this successor
			}

			openList[successor.Self.Coords] = successor
			closedList[q.Self.Coords] = q
		}

		if found != (TraversedSquare{}) {
			break
		}

		closedList[q.Self.Coords] = q
	}

	return found, found != (TraversedSquare{})
}

func PartOne() any {
	hmap, start, end := getPartOneData()
	foundPath, found := FindShortestPath(hmap, start, end)
	if !found {
		panic("no path found from start to end")
	}
	return foundPath.G
}

func PartTwo() any {
	hmap, start, end := getPartOneData()
	minFromLowest, found := FindShortestPath(hmap, start, end)
	if !found {
		panic("start should have a path to the end")
	}

	for _, testSquare := range hmap {
		if testSquare.Elevation > 0 {
			continue
		}

		if result, found := FindShortestPath(hmap, testSquare, end); found && result.G < minFromLowest.G {
			minFromLowest = result
		}
	}

	return minFromLowest.G
}

func main() {
	log.Infow(fmt.Sprintf("Answer: %v", PartOne()), "part", 1)
	log.Infow(fmt.Sprintf("Answer: %v", PartTwo()), "part", 2)
}
