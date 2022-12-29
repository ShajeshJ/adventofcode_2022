package main

import (
	"embed"
	"fmt"

	ds "github.com/ShajeshJ/adventofcode_2022/common/datastructures"
	"github.com/ShajeshJ/adventofcode_2022/common/logging"
	"github.com/ShajeshJ/adventofcode_2022/common/util"
	"golang.org/x/exp/constraints"
)

var log = logging.GetLogger()

//go:embed input.txt
var files embed.FS

func getSum[T constraints.Ordered](l *ds.TopList[T]) T {
	var total T
	for _, item := range l.Values {
		total += item
	}
	return total
}

func PartOne() any {
	mostCalories := ds.NewTopList[int](1)
	curElfCalories := 0

	for _, food := range util.ReadProblemInput(files, 1) {
		if food == "" {
			mostCalories.TryPush(curElfCalories)
			curElfCalories = 0
			continue
		}
		curElfCalories += util.AtoiNoError(food)
	}

	return getSum(mostCalories)
}

func PartTwo() any {
	mostCalories := ds.NewTopList[int](3)
	curElfCalories := 0

	for _, food := range util.ReadProblemInput(files, 1) {
		if food == "" {
			mostCalories.TryPush(curElfCalories)
			curElfCalories = 0
			continue
		}
		curElfCalories += util.AtoiNoError(food)
	}

	return getSum(mostCalories)
}

func main() {
	log.Infow(fmt.Sprintf("%v", PartOne()), "part", 1)
	log.Infow(fmt.Sprintf("%v", PartTwo()), "part", 2)
}
