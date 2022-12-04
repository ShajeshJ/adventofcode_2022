package day4

import (
	"embed"
	"regexp"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

//go:embed part1.txt
var files embed.FS

type Elf struct {
	start int
	end   int
}

func (e *Elf) IsRangeSubset(other Elf) bool {
	return e.start >= other.start && e.end <= other.end
}

func (e *Elf) HasOverlap(other Elf) bool {
	return (e.end >= other.start && e.start <= other.end ||
		other.end >= e.start && other.start <= e.end)
}

var elfPairRegex = `(\d+)-(\d+),(\d+)-(\d+)`

func getPartOneData() (elfPairs [][]Elf) {
	bytes, _ := files.ReadFile("part1.txt")
	for _, pairStr := range strings.Split(string(bytes), "\n") {
		r := regexp.MustCompile(elfPairRegex)
		matches := r.FindStringSubmatch(pairStr)

		leftStart, _ := strconv.Atoi(matches[1])
		leftEnd, _ := strconv.Atoi(matches[2])
		rightStart, _ := strconv.Atoi(matches[3])
		rightEnd, _ := strconv.Atoi(matches[4])

		elfPairs = append(elfPairs, []Elf{{leftStart, leftEnd}, {rightStart, rightEnd}})
	}
	return
}

func PartOne(logger *zap.SugaredLogger) {
	elfPairs := getPartOneData()

	total := 0

	for _, pair := range elfPairs {
		if pair[0].IsRangeSubset(pair[1]) || pair[1].IsRangeSubset(pair[0]) {
			total += 1
		}
	}

	logger.Info(total)
}

func PartTwo(logger *zap.SugaredLogger) {
	elfPairs := getPartOneData()

	total := 0

	for _, pair := range elfPairs {
		if pair[0].HasOverlap(pair[1]) {
			total += 1
		}
	}

	logger.Info(total)
}
