package main

import (
	"embed"
	"fmt"

	"github.com/ShajeshJ/adventofcode_2022/common/logging"
	"github.com/ShajeshJ/adventofcode_2022/common/util"
	"golang.org/x/exp/slices"
)

var log = logging.GetLogger()

//go:embed input.txt
var files embed.FS

type SNAFUDIGIT rune

var DIGITS = []SNAFUDIGIT{'=', '-', '0', '1', '2'}

const OFFSET = 2 // The negative offset of a digit's actual value to its index in DIGITS

type Snafu []SNAFUDIGIT

func (s *Snafu) Int() int {
	val := 0
	for i, r := range *s {
		val += util.Pow(5, i) * (slices.Index(DIGITS, r) - OFFSET)
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

func GetNextDigit(n int) (SNAFUDIGIT, int) {
	rem := n % 5
	carry := 0

	if rem > 2 {
		rem = rem - 5
		carry = 1
	}

	return DIGITS[rem+OFFSET], carry
}

func ConvertToSnafu(n int) Snafu {
	snafu := make(Snafu, 0)
	for n > 0 {
		digit, carry := GetNextDigit(n)
		snafu = append(snafu, digit)
		n = n/5 + carry
	}
	return snafu
}

func getPartOneData() []Snafu {
	snafus := make([]Snafu, 0)
	for _, line := range util.ReadProblemInput(files) {
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
