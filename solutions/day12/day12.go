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
	Square
	F, G, H int
}

func getPartOneData() (map[Coordinates]Square, Square, Square) {
	heightmap := map[Coordinates]Square{}

	var start, end Square

	for row, line := range util.ReadProblemInput(files) {
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

// FindShortestPath finds the shortest path from start to end using A* algorithm
func FindShortestPath(
	hmap map[Coordinates]Square,
	start Square,
	calcHeuristic func(s, q TraversedSquare) int,
	isValidSuccessor func(s, q TraversedSquare) bool,
	isEnd func(s TraversedSquare) bool,
) (TraversedSquare, bool) {

	openList := map[Coordinates]TraversedSquare{start.Coords: {Square: start}}
	closedList := map[Coordinates]TraversedSquare{}

	var found TraversedSquare

	for len(openList) > 0 {
		var q TraversedSquare

		for _, s := range openList {
			if q == (TraversedSquare{}) || s.F < q.F {
				q = s
			}
		}

		delete(openList, q.Coords)

		candidateSuccessorCoords := []Coordinates{
			{q.Coords.Row() - 1, q.Coords.Col()},
			{q.Coords.Row() + 1, q.Coords.Col()},
			{q.Coords.Row(), q.Coords.Col() - 1},
			{q.Coords.Row(), q.Coords.Col() + 1},
		}

		for _, c := range candidateSuccessorCoords {
			s, ok := hmap[c]

			if !ok {
				continue // Invalid coord
			}

			successor := TraversedSquare{Square: s}
			successor.G = q.G + 1
			successor.H = calcHeuristic(successor, q)
			successor.F = successor.G + successor.H

			if !isValidSuccessor(successor, q) {
				continue // successor is unreachable for whatever reason
			}

			if isEnd(successor) {
				found = successor
				break
			}

			if opened, exists := openList[successor.Coords]; exists && opened.F < successor.F {
				continue // There's a to-be-processed path that's better than this successor
			}

			if closed, exists := closedList[successor.Coords]; exists && closed.F < successor.F {
				continue // There's an already processed path that was better than this successor
			}

			openList[successor.Coords] = successor
			closedList[q.Coords] = q
		}

		if found != (TraversedSquare{}) {
			break
		}

		closedList[q.Coords] = q
	}

	return found, found != (TraversedSquare{})
}

func PartOne() any {
	hmap, start, end := getPartOneData()
	foundPath, found := FindShortestPath(
		hmap,
		start,
		// Prioritize paths closer to the end tile, both on the 2D plane and in elevation
		func(s, q TraversedSquare) int { return ManhattanDist(s.Square, end) + (q.Elevation - s.Elevation) },
		func(s, q TraversedSquare) bool { return s.Elevation <= q.Elevation+1 },
		func(s TraversedSquare) bool { return s.IsEnd },
	)
	if !found {
		panic("no path found from start to end")
	}
	return foundPath.G
}

func PartTwo() any {
	hmap, _, end := getPartOneData()
	// To find the shortest path starting from any 0-elevation square, we instead
	// find the shorest path starting from the end tile and aim towards any arbitrary 0-elevation tile
	foundPath, found := FindShortestPath(
		hmap,
		end,
		// We only care about elevation now, since we're targetting any arbitrary 0-elevation tile
		func(s, q TraversedSquare) int { return s.Elevation },
		func(s, q TraversedSquare) bool { return s.Elevation >= q.Elevation-1 },
		func(s TraversedSquare) bool { return s.Elevation == 0 },
	)
	if !found {
		panic("no path found from start to end")
	}
	return foundPath.G
}

func main() {
	log.Infow(fmt.Sprintf("Answer: %v", PartOne()), "part", 1)
	log.Infow(fmt.Sprintf("Answer: %v", PartTwo()), "part", 2)
}
