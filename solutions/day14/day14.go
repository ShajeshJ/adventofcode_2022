package main

import (
	"embed"
	"fmt"
	"regexp"

	"github.com/ShajeshJ/adventofcode_2022/common/logging"
	"github.com/ShajeshJ/adventofcode_2022/common/util"
)

var log = logging.GetLogger()

//go:embed input.txt
var files embed.FS

type MapFeature int

const (
	Air = iota
	Rock
	Sand
	SandSource
)

type CaveMap struct {
	Features [][]MapFeature
	Xoffset  int
}

func (c *CaveMap) Depth() int {
	return len(c.Features)
}

func (c *CaveMap) Get(x int, y int) MapFeature {
	return c.Features[y][x-c.Xoffset]
}

// Set will set the given `val` at the coordinates
// Note: This operation will scale the 2D slices as needed to fit new coordinates
func (c *CaveMap) Set(x int, y int, val MapFeature) {
	// Increase X to the left to always keep 1 column as buffer
	for x-c.Xoffset < 1 {
		for eachY := 0; eachY < len(c.Features); eachY++ {
			c.Features[eachY] = append([]MapFeature{Air}, c.Features[eachY]...)
		}
		c.Xoffset--
	}

	// Increase X to the right to always keep 1 column as buffer
	for x-c.Xoffset+1 >= len(c.Features[0]) {
		for eachY := 0; eachY < len(c.Features); eachY++ {
			c.Features[eachY] = append(c.Features[eachY], Air)
		}
	}

	// Increase Y Downwards as needed
	for y >= len(c.Features) {
		c.Features = append(c.Features, make([]MapFeature, len(c.Features[0])))
	}

	c.Features[y][x-c.Xoffset] = val
}

func NewCaveMap(x, y int, val MapFeature) CaveMap {
	cavemap := CaveMap{
		Xoffset:  x - 1,
		Features: make([][]MapFeature, y+1),
	}
	for i := 0; i < len(cavemap.Features); i++ {
		cavemap.Features[i] = make([]MapFeature, 3)
	}
	cavemap.Set(x, y, val)
	return cavemap
}

func getPartOneData() CaveMap {
	cavemap := NewCaveMap(500, 0, SandSource)
	pointRegex := regexp.MustCompile(`(\d+),(\d+)`)

	for _, line := range util.ReadProblemInput(files) {
		formationPoints := pointRegex.FindAllStringSubmatch(line, -1)
		prev := []int{
			util.AtoiNoError(formationPoints[0][1]),
			util.AtoiNoError(formationPoints[0][2]),
		}

		for _, nextMatches := range formationPoints[1:] {
			next := []int{
				util.AtoiNoError(nextMatches[1]),
				util.AtoiNoError(nextMatches[2]),
			}

			for x := util.Min(prev[0], next[0]); x <= util.Max(prev[0], next[0]); x++ {
				for y := util.Min(prev[1], next[1]); y <= util.Max(prev[1], next[1]); y++ {
					cavemap.Set(x, y, Rock)
				}
			}

			prev = next
		}
	}

	return cavemap
}

func PrintCaveMap(cavemap CaveMap) {
	visual := ""
	for _, col := range cavemap.Features {
		visual += "\n"
		for _, feature := range col {
			if feature == Air {
				visual += "."
			} else if feature == Rock {
				visual += "#"
			} else if feature == Sand {
				visual += "o"
			} else {
				visual += "+"
			}
		}
	}
	log.Info(visual)
}

func SimulateSand(cavemap CaveMap, withfloor bool) int {
	numsand := 0

	for true {
		sandX, sandY := 500, 0
		if cavemap.Get(sandX, sandY) != SandSource {
			return numsand // Source is blocked
		}

		for true {
			sandY++
			if sandY >= cavemap.Depth() {
				if !withfloor {
					return numsand // Into the abyss
				} else {
					cavemap.Set(sandX, sandY-1, Sand)
					numsand++
					break // Rests at rock bottom
				}
			}

			if cavemap.Get(sandX, sandY) == Air {
				continue // Straight down
			}

			if cavemap.Get(sandX-1, sandY) == Air {
				sandX--
				continue // Down left
			}

			if cavemap.Get(sandX+1, sandY) == Air {
				sandX++
				continue // Down right
			}

			cavemap.Set(sandX, sandY-1, Sand)
			numsand++
			break // At rest
		}
	}

	panic("should not have reached here!")
}

func PartOne() any {
	cavemap := getPartOneData()
	numsand := SimulateSand(cavemap, false)
	// PrintCaveMap(cavemap)
	return numsand
}

func PartTwo() any {
	cavemap := getPartOneData()
	cavemap.Set(500, cavemap.Depth(), Air) // Add 1 layer of air before the bottom
	numsand := SimulateSand(cavemap, true)
	// PrintCaveMap(cavemap)
	return numsand
}

func main() {
	log.Infow(fmt.Sprintf("Answer: %v", PartOne()), "part", 1)
	log.Infow(fmt.Sprintf("Answer: %v", PartTwo()), "part", 2)
}
