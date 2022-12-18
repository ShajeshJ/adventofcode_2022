package main

import (
	"embed"
	"fmt"

	"github.com/ShajeshJ/adventofcode_2022/common/logging"
	"github.com/ShajeshJ/adventofcode_2022/common/util"
	"golang.org/x/exp/slices"
)

var log = logging.GetLogger()

//go:embed part1.txt
var files embed.FS

type Direction int

const (
	Down = iota
	Left
	Right
)

const ChamberWidth = 7
const ChamberVerticalBuffer = 7

const (
	RockRune       = '#'
	EmptySpaceRune = '.'
)

// GetWindGenerator returns a function that will return the wind direction
// at the given index, and the index for the next wind direction
func GetWindGenerator(dirs string) func(i int) (Direction, int) {
	wind := []rune(dirs)
	return func(i int) (Direction, int) {
		if wind[i%len(wind)] == '<' {
			return Left, (i + 1) % len(wind)
		} else {
			return Right, (i + 1) % len(wind)
		}
	}
}

// CreateChamber returns a 2D rune array representing the falling rock chamber
func CreateChamber() *[][]rune {
	// Start 7 tall to accomodate tallest vertical rock
	chamber := make([][]rune, ChamberVerticalBuffer)

	for y := 0; y < len(chamber); y++ {
		chamber[y] = make([]rune, ChamberWidth)
		for x := 0; x < len(chamber[y]); x++ {
			chamber[y][x] = EmptySpaceRune
		}
	}

	return &chamber
}

// ExpandChamber expands the vertical height of the chamber up to the minimum necessary buffer
func ExpandChamber(height int, chamber *[][]rune) {
	// appendAmt = minBuffer - numEmptyRowsFromTop
	// numEmptyRowsFromTop = chamberHeight - highestRockHeight
	appendAmt := ChamberVerticalBuffer - (len(*chamber) - height)

	for i := 0; i < appendAmt; i++ {
		row := make([]rune, ChamberWidth)
		for x := 0; x < len(row); x++ {
			row[x] = EmptySpaceRune
		}
		*chamber = append(*chamber, row)
	}
}

// GetHeight returns the height of the highest rock in `chamber`
func GetHeight(chamber [][]rune) int {
	for i := len(chamber) - 1; i >= 0; i-- {
		if slices.Contains(chamber[i], RockRune) {
			return i + 1
		}
	}
	return 0
}

func PrintChamber(chamber [][]rune) {
	for i := len(chamber) - 1; i >= 0; i-- {
		log.Info("|", string(chamber[i]), "|", i)
	}
	log.Info("+-------+")
}

type RockWindIndexes struct {
	Rock, Wind int
}

// RunRockSimulation will run the run simulation from the `start` indexes,
// and will attempt to drop `numRocks` until the `end` indexes. If the `end` indexes
// are -1, then the simulation will go until all `numRocks` are thrown.
// The rock formation height and number of rocks thrown will be returned
func RunRockSimulation(start, end RockWindIndexes, numRocks int, chamber *[][]rune) (int, int) {
	input := util.ReadProblemInput(files, 1)[0]
	getWind := GetWindGenerator(input)
	getRock := GetRockGenerator()

	rockLoopIdx, windLoopIdx := start.Rock, start.Wind
	startHeight := GetHeight(*chamber)

	for i := 1; i <= numRocks; i++ {
		height := GetHeight(*chamber)
		ExpandChamber(height, chamber)

		var rock Rock
		rock, rockLoopIdx = getRock(rockLoopIdx)
		rock.InitPosition(2, height+3)

		for {
			var wind Direction
			wind, windLoopIdx = getWind(windLoopIdx)
			rock.Move(wind, *chamber)
			stillFalling := rock.Move(Down, *chamber)
			if !stillFalling {
				rock.PlaceRock(chamber)
				break
			}
		}

		if end.Rock == -1 || end.Wind == -1 {
			continue
		}

		if end.Rock == rockLoopIdx && end.Wind == windLoopIdx {
			// PrintChamber(chamber)
			return GetHeight(*chamber) - startHeight, numRocks - (numRocks - i)
		}
	}

	// PrintChamber(chamber)
	return GetHeight(*chamber) - startHeight, numRocks
}

func PartOne() any {
	height, _ := RunRockSimulation(
		RockWindIndexes{0, 0},
		RockWindIndexes{-1, -1},
		2022,
		CreateChamber(),
	)
	return height
}

// GetRepeatingIndexes simulates the rock fall from part 1, until
// it finds a repeating index pattern
func GetRepeatingIndexes() (rockIdx, windIdx int) {
	input := util.ReadProblemInput(files, 1)[0]
	getWind := GetWindGenerator(input)
	getRock := GetRockGenerator()
	chamber := *CreateChamber()

	foundPairs := []RockWindIndexes{}

	for {
		height := GetHeight(chamber)
		ExpandChamber(height, &chamber)

		var rock Rock
		rock, rockIdx = getRock(rockIdx)
		rock.InitPosition(2, height+3)

		for {
			var wind Direction
			wind, windIdx = getWind(windIdx)
			rock.Move(wind, chamber)
			stillFalling := rock.Move(Down, chamber)
			if !stillFalling {
				rock.PlaceRock(&chamber)
				break
			}
		}

		curPair := RockWindIndexes{rockIdx, windIdx}
		if slices.Contains(foundPairs, curPair) {
			return
		}
		foundPairs = append(foundPairs, curPair)
	}
}

func PartTwo() any {
	totalHeight, totalRocks := 0, 1_000_000_000_000
	loopRockIdx, loopWindIdx := GetRepeatingIndexes()
	chamber := CreateChamber()

	// Simulate the first bit, just before the looped portion
	startSimHeight, startSimThrown := RunRockSimulation(
		RockWindIndexes{0, 0},
		RockWindIndexes{loopRockIdx, loopWindIdx},
		totalRocks,
		chamber,
	)
	totalHeight += startSimHeight
	totalRocks -= startSimThrown

	// Simulate the first loop explicitly so the final chamber terrain matches the end of loop
	l1SimHeight, l1SimThrown := RunRockSimulation(
		RockWindIndexes{loopRockIdx, loopWindIdx},
		RockWindIndexes{loopRockIdx, loopWindIdx},
		totalRocks,
		chamber,
	)
	totalHeight += l1SimHeight
	totalRocks -= l1SimThrown

	// Simulate loop again, but this time we can simply multiple the resulting outputs to
	// quickly multiple/add the looped portion without literally simulating it
	loopSimHeight, loopSimThrown := RunRockSimulation(
		RockWindIndexes{loopRockIdx, loopWindIdx},
		RockWindIndexes{loopRockIdx, loopWindIdx},
		totalRocks,
		chamber,
	)
	requiredIterations := totalRocks / loopSimThrown
	totalHeight += loopSimHeight * requiredIterations
	totalRocks -= loopSimThrown * requiredIterations

	// Simulate the remaining rocks, after the final looped portion
	endSimHeight, endSimThrown := RunRockSimulation(
		RockWindIndexes{loopRockIdx, loopWindIdx},
		RockWindIndexes{-1, -1},
		totalRocks,
		chamber,
	)

	totalHeight += endSimHeight
	totalRocks -= endSimThrown

	if totalRocks != 0 {
		log.Info(totalRocks)
		panic("something went wrong")
	}

	return totalHeight
}

func main() {
	log.Infow(fmt.Sprintf("Answer: %v", PartOne()), "part", 1)
	log.Infow(fmt.Sprintf("Answer: %v", PartTwo()), "part", 2)
}
