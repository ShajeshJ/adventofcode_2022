package main

import (
	"embed"
	"fmt"

	"github.com/ShajeshJ/adventofcode_2022/common/logging"
)

var log = logging.GetLogger()

//go:embed input.txt
var files embed.FS

func PartOne() any {
	return 0
}

func PartTwo() any {
	return 0
}

func main() {
	log.Infow(fmt.Sprintf("Answer: %v", PartOne()), "part", 1)
	log.Infow(fmt.Sprintf("Answer: %v", PartTwo()), "part", 2)
}
