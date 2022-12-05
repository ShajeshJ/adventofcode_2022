package day2

import (
	"embed"

	"github.com/ShajeshJ/adventofcode_2022/common/util"
	"go.uber.org/zap"
)

//go:embed part1.txt
var files embed.FS

func lookupAndSum[K comparable](lookupTable *map[K]int, lookupKeys *[]K) int {
	total := 0
	for _, k := range *lookupKeys {
		total += (*lookupTable)[k]
	}
	return total
}

var p1Scores = map[string]int{
	"A X": 4,
	"B X": 1,
	"C X": 7,
	"A Y": 8,
	"B Y": 5,
	"C Y": 2,
	"A Z": 3,
	"B Z": 9,
	"C Z": 6,
}

func PartOne(logger *zap.SugaredLogger) {
	allRounds := util.ReadProblemInput(files, 1)
	logger.Info(lookupAndSum(&p1Scores, &allRounds))
}

var p2Scores = map[string]int{
	"A X": 3,
	"B X": 1,
	"C X": 2,
	"A Y": 4,
	"B Y": 5,
	"C Y": 6,
	"A Z": 8,
	"B Z": 9,
	"C Z": 7,
}

func PartTwo(logger *zap.SugaredLogger) {
	allRounds := util.ReadProblemInput(files, 1)
	logger.Info(lookupAndSum(&p2Scores, &allRounds))
}
