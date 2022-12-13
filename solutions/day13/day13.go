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

type Packet []any

type PacketParser struct {
	Tokens []rune
}

func (p *PacketParser) Next() rune {
	val := p.Tokens[0]
	p.Tokens = p.Tokens[1:]
	return val
}

func (p *PacketParser) ParseExpression() any {
	if p.Tokens[0] == '[' {
		return p.ParseList()
	} else {
		return p.ParseValue()
	}
}

func (p *PacketParser) ParseValue() int {
	builder := ""
	for p.Tokens[0] >= '0' && p.Tokens[0] <= '9' {
		builder += string(p.Next()) // pop digit and append to builder
	}
	return util.AtoiNoError(builder)
}

func (p *PacketParser) ParseList() []any {
	val := []any{}

	p.Next() // pop '['

	for p.Tokens[0] != ']' {
		val = append(val, p.ParseExpression())
		if p.Tokens[0] == ',' {
			p.Next()
		}
	}

	p.Next() // pop ']'

	return val
}

func getPartOneData() [][2]Packet {
	var receivedPairs [][2]Packet
	var curPair [2]Packet
	fillIndex := 0

	inputLines := util.ReadProblemInput(files, 1)
	if inputLines[len(inputLines)-1] != "" {
		// Simplify for-loop parsing
		inputLines = append(inputLines, "")
	}

	for _, line := range inputLines {
		if line == "" {
			receivedPairs = append(receivedPairs, curPair)
			curPair = [2]Packet{}
			fillIndex = 0
			continue
		}

		p := PacketParser{[]rune(line)}
		curPair[fillIndex] = p.ParseList()
		fillIndex++
	}

	return receivedPairs
}

type CompareResult int

const (
	Invalid = iota
	Inconclusive
	Valid
)

func CompareItems(left any, right any) CompareResult {
	leftInt, leftIsInt := left.(int)
	rightInt, rightIsInt := right.(int)

	if leftIsInt && rightIsInt {
		return CompareValues(leftInt, rightInt)
	}

	leftSlice := []any{leftInt}
	rightSlice := []any{rightInt}

	if !leftIsInt {
		leftSlice = left.([]any)
	}

	if !rightIsInt {
		rightSlice = right.([]any)
	}

	return CompareLists(leftSlice, rightSlice)
}

func CompareValues(left int, right int) CompareResult {
	if left < right {
		return Valid
	} else if left > right {
		return Invalid
	} else {
		return Inconclusive
	}
}

func CompareLists(left []any, right []any) CompareResult {
	for i := 0; i < util.Min(len(left), len(right)); i++ {
		if result := CompareItems(left[i], right[i]); result != Inconclusive {
			return result
		}
	}

	if len(left) < len(right) {
		return Valid
	} else if len(left) > len(right) {
		return Invalid
	} else {
		return Inconclusive
	}
}

func PartOne() any {
	total := 0
	for i, pair := range getPartOneData() {
		if result := CompareLists(pair[0], pair[1]); result == Valid {
			total += i + 1
		}
	}
	return total
}

func getPartTwoData() []Packet {
	var receivedPairs []Packet

	for _, line := range util.ReadProblemInput(files, 1) {
		if line == "" {
			continue
		}

		p := PacketParser{[]rune(line)}
		receivedPairs = append(receivedPairs, p.ParseList())
	}

	return receivedPairs
}

func FindPacketIndex(packet Packet, others []Packet) int {
	index := 1
	for _, other := range others {
		if result := CompareLists(other, packet); result == Valid {
			index++
		}
	}
	return index
}

func PartTwo() any {
	divPackets := []Packet{
		{[]any{2}},
		{[]any{6}},
	}

	packets := getPartTwoData()

	decoderKey := FindPacketIndex(divPackets[0], append(packets, divPackets[1]))
	decoderKey *= FindPacketIndex(divPackets[1], append(packets, divPackets[0]))

	return decoderKey
}

func main() {
	log.Infow(fmt.Sprintf("Answer: %v", PartOne()), "part", 1)
	log.Infow(fmt.Sprintf("Answer: %v", PartTwo()), "part", 2)
}
