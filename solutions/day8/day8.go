package main

import (
	"embed"
	"fmt"

	"github.com/ShajeshJ/adventofcode_2022/common/logging"
	"github.com/ShajeshJ/adventofcode_2022/common/util"
)

var log = logging.GetLogger()

//go:embed input.txt
var files embed.FS

func getPartOneInput() (data [][]int) {
	for i, line := range util.ReadProblemInput(files) {
		data = append(data, []int{})
		for _, r := range line {
			data[i] = append(data[i], util.AtoiNoError(string(r)))
		}
	}
	return
}

func hiddenHorizontally(x int, y int, trees [][]int) bool {
	hidden := false

	for j := 0; j < y; j++ {
		if trees[x][j] >= trees[x][y] {
			hidden = true
			break
		}
	}

	// Visible from the left
	if !hidden {
		return false
	}

	for j := y + 1; j < len(trees[x]); j++ {
		if trees[x][j] >= trees[x][y] {
			return true
		}
	}

	// Visible from the right
	return false
}

func hiddenVertically(x int, y int, trees [][]int) bool {
	hidden := false

	for i := 0; i < x; i++ {
		if trees[i][y] >= trees[x][y] {
			hidden = true
			break
		}
	}

	// Visible from above
	if !hidden {
		return false
	}

	for i := x + 1; i < len(trees); i++ {
		if trees[i][y] >= trees[x][y] {
			return true
		}
	}

	// Visible from below
	return false
}

func calculateScenicScore(x int, y int, trees [][]int) int {
	score := 1
	numViewable := 0

	// From above
	for i := x - 1; i >= 0; i-- {
		numViewable++
		if trees[i][y] >= trees[x][y] {
			break
		}
	}

	score *= numViewable
	numViewable = 0

	// From below
	for i := x + 1; i < len(trees); i++ {
		numViewable++
		if trees[i][y] >= trees[x][y] {
			break
		}
	}

	score *= numViewable
	numViewable = 0

	// From the left
	for j := y - 1; j >= 0; j-- {
		numViewable++
		if trees[x][j] >= trees[x][y] {
			break
		}
	}

	score *= numViewable
	numViewable = 0

	// From the right
	for j := y + 1; j < len(trees[x]); j++ {
		numViewable++
		if trees[x][j] >= trees[x][y] {
			break
		}
	}

	score *= numViewable
	return score
}

func PartOne() any {
	data := getPartOneInput()
	total := 0

	for i, row := range data {
		if i == 0 || i == len(data)-1 {
			total += len(row)
			continue
		}

		for j := range data {
			if j == 0 || j == len(row)-1 {
				total++
				continue
			}

			if !hiddenHorizontally(i, j, data) || !hiddenVertically(i, j, data) {
				total++
			}
		}
	}

	return total
}

func PartTwo() any {
	data := getPartOneInput()
	maxScore := 0

	for i, row := range data {
		if i == 0 || i == len(data)-1 {
			continue // trees at the edge get an automatic score of 0 (0 * x = 0)
		}

		for j := range data {
			if j == 0 || j == len(row)-1 {
				continue // trees at the edge get an automatic score of 0 (0 * x = 0)
			}

			if res := calculateScenicScore(i, j, data); res > maxScore {
				maxScore = res
			}
		}
	}

	return maxScore
}

func main() {
	log.Infow(fmt.Sprintf("Answer: %v", PartOne()), "part", 1)
	log.Infow(fmt.Sprintf("Answer: %v", PartTwo()), "part", 2)
}
