package main

import (
	"embed"
	"fmt"
	"regexp"

	"github.com/ShajeshJ/adventofcode_2022/common/logging"
	"github.com/ShajeshJ/adventofcode_2022/common/util"
	"golang.org/x/exp/slices"
)

var log = logging.GetLogger()

//go:embed input.txt
var files embed.FS

type Sensor struct {
	Point         []int
	Beacon        []int
	NoBeaconRange int
}

func ManhattanDist(p1, p2 []int) int {
	return util.Abs(p1[0]-p2[0]) + util.Abs(p1[1]-p2[1])
}

func getPartOneData() []Sensor {
	xregex := regexp.MustCompile(`x=(-?\d+)`)
	yregex := regexp.MustCompile(`y=(-?\d+)`)

	var sensors []Sensor
	for _, line := range util.ReadProblemInput(files) {
		xpoints := xregex.FindAllStringSubmatch(line, 2)
		ypoints := yregex.FindAllStringSubmatch(line, 2)
		s := Sensor{
			Point:  []int{util.AtoiNoError(xpoints[0][1]), util.AtoiNoError(ypoints[0][1])},
			Beacon: []int{util.AtoiNoError(xpoints[1][1]), util.AtoiNoError(ypoints[1][1])},
		}
		s.NoBeaconRange = ManhattanDist(s.Point, s.Beacon)
		sensors = append(sensors, s)
	}
	return sensors
}

type Range struct {
	Min int
	Max int
}

func CombineOrderedRanges(r Range, ranges []Range) []Range {
	combined := []Range{}

	breakIndex := -1
	for i, r2 := range ranges {
		if r2.Max+1 < r.Min {
			combined = append(combined, r2)
			continue
		}
		if r.Max+1 < r2.Min {
			breakIndex = i
			break
		}
		r.Min = util.Min(r.Min, r2.Min)
		r.Max = util.Max(r.Max, r2.Max)
	}
	combined = append(combined, r)
	if breakIndex != -1 {
		combined = append(combined, ranges[breakIndex:]...)
	}
	return combined
}

func PartOne() any {
	sensors := getPartOneData()
	targetY := 2_000_000
	atTargetY := []Range{}
	beaconsAtTargetY := []int{}

	noBeaconCount := 0
	for _, s := range sensors {
		if s.Beacon[1] == targetY && !slices.Contains(beaconsAtTargetY, s.Beacon[0]) {
			noBeaconCount--
			beaconsAtTargetY = append(beaconsAtTargetY, s.Beacon[0])
		}

		// sub the Y distance, and focus only on X
		beaconlessRange := s.NoBeaconRange - util.Abs(s.Point[1]-targetY)
		if beaconlessRange < 0 {
			continue
		}

		atTargetY = CombineOrderedRanges(Range{s.Point[0] - beaconlessRange, s.Point[0] + beaconlessRange}, atTargetY)
	}

	for _, r := range atTargetY {
		noBeaconCount += r.Max - r.Min + 1
	}
	return noBeaconCount
}

// InSensorRange returns true if the point is in range of a sensor, and returns the next
// y-coordinate that's outside of that sensor's range; otherwise returns false and `y`
func InSensorRange(x, y int, sensors []Sensor) (bool, int) {
	for _, s := range sensors {
		if ManhattanDist(s.Point, []int{x, y}) <= s.NoBeaconRange {
			// nextY = {no-beacon range around the sensor} - (x dist from sensor) + 1
			return true, s.Point[1] + (s.NoBeaconRange - util.Abs(s.Point[0]-x)) + 1
		}
	}
	return false, y
}

func PartTwo() any {
	sensors := getPartOneData()
	maxCoords := 4_000_000

	distressX, distressY := -1, -1
	for x := 0; x <= maxCoords; x++ {
		y := 0
		for y <= maxCoords {
			var inRange bool
			if inRange, y = InSensorRange(x, y, sensors); !inRange {
				distressX, distressY = x, y
				break
			}
		}

		if distressX != -1 {
			break
		}
	}

	return distressX*4_000_000 + distressY
}

func main() {
	log.Infow(fmt.Sprintf("Answer: %v", PartOne()), "part", 1)
	log.Infow(fmt.Sprintf("Answer: %v", PartTwo()), "part", 2)
}
