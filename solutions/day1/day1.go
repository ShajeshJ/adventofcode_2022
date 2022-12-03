package day1

import (
	"embed"
	"strconv"
	"strings"

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
	allFoods, err := getPartOneData()
	if err != nil {
		logger.Error(err)
		return
	}

	mostCalories := NewTopList[int](1)
	curElfCalories := 0

	for _, food := range allFoods {
		if food == "" {
			mostCalories.Push(curElfCalories)
			curElfCalories = 0
			continue
		}

		calories, err := strconv.Atoi(food)
		if err != nil {
			logger.Errorw(err.Error(), "calories", calories)
			return
		}
		curElfCalories += calories
	}

	logger.Info(mostCalories.Sum())
}

func PartTwo(logger *zap.SugaredLogger) {
	allFoods, err := getPartOneData()
	if err != nil {
		logger.Error(err)
		return
	}

	mostCalories := NewTopList[int](3)
	curElfCalories := 0

	for _, food := range allFoods {
		if food == "" {
			mostCalories.Push(curElfCalories)
			curElfCalories = 0
			continue
		}

		calories, err := strconv.Atoi(food)
		if err != nil {
			logger.Errorw(err.Error(), "calories", calories)
			return
		}
		curElfCalories += calories
	}

	logger.Info(mostCalories.Sum())
}
