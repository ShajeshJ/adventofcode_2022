package main

import (
	"embed"
	"fmt"
	"regexp"
	"strings"

	"github.com/ShajeshJ/adventofcode_2022/common/logging"
	"github.com/ShajeshJ/adventofcode_2022/common/util"
)

var log = logging.GetLogger()

//go:embed part1.txt
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
	return c.Features[y][x-c.Xoffset+2]
}

func (c *CaveMap) Set(x int, y int, val MapFeature) {
	// Increase X to the left to always keep 1 column as buffer
	for x-c.Xoffset+2 < 1 && val != Air {
		for eachY := 0; eachY < len(c.Features); eachY++ {
			c.Features[eachY] = append([]MapFeature{Air}, c.Features[eachY]...)
		}
		c.Xoffset--
	}

	// Increase X to the right to always keep 1 column as buffer
	for x-c.Xoffset+2 >= len(c.Features[y])-1 && val != Air {
		for eachY := 0; eachY < len(c.Features); eachY++ {
			c.Features[eachY] = append(c.Features[eachY], Air)
		}
	}
	c.Features[y][x-c.Xoffset+2] = val
}

func NewCaveMap(minX, maxX, maxY int) CaveMap {
	cavemap := CaveMap{
		Xoffset:  minX + 1,
		Features: make([][]MapFeature, maxY+1), // 0 <= depth <= maxY
	}
	for i := 0; i < len(cavemap.Features); i++ {
		// maxX <= len <= maxX && Add 1 column on either side
		cavemap.Features[i] = make([]MapFeature, maxX-minX+1+2)
	}
	return cavemap
}

func getPartOneData() CaveMap {
	// Get size our map
	minX, maxX := 1000, 0
	maxY := 0
	pointRegex := regexp.MustCompile(`(\d+),(\d+)`)

	data := util.ReadProblemInput(files, 1)

	for _, line := range data {
		for _, point := range pointRegex.FindAllStringSubmatch(line, -1) {
			x := util.AtoiNoError(point[1])
			if x < minX {
				minX = x
			}
			if x > maxX {
				maxX = x
			}

			y := util.AtoiNoError(point[2])
			if y > maxY {
				maxY = y
			}
		}
	}

	cavemap := NewCaveMap(minX, maxX, maxY)

	for _, line := range data {
		formationPoints := strings.Split(line, " -> ")
		prev := util.Map(
			strings.Split(formationPoints[0], ","),
			func(x string) int { return util.AtoiNoError(x) },
		)

		for _, p := range formationPoints[1:] {
			next := util.Map(
				strings.Split(p, ","),
				func(x string) int { return util.AtoiNoError(x) },
			)

			for x := util.Min(prev[0], next[0]); x <= util.Max(prev[0], next[0]); x++ {
				for y := util.Min(prev[1], next[1]); y <= util.Max(prev[1], next[1]); y++ {
					cavemap.Set(x, y, Rock)
				}
			}

			prev = next
		}
	}

	cavemap.Set(500, 0, SandSource)

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
	sandX, sandY := 500, 0

	for true {
		sandX, sandY = 500, 0
		if cavemap.Get(sandX, sandY) != SandSource {
			// PrintCaveMap(cavemap)
			return numsand // Source is blocked
		}

		for true {
			sandY++
			if sandY >= cavemap.Depth() {
				if !withfloor {
					// PrintCaveMap(cavemap)
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
	return numsand
}

func PartTwo() any {
	cavemap := getPartOneData()
	// Add 1 extra space at the bottom for the floor
	cavemap.Features = append(cavemap.Features, make([]MapFeature, len(cavemap.Features[0])))
	numsand := SimulateSand(cavemap, true)
	return numsand
}

func main() {
	log.Infow(fmt.Sprintf("Answer: %v", PartOne()), "part", 1)
	log.Infow(fmt.Sprintf("Answer: %v", PartTwo()), "part", 2)
}
