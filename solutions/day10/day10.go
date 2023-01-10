package main

import (
	"embed"
	"fmt"
	"strings"

	ds "github.com/ShajeshJ/adventofcode_2022/common/datastructures"
	"github.com/ShajeshJ/adventofcode_2022/common/logging"
	"github.com/ShajeshJ/adventofcode_2022/common/util"
)

var log = logging.GetLogger()

//go:embed input.txt
var files embed.FS

func GetDelay(instruction string) int {
	if instruction == "noop" {
		return 1
	} else {
		return 2
	}
}

func RunCRT(doCycleProcessing func(cycle, x int)) {
	x := 1
	cycle := 0
	delay := 0
	curInstruction := ""
	instructions := ds.Stack[string](util.ReadProblemInput(files))

	for len(instructions) != 0 {
		cycle++
		doCycleProcessing(cycle, x)

		if curInstruction == "" {
			curInstruction, _ = instructions.Pop()
			delay = GetDelay(curInstruction)
		}

		delay--
		if delay > 0 {
			continue
		}

		if curInstruction != "noop" {
			x += util.AtoiNoError(strings.Split(curInstruction, " ")[1])
		}
		curInstruction = ""
	}
}

func PartOne() any {
	total := 0
	RunCRT(func(cycle, x int) {
		if (cycle-20)%40 == 0 {
			total += cycle * x
		}
	})
	return total
}

func PartTwo() any {
	screen := ""

	RunCRT(func(cycle, x int) {
		if cycle%40 == 1 {
			screen += "\n"
		}

		curPixel := (cycle - 1) % 40
		if curPixel >= x-1 && curPixel <= x+1 {
			screen += "#"
		} else {
			screen += "."
		}
	})

	return screen
}

func main() {
	log.Infow(fmt.Sprintf("Answer: %v", PartOne()), "part", 1)
	log.Infow(fmt.Sprintf("Answer: %v", PartTwo()), "part", 2)
}
