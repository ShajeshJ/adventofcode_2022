package main

import (
	"embed"
	"fmt"
	"strings"

	"github.com/ShajeshJ/adventofcode_2022/common/logging"
	"github.com/ShajeshJ/adventofcode_2022/common/util"
)

var log = logging.GetLogger()

//go:embed part1.txt
var files embed.FS

func GetDelay(instruction string) int {
	if instruction == "noop" {
		return 1
	} else {
		return 2
	}
}

func Pop(instructions []string) (string, []string) {
	return instructions[0], instructions[1:]
}

func PartOne() any {
	instructions := util.ReadProblemInput(files, 1)
	x := 1
	cycle := 0
	delay := 0
	total := 0
	curInstruction := ""

	for len(instructions) != 0 {
		cycle++
		if (cycle-20)%40 == 0 {
			fmt.Println(cycle)
			total += cycle * x
		}

		if curInstruction == "" {
			curInstruction, instructions = Pop(instructions)
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

	return total
}

func PartTwo() any {
	instructions := util.ReadProblemInput(files, 1)
	x := 1
	cycle := 0
	delay := 0
	curInstruction := ""
	screen := ""

	for len(instructions) != 0 {
		cycle++
		if cycle%40 == 1 {
			screen += "\n"
		}

		if cycle%40-1 >= x-1 && cycle%40-1 <= x+1 {
			screen += "#"
		} else {
			screen += "."
		}

		if curInstruction == "" {
			curInstruction, instructions = Pop(instructions)
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

	return screen
}

func main() {
	log.Infow(fmt.Sprintf("Answer: %v", PartOne()), "part", 1)
	log.Infow(fmt.Sprintf("Answer: %v", PartTwo()), "part", 2)
}
