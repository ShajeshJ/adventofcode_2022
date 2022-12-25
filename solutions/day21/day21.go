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

const (
	Add = iota
	Subtract
	Multiply
	Divide
)

type Operator int

type Monkey struct {
	IsVal     bool
	IsUnknown bool
	Val       int
	Left      string
	Right     string
	Op        Operator
}

func getPartOneData() map[string]*Monkey {
	OpLookup := map[string]Operator{
		"+": Add,
		"-": Subtract,
		"*": Multiply,
		"/": Divide,
	}

	monkeys := make(map[string]*Monkey)
	for _, line := range util.ReadProblemInput(files, 1) {
		tokens := strings.Split(line, " ")
		name := strings.TrimRight(tokens[0], ":")
		if len(tokens) == 2 {
			monkeys[name] = &Monkey{IsVal: true, Val: util.AtoiNoError(tokens[1])}
		} else {
			monkeys[name] = &Monkey{IsVal: false, Left: tokens[1], Right: tokens[3], Op: OpLookup[tokens[2]]}
		}
	}

	return monkeys
}

func Compute(m *Monkey, monkeys map[string]*Monkey) int {
	if m.IsVal {
		return m.Val
	}

	left := monkeys[m.Left]
	right := monkeys[m.Right]

	switch m.Op {
	case Add:
		return Compute(left, monkeys) + Compute(right, monkeys)
	case Subtract:
		return Compute(left, monkeys) - Compute(right, monkeys)
	case Multiply:
		return Compute(left, monkeys) * Compute(right, monkeys)
	case Divide:
		return Compute(left, monkeys) / Compute(right, monkeys)
	}

	panic("Invalid operator")
}

func PartOne() any {
	monkeys := getPartOneData()
	return Compute(monkeys["root"], monkeys)
}

type PartialResult struct {
	IsUnknown bool
	Val       int
	Left      *PartialResult
	Op        Operator
	Right     *PartialResult
}

func ComputePartial(m *Monkey, monkeys map[string]*Monkey) *PartialResult {
	if m.IsUnknown {
		return &PartialResult{IsUnknown: true}
	}
	if m.IsVal {
		return &PartialResult{Val: m.Val}
	}

	result := &PartialResult{}
	result.Left = ComputePartial(monkeys[m.Left], monkeys)
	result.Right = ComputePartial(monkeys[m.Right], monkeys)
	result.Op = m.Op
	result.IsUnknown = result.Left.IsUnknown || result.Right.IsUnknown

	if result.IsUnknown {
		return result
	}

	switch m.Op {
	case Add:
		result.Val = result.Left.Val + result.Right.Val
	case Subtract:
		result.Val = result.Left.Val - result.Right.Val
	case Multiply:
		result.Val = result.Left.Val * result.Right.Val
	case Divide:
		result.Val = result.Left.Val / result.Right.Val
	}

	return result
}

func SolveStep(target int, equation *PartialResult) int {
	if equation.Left == nil && equation.Right == nil {
		return target
	}

	var unknown *PartialResult

	if equation.Left.IsUnknown {
		unknown = equation.Left
		switch equation.Op {
		case Add:
			target -= equation.Right.Val
		case Subtract:
			target += equation.Right.Val
		case Multiply:
			target /= equation.Right.Val
		case Divide:
			target *= equation.Right.Val
		}
	} else {
		unknown = equation.Right
		switch equation.Op {
		case Add:
			target -= equation.Left.Val
		case Subtract:
			target = equation.Left.Val - target
		case Multiply:
			target /= equation.Left.Val
		case Divide:
			target = equation.Left.Val / target
		}
	}

	return SolveStep(target, unknown)
}

// SolveEquation solves by assuming there is exactly 1 deeply nested unknown to unwrap solve naively
func SolveEquation(left, right *PartialResult) int {
	var target, unknown *PartialResult
	if left.IsUnknown {
		unknown = left
		target = right
	} else {
		unknown = right
		target = left
	}

	return SolveStep(target.Val, unknown)
}

func PartTwo() any {
	monkeys := getPartOneData()
	monkeys["humn"].IsUnknown = true

	left := ComputePartial(monkeys[monkeys["root"].Left], monkeys)
	right := ComputePartial(monkeys[monkeys["root"].Right], monkeys)

	return SolveEquation(left, right)
}

func main() {
	log.Infow(fmt.Sprintf("Answer: %v", PartOne()), "part", 1)
	log.Infow(fmt.Sprintf("Answer: %v", PartTwo()), "part", 2)
}
