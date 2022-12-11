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

//go:embed part1.txt
var files embed.FS

type Monkey struct {
	items            []int
	inspectOperands  []string
	inspectOperation string
	divisibleNum     int
	trueTarget       int
	falseTarget      int
	inspectionCount  int
}

func (m *Monkey) InspectItem() {
	operands := [2]int{}
	for i, operandStr := range m.inspectOperands {
		if operandStr == "old" {
			operands[i] = m.items[0]
		} else {
			operands[i] = util.AtoiNoError(operandStr)
		}
	}

	if m.inspectOperation == "+" {
		m.items[0] = operands[0] + operands[1]
	} else {
		// Assume multiply is the only other allowed
		m.items[0] = operands[0] * operands[1]
	}

	m.inspectionCount++
}

func (m *Monkey) CalcBoredom() {
	m.items[0] = m.items[0] / 3
}

func (m *Monkey) GetNextMonkey() int {
	if m.items[0]%m.divisibleNum == 0 {
		return m.trueTarget
	} else {
		return m.falseTarget
	}
}

func (m *Monkey) Pop() int {
	val := m.items[0]
	m.items = m.items[1:]
	return val
}

func (m *Monkey) Push(item int) {
	m.items = append(m.items, item)
}

func getPartOneData() []Monkey {
	monkeys := []Monkey{}
	var curMonkey Monkey

	for _, line := range util.ReadProblemInput(files, 1) {
		tokens := strings.Split(strings.TrimSpace(line), " ")

		if tokens[0] == "Monkey" {
			// Monkey {index}:
			curMonkey = Monkey{}
			continue
		} else if tokens[0] == "Starting" {
			// Start items: {items...}
			for _, item := range tokens[2:] {
				curMonkey.items = append(
					curMonkey.items,
					util.AtoiNoError(strings.Trim(item, ",")),
				)
			}
		} else if tokens[0] == "Operation:" {
			// Operation: new = {operand} {operator} {operand}
			curMonkey.inspectOperands = []string{tokens[3], tokens[5]}
			curMonkey.inspectOperation = tokens[4]
		} else if tokens[0] == "Test:" {
			// Test: divisible by {divisbleNum}
			curMonkey.divisibleNum = util.AtoiNoError(tokens[3])
		} else if tokens[0] == "If" && tokens[1] == "true:" {
			// If true: throw to monkey {trueTarget}
			curMonkey.trueTarget = util.AtoiNoError(tokens[5])
		} else if tokens[0] == "If" && tokens[1] == "false:" {
			// If false: throw to monkey {falseTarget}
			curMonkey.falseTarget = util.AtoiNoError(tokens[5])
		} else if tokens[0] == "" {
			// Space between monkeys
			monkeys = append(monkeys, curMonkey)
		} else {
			panic(fmt.Sprintf("Unexpected string %v", tokens))
		}
	}
	monkeys = append(monkeys, curMonkey) // DOn't forget last monkey
	return monkeys
}

func PartOne() any {
	monkeys := getPartOneData()

	for i := 0; i < 20; i++ {
		for m := 0; m < len(monkeys); m++ {
			for len(monkeys[m].items) > 0 {
				monkeys[m].InspectItem()
				monkeys[m].CalcBoredom()
				nextM := monkeys[m].GetNextMonkey()
				monkeys[nextM].Push(monkeys[m].Pop())
			}
		}
	}

	topTwoActive := *ds.NewTopList[int](2)
	for _, m := range monkeys {
		topTwoActive.TryPush(m.inspectionCount)
	}

	return topTwoActive.Values[0] * topTwoActive.Values[1]
}

func PartTwo() any {
	monkeys := getPartOneData()

	// We need to reduce using a modulo to avoid integer overflow
	// But reducing item worry by a particular modulo will affect
	// a monkey's divisible check if the modulo does not contain
	// the monkey's divisibleNum as a factor; so we make the
	// the modulo equal to a production of all monkey's divisibleNum
	modReducer := 1
	for _, m := range monkeys {
		modReducer *= m.divisibleNum
	}

	for i := 0; i < 10000; i++ {
		for m := 0; m < len(monkeys); m++ {
			for len(monkeys[m].items) > 0 {
				monkeys[m].InspectItem()
				nextM := monkeys[m].GetNextMonkey()
				monkeys[nextM].Push(monkeys[m].Pop() % modReducer)
			}
		}
	}

	topTwoActive := *ds.NewTopList[int](2)
	for _, m := range monkeys {
		topTwoActive.TryPush(m.inspectionCount)
	}

	return topTwoActive.Values[0] * topTwoActive.Values[1]
}

func main() {
	log.Infow(fmt.Sprintf("Answer: %v", PartOne()), "part", 1)
	log.Infow(fmt.Sprintf("Answer: %v", PartTwo()), "part", 2)
}
