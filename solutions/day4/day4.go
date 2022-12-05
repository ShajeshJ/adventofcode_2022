package day4

import (
	"embed"
	"regexp"
	"strconv"

	"github.com/ShajeshJ/adventofcode_2022/common/util"
	"go.uber.org/zap"
)

//go:embed part1.txt
var files embed.FS

type Elf struct {
	start int
	end   int
}

func (e *Elf) Contains(other Elf) bool {
	return e.start >= other.start && e.end <= other.end
}

func (e *Elf) Intersects(other Elf) bool {
	return (e.end >= other.start && e.start <= other.end ||
		other.end >= e.start && other.start <= e.end)
}

var inputRegex = regexp.MustCompile(`(\d+)-(\d+),(\d+)-(\d+)`)

func getPartOneData() (elfPairs [][]Elf) {
	for _, line := range util.ReadProblemInput(files, 1) {
		var p []int
		for _, m := range inputRegex.FindStringSubmatch(line)[1:] {
			temp, _ := strconv.Atoi(m)
			p = append(p, temp)
		}
		elfPairs = append(elfPairs, []Elf{{p[0], p[1]}, {p[2], p[3]}})
	}
	return
}

func PartOne(logger *zap.SugaredLogger) {
	elfPairs := getPartOneData()

	total := 0

	for _, pair := range elfPairs {
		if pair[0].Contains(pair[1]) || pair[1].Contains(pair[0]) {
			total += 1
		}
	}

	logger.Info(total)
}

func PartTwo(logger *zap.SugaredLogger) {
	elfPairs := getPartOneData()

	total := 0

	for _, pair := range elfPairs {
		if pair[0].Intersects(pair[1]) {
			total += 1
		}
	}

	logger.Info(total)
}
