package main

import (
	"embed"
	"fmt"
	"regexp"

	ds "github.com/ShajeshJ/adventofcode_2022/common/datastructures"
	"github.com/ShajeshJ/adventofcode_2022/common/logging"
	"github.com/ShajeshJ/adventofcode_2022/common/util"
	"golang.org/x/exp/slices"
)

var log = logging.GetLogger()

//go:embed input.txt
var files embed.FS

var stepRegex = regexp.MustCompile(`move (\d+) from (\d+) to (\d+)`)

// Step is a slice of 3 ints, where index
// 0 is the # of containers to move,
// 1 is the stack to move from,
// and 2 is the stack to move to
type Step []int

func (s *Step) GetAmount() int {
	return (*s)[0]
}

func (s *Step) GetFromIndex() int {
	return (*s)[1] - 1
}

func (s *Step) GetToIndex() int {
	return (*s)[2] - 1
}

func getPartOneInput() ([]ds.Stack[rune], []Step) {
	lines := util.ReadProblemInput(files)
	sepIndex := slices.Index(lines, "")

	stacks := make([]ds.Stack[rune], 9)

	for i := sepIndex - 2; i >= 0; i-- {
		for j := 0; j < len(stacks); j++ {
			strRunes := []rune(lines[i])
			if strRunes[j*4+1] == ' ' {
				continue
			}
			stacks[j].Push(strRunes[j*4+1])
		}
	}

	steps := []Step{}
	for i := sepIndex + 1; i < len(lines); i++ {
		var step Step
		for _, m := range stepRegex.FindStringSubmatch(lines[i])[1:] {
			step = append(step, util.AtoiNoError(m))
		}
		steps = append(steps, step)
	}

	return stacks, steps
}

func PartOne() any {
	stacks, steps := getPartOneInput()

	for _, step := range steps {
		for i := 0; i < step.GetAmount(); i++ {
			val, _ := stacks[step.GetFromIndex()].Pop()
			stacks[step.GetToIndex()].Push(val)
		}
	}

	var output string
	for _, s := range stacks {
		top, _ := s.Pop()
		output += string(top)
	}
	return output
}

func PartTwo() any {
	stacks, steps := getPartOneInput()

	for _, step := range steps {
		val, _ := stacks[step.GetFromIndex()].PopN(step.GetAmount())
		stacks[step.GetToIndex()].PushN(val)
	}

	var output string
	for _, s := range stacks {
		top, _ := s.Pop()
		output += string(top)
	}
	return output
}

func main() {
	log.Infow(fmt.Sprintf("%v", PartOne()), "part", 1)
	log.Infow(fmt.Sprintf("%v", PartTwo()), "part", 2)
}
