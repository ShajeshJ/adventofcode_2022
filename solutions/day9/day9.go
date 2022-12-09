package main

import (
	"embed"
	"fmt"
	"math"
	"strings"

	"github.com/ShajeshJ/adventofcode_2022/common/logging"
	"github.com/ShajeshJ/adventofcode_2022/common/util"
)

var log = logging.GetLogger()

//go:embed part1.txt
var files embed.FS

type Position [2]int
type Movement [2]int

func (p *Position) Move(m Movement) {
	(*p)[0] += m[0]
	(*p)[1] += m[1]
}

var moveVec = map[string]Movement{
	"L": {-1, 0},
	"R": {1, 0},
	"U": {0, 1},
	"D": {0, -1},
}

func GetTailMove(posH, posT Position) Movement {
	m := Movement{posH[0] - posT[0], posH[1] - posT[1]}

	if m[0] >= -1 && m[0] <= 1 && m[1] >= -1 && m[1] <= 1 {
		return Movement{0, 0} // Still touching within diagonally; no movement
	}

	// Normalize movement
	if m[0] != 0 {
		m[0] = m[0] / int(math.Abs(float64(m[0])))
	}
	if m[1] != 0 {
		m[1] = m[1] / int(math.Abs(float64(m[1])))
	}

	return m
}

func PartOne() any {
	posH := Position{0, 0}
	posT := Position{0, 0}
	visited := map[Position]int{posT: 1}

	for _, m := range util.ReadProblemInput(files, 1) {
		mParts := strings.Split(m, " ")
		dir, amt := mParts[0], util.AtoiNoError(mParts[1])

		for i := 0; i < amt; i++ {
			posH.Move(moveVec[dir])
			posT.Move(GetTailMove(posH, posT))
			visited[posT]++
		}
	}

	numVisited := 0
	for range visited {
		numVisited++
	}
	return numVisited
}

func PartTwo() any {
	knots := make([]Position, 10)
	visited := map[Position]int{knots[9]: 1}

	for _, m := range util.ReadProblemInput(files, 1) {
		mParts := strings.Split(m, " ")
		dir, amt := mParts[0], util.AtoiNoError(mParts[1])

		for i := 0; i < amt; i++ {
			knots[0].Move(moveVec[dir])

			for j := 1; j < len(knots); j++ {
				knots[j].Move(GetTailMove(knots[j-1], knots[j]))
			}

			visited[knots[len(knots)-1]]++
		}
	}

	numVisited := 0
	for range visited {
		numVisited++
	}
	return numVisited
}

func main() {
	log.Infow(fmt.Sprintf("Answer: %v", PartOne()), "part", 1)
	log.Infow(fmt.Sprintf("Answer: %v", PartTwo()), "part", 2)
}
