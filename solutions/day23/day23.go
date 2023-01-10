package main

import (
	"embed"
	"fmt"

	"github.com/ShajeshJ/adventofcode_2022/common/logging"
	"github.com/ShajeshJ/adventofcode_2022/common/util"
	"golang.org/x/exp/maps"
)

var log = logging.GetLogger()

type Direction int

const (
	N = Direction(iota)
	S
	W
	E
)

type Point struct {
	X int
	Y int
}

type Elf struct {
	Pos Point
}

func GetNextCardinal(d Direction) Direction {
	return Direction((d + 1) % 4)
}

type ElfMap map[Point]Elf

func (em *ElfMap) GetAdjacent(p Point) map[Direction]int {
	adj := map[Direction]int{}
	if _, ok := (*em)[Point{p.X, p.Y - 1}]; ok { // N
		adj[N]++
	}
	if _, ok := (*em)[Point{p.X, p.Y + 1}]; ok { // S
		adj[S]++
	}
	if _, ok := (*em)[Point{p.X - 1, p.Y}]; ok { // W
		adj[W]++
	}
	if _, ok := (*em)[Point{p.X + 1, p.Y}]; ok { // E
		adj[E]++
	}
	if _, ok := (*em)[Point{p.X - 1, p.Y - 1}]; ok { // NW
		adj[N]++
		adj[W]++
	}
	if _, ok := (*em)[Point{p.X + 1, p.Y - 1}]; ok { // NE
		adj[N]++
		adj[E]++
	}
	if _, ok := (*em)[Point{p.X - 1, p.Y + 1}]; ok { // SW
		adj[S]++
		adj[W]++
	}
	if _, ok := (*em)[Point{p.X + 1, p.Y + 1}]; ok { // SE
		adj[S]++
		adj[E]++
	}
	return adj
}

func (em *ElfMap) GetArea() int {
	points := maps.Keys(*em)
	minX := points[0].X
	maxX := points[0].X
	minY := points[0].Y
	maxY := points[0].Y

	for _, p := range points[1:] {
		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}

	return (maxX - minX + 1) * (maxY - minY + 1)
}

func GetPoint(d Direction, p Point) Point {
	switch d {
	case N:
		return Point{p.X, p.Y - 1}
	case S:
		return Point{p.X, p.Y + 1}
	case W:
		return Point{p.X - 1, p.Y}
	case E:
		return Point{p.X + 1, p.Y}
	}
	panic("invalid direction")
}

func getPartOneData() ElfMap {
	elves := make(ElfMap)
	for i, line := range util.ReadProblemInput(files) {
		for j, c := range line {
			if c == '#' {
				p := Point{j, i}
				elves[p] = Elf{p}
			}
		}
	}
	return elves
}

//go:embed input.txt
var files embed.FS

func getProposals(elves ElfMap, curDir Direction) map[Point][]Elf {
	proposals := map[Point][]Elf{}

	for p := range elves {
		adj := elves.GetAdjacent(p)
		if len(adj) == 0 {
			continue
		}

		canMove := false
		proposalDir := curDir
		for i := 0; i < 4; i++ {
			if adj[proposalDir] == 0 {
				canMove = true
				break
			}
			proposalDir = GetNextCardinal(proposalDir)
		}
		if !canMove {
			continue
		}

		nextPoint := GetPoint(proposalDir, p)
		if _, ok := proposals[nextPoint]; !ok {
			proposals[nextPoint] = []Elf{}
		}
		proposals[nextPoint] = append(proposals[nextPoint], elves[p])
	}

	return proposals
}

func SimulateCoordination(maxRound int) int {
	elves := getPartOneData()
	curDir := N
	curRound := 0

	for maxRound <= 0 || curRound < maxRound {
		curRound++
		proposals := getProposals(elves, curDir)

		elfMoved := false
		for p, elvesToMove := range proposals {
			if len(elvesToMove) != 1 {
				continue
			}
			elfMoved = true
			elf := elvesToMove[0]
			delete(elves, elf.Pos)
			elf.Pos = p
			elves[p] = elf
		}

		if !elfMoved {
			break
		}

		curDir = GetNextCardinal(curDir)
	}

	if maxRound <= 0 {
		return curRound
	}
	return elves.GetArea() - len(elves)
}

func PartOne() any {
	return SimulateCoordination(10)
}

func PartTwo() any {
	return SimulateCoordination(-1)
}

func main() {
	log.Infow(fmt.Sprintf("Answer: %v", PartOne()), "part", 1)
	log.Infow(fmt.Sprintf("Answer: %v", PartTwo()), "part", 2)
}
