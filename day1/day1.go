package day1

import (
	"embed"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"golang.org/x/exp/constraints"
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

	mostCalories := 0
	curElfCalories := 0

	for _, food := range allFoods {
		if food == "" {
			if curElfCalories > mostCalories {
				mostCalories = curElfCalories
			}
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

	logger.Info(mostCalories)
}

type TopList[T constraints.Ordered] struct {
	Size int
	data []T
}

func NewTopList[T constraints.Ordered](size int) *TopList[T] {
	if size < 1 {
		return &TopList[T]{Size: 1, data: make([]T, 1)}
	}
	return &TopList[T]{Size: size, data: make([]T, size)}
}

// Add will add `val` to the list if it's greater than at least 1 number in the list in order, and returns true.
// Otherwise it does nothing and returns false
func (t *TopList[T]) Add(val T) bool {
	if val <= t.data[t.Size-1] {
		return false
	}

	t.data[t.Size-1] = val

	for i := t.Size - 2; i >= 0; i-- {
		if t.data[i] >= t.data[i+1] {
			break
		}
		t.data[i], t.data[i+1] = t.data[i+1], t.data[i]
	}

	return true
}

// Sum will return the sum of all values in the list
func (t *TopList[T]) Sum() (total T) {
	for _, item := range t.data {
		total += item
	}
	return
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
			mostCalories.Add(curElfCalories)
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
