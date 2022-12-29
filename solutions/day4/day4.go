package main

import (
	"embed"
	"fmt"
	"regexp"

	"github.com/ShajeshJ/adventofcode_2022/common/logging"
	"github.com/ShajeshJ/adventofcode_2022/common/util"
)

var log = logging.GetLogger()

//go:embed input.txt
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
			p = append(p, util.AtoiNoError(m))
		}
		elfPairs = append(elfPairs, []Elf{{p[0], p[1]}, {p[2], p[3]}})
	}
	return
}

func PartOne() any {
	elfPairs := getPartOneData()

	total := 0

	for _, pair := range elfPairs {
		if pair[0].Contains(pair[1]) || pair[1].Contains(pair[0]) {
			total += 1
		}
	}

	return total
}

func PartTwo() any {
	elfPairs := getPartOneData()

	total := 0

	for _, pair := range elfPairs {
		if pair[0].Intersects(pair[1]) {
			total += 1
		}
	}

	return total
}

func main() {
	log.Infow(fmt.Sprintf("Answer: %v", PartOne()), "part", 1)
	log.Infow(fmt.Sprintf("Answer: %v", PartTwo()), "part", 2)
}
