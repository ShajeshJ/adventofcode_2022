package main

import (
	"embed"
	"fmt"
	"math"

	"github.com/ShajeshJ/adventofcode_2022/common/logging"
	"github.com/ShajeshJ/adventofcode_2022/common/util"
)

var log = logging.GetLogger()

//go:embed part1.txt
var files embed.FS

type SNAFUDIGIT rune

var DIGITVAL = map[SNAFUDIGIT]int{
	'=': -2,
	'-': -1,
	'0': 0,
	'1': 1,
	'2': 2,
}

type Snafu []SNAFUDIGIT

func (s *Snafu) Int() int {
	val := 0
	for i, r := range *s {
		val += int(math.Pow(5, float64(i))) * DIGITVAL[r]
	}
	return val
}

func (s *Snafu) String() string {
	str := ""
	for i := len(*s) - 1; i >= 0; i-- {
		str += string((*s)[i])
	}
	return str
}

func ConvertToSnafu(n int) Snafu {
	snafu := make(Snafu, 0)
	carry := 0
	for n > 0 {
		rem := n % 5
		if rem > 2 {
			rem = rem - 5
			carry = 1
		} else {
			carry = 0
		}

		for k, v := range DIGITVAL {
			if v == rem {
				snafu = append(snafu, k)
			}
		}
		n = n/5 + carry
	}
	if carry == 1 {
		snafu = append(snafu, SNAFUDIGIT('1'))
	}
	return snafu
}

func getPartOneData() []Snafu {
	snafus := make([]Snafu, 0)
	for _, line := range util.ReadProblemInput(files, 1) {
		snafu := make(Snafu, 0)
		for j := len(line) - 1; j >= 0; j-- {
			snafu = append(snafu, SNAFUDIGIT(line[j]))
		}
		snafus = append(snafus, snafu)
	}
	return snafus
}

func PartOne() any {
	total := 0
	for _, snafu := range getPartOneData() {
		total += snafu.Int()
	}
	totalAsSnafu := ConvertToSnafu(total)
	return totalAsSnafu.String()
}

func PartTwo() any {
	return "No puzzle for part 2"
}

func main() {
	log.Infow(fmt.Sprintf("Answer: %v", PartOne()), "part", 1)
	log.Infow(fmt.Sprintf("Answer: %v", PartTwo()), "part", 2)
}
