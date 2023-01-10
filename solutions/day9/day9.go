package main

import (
	"embed"
	"fmt"
	"strings"

	"github.com/ShajeshJ/adventofcode_2022/common/logging"
	"github.com/ShajeshJ/adventofcode_2022/common/util"
)

var log = logging.GetLogger()

//go:embed input.txt
var files embed.FS

type Position [2]int
type Movement [2]int

func (p *Position) Move(m Movement) {
	(*p)[0] += m[0]
	(*p)[1] += m[1]
}

var moveUnitVec = map[string]Movement{
	"L": {-1, 0},
	"R": {1, 0},
	"U": {0, 1},
	"D": {0, -1},
}

func GetKnotMove(prevKnot, nextKnot Position) Movement {
	m := Movement{prevKnot[0] - nextKnot[0], prevKnot[1] - nextKnot[1]}

	if m[0] >= -1 && m[0] <= 1 && m[1] >= -1 && m[1] <= 1 {
		return Movement{0, 0} // Still touching within diagonal distance
	}

	if m[0] != 0 {
		m[0] = util.Normalize(m[0])
	}
	if m[1] != 0 {
		m[1] = util.Normalize(m[1])
	}

	return m
}

func SimulateRope(numKnots int) int {
	knots := make([]Position, numKnots)
	visited := map[Position]int{knots[numKnots-1]: 1}

	for _, m := range util.ReadProblemInput(files) {
		headMoves := strings.Split(m, " ")
		dir, amt := headMoves[0], util.AtoiNoError(headMoves[1])

		for i := 0; i < amt; i++ {
			knots[0].Move(moveUnitVec[dir])

			for j := 1; j < len(knots); j++ {
				knots[j].Move(GetKnotMove(knots[j-1], knots[j]))
			}

			visited[knots[numKnots-1]]++
		}
	}

	return len(visited)
}

func PartOne() any {
	return SimulateRope(2)
}

func PartTwo() any {
	return SimulateRope(10)
}

func main() {
	log.Infow(fmt.Sprintf("Answer: %v", PartOne()), "part", 1)
	log.Infow(fmt.Sprintf("Answer: %v", PartTwo()), "part", 2)
}
