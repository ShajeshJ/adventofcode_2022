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

func hasDuplicates(counter map[rune]int) bool {
	for _, count := range counter {
		if count > 1 {
			return true
		}
	}
	return false
}

func PartOne(n int) any {
	seq := util.ReadProblemInput(files, 1)[0]

	buffer := []rune{}
	runeCount := map[rune]int{}

	// Add first n-1 chars to reach the required length
	for _, c := range seq[:n-1] {
		buffer = append(buffer, c)
		runeCount[c]++
	}

	counter := n - 1

	for _, c := range seq[n-1:] {
		buffer = append(buffer, c)
		runeCount[c]++
		counter++

		if !hasDuplicates(runeCount) {
			break
		}

		runeCount[buffer[0]]--
		buffer = buffer[1:]
	}

	return counter
}

func PartTwo() any {
	return PartOne(14)
}

func main() {
	log.Infow(fmt.Sprintf("Answer: %v", PartOne(4)), "part", 1)
	log.Infow(fmt.Sprintf("Answer: %v", PartTwo()), "part", 2)
}
