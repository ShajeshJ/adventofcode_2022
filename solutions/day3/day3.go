package day3

import (
	"embed"
	"errors"
	"fmt"
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

func getPartTwoData() (data [][]string, err error) {
	rucksacks, err := getPartOneData()
	if err != nil {
		return
	}

	if len(rucksacks)%3 != 0 {
		err = errors.New("cannot evenly divide elves into groups of 3")
		return
	}

	for i := 0; i < len(rucksacks); i += 3 {
		data = append(data, []string{rucksacks[i], rucksacks[i+1], rucksacks[i+2]})
	}

	return
}

func getCommonLetter(strs ...string) (rune, error) {
	if len(strs) <= 1 {
		return 0, errors.New("must give > 1 strings to find the common letter")
	}

	runningOverlaps := map[rune]bool{}

	for _, r := range strs[0] {
		runningOverlaps[r] = true
	}

	for _, str := range strs[1:] {
		prevOverlaps := runningOverlaps
		runningOverlaps = map[rune]bool{}

		for _, r := range str {
			if _, ok := prevOverlaps[r]; ok {
				runningOverlaps[r] = true
			}
		}
	}

	var overlaps []string
	for r := range runningOverlaps {
		overlaps = append(overlaps, string(r))
	}

	if len(overlaps) == 0 {
		return 0, fmt.Errorf(`no overlap found across strings %v`, strs)
	}

	if len(overlaps) > 1 {
		return 0, fmt.Errorf(`multiple overlaps %v found in %v`, overlaps, strs)
	}

	return rune(overlaps[0][0]), nil
}

func getPriority(r rune) (int, error) {
	if 'a' <= r && r <= 'z' {
		return int(r - 'a' + 1), nil
	}

	if 'A' <= r && r <= 'Z' {
		return int(r - 'A' + 27), nil
	}

	return 0, fmt.Errorf("rune %v (%v) cannot be prioritized", string(r), r)
}

func PartOne(logger *zap.SugaredLogger) {
	rucksacks, err := getPartOneData()
	if err != nil {
		logger.Error(err)
		return
	}

	total := 0

	for _, items := range rucksacks {
		common, err := getCommonLetter(items[:len(items)/2], items[len(items)/2:])
		if err != nil {
			logger.Error(err)
			return
		}

		priority, err := getPriority(common)
		if err != nil {
			logger.Error(err)
			return
		}

		total += priority
	}

	logger.Info(total)
}

func PartTwo(logger *zap.SugaredLogger) {
	groups, err := getPartTwoData()
	if err != nil {
		logger.Error(err)
		return
	}

	total := 0

	for _, group := range groups {
		common, err := getCommonLetter(group...)
		if err != nil {
			logger.Error(err)
			return
		}

		priority, err := getPriority(common)
		if err != nil {
			logger.Error(err)
			return
		}

		total += priority
	}

	logger.Info(total)
}
