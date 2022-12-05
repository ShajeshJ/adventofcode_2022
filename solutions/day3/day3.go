package day3

import (
	"embed"

	"github.com/ShajeshJ/adventofcode_2022/common/util"
	"go.uber.org/zap"
)

//go:embed part1.txt
var files embed.FS

func getPartTwoData() (data [][]string) {
	rucksacks := util.ReadProblemInput(files, 1)

	for i := 0; i < len(rucksacks); i += 3 {
		data = append(data, []string{rucksacks[i], rucksacks[i+1], rucksacks[i+2]})
	}

	return
}

func getCommonLetter(strs ...string) rune {
	overlapCount := map[rune]int{}

	for _, str := range strs[:len(strs)-1] {
		for _, r := range str {
			overlapCount[r]++
		}
	}

	for _, r := range strs[len(strs)-1] {
		if overlapCount[r] == len(strs)-1 {
			return r
		}
	}

	return 0 // Shouldn't happen
}

func getPriority(r rune) int {
	if 'a' <= r && r <= 'z' {
		return int(r - 'a' + 1)
	} else { // Should only catch: 'A' <= r && r <= 'Z'
		return int(r - 'A' + 27)
	}
}

func PartOne(logger *zap.SugaredLogger) {
	rucksacks := util.ReadProblemInput(files, 1)

	total := 0

	for _, items := range rucksacks {
		rucksackCompartments := []string{items[:len(items)/2], items[len(items)/2:]}
		total += getPriority(getCommonLetter(rucksackCompartments...))
	}

	logger.Info(total)
}

func PartTwo(logger *zap.SugaredLogger) {
	groups := getPartTwoData()

	total := 0

	for _, group := range groups {
		total += getPriority(getCommonLetter(group...))
	}

	logger.Info(total)
}
