package day1

import (
	"embed"
	"strconv"
	"strings"

	"github.com/ShajeshJ/adventofcode_2022/common/utility"
	"go.uber.org/zap"
)

//go:embed part1.txt
var files embed.FS

func getPartOneData() (data []string, err error) {
	bytes, err := files.ReadFile("part1.txt")
	if err != nil {
		return
	}

	data = strings.Split(string(bytes), "\n")
	return
}

func PartOne(logger *zap.SugaredLogger) {
	mostCalories := NewTopList[int](1)
	curElfCalories := 0

	for _, food := range utility.ReadProblemInput(files, 1) {
		if food == "" {
			mostCalories.TryPush(curElfCalories)
			curElfCalories = 0
			continue
		}

		calories, _ := strconv.Atoi(food)
		curElfCalories += calories
	}

	logger.Info(mostCalories.Sum())
}

func PartTwo(logger *zap.SugaredLogger) {
	mostCalories := NewTopList[int](3)
	curElfCalories := 0

	for _, food := range utility.ReadProblemInput(files, 1) {
		if food == "" {
			mostCalories.TryPush(curElfCalories)
			curElfCalories = 0
			continue
		}

		calories, _ := strconv.Atoi(food)
		curElfCalories += calories
	}

	logger.Info(mostCalories.Sum())
}
