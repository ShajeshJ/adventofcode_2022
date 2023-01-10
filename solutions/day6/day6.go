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

func hasDuplicateRunes(s string) bool {
	for _, c := range s {
		if strings.Count(s, string(c)) > 1 {
			return true
		}
	}
	return false
}

func PartOne(n int) any {
	seq := util.ReadProblemInput(files)[0]

	// Must be a minimum of n characters long
	buffer := seq[:n]
	counter := n

	for _, c := range seq[n:] {
		if !hasDuplicateRunes(buffer) {
			break
		}
		buffer = buffer[1:]
		buffer += string(c)
		counter++
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
