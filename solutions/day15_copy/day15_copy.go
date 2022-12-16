package main

import (
	"embed"
	"fmt"
	"regexp"

	"github.com/ShajeshJ/adventofcode_2022/common/logging"
	"github.com/ShajeshJ/adventofcode_2022/common/util"
)

var log = logging.GetLogger()

//go:embed part1.txt
var files embed.FS

type State int

const (
	Unknown = iota
	HasBeacon
	HasSensor
	NoBeacon
)

type Sensor struct {
	Point  [2]int
	Beacon [2]int
}

func getPartOneData() []Sensor {
	xregex := regexp.MustCompile(`x=(-?\d+)`)
	yregex := regexp.MustCompile(`y=(-?\d+)`)

	var sensors []Sensor
	for _, line := range util.ReadProblemInput(files, 1) {
		xpoints := xregex.FindAllStringSubmatch(line, 2)
		ypoints := yregex.FindAllStringSubmatch(line, 2)
		sensors = append(sensors, Sensor{
			Point:  [2]int{util.AtoiNoError(xpoints[0][1]), util.AtoiNoError(ypoints[0][1])},
			Beacon: [2]int{util.AtoiNoError(xpoints[1][1]), util.AtoiNoError(ypoints[1][1])},
		})
	}
	return sensors
}

func ManhattanDist(p1, p2 [2]int) int {
	return util.Abs(p1[0]-p2[0]) + util.Abs(p1[1]-p2[1])
}

func PartOne() any {
	sensors := getPartOneData()
	targetY := 2_000_000
	atTargetY := map[int]State{}

	for _, s := range sensors {
		if s.Point[1] == targetY {
			atTargetY[s.Point[0]] = HasSensor
		}
		if s.Beacon[1] == targetY {
			atTargetY[s.Beacon[0]] = HasBeacon
		}

		beaconlessRange := ManhattanDist(s.Point, s.Beacon)
		beaconlessRange -= util.Abs(s.Point[1] - targetY) // sub the Y distance, and focus only on X

		for x := s.Point[0] - beaconlessRange; x <= s.Point[0]+beaconlessRange; x++ {
			if atTargetY[x] == Unknown {
				atTargetY[x] = NoBeacon
			}
		}
	}

	noBeaconCount := 0
	for _, s := range atTargetY {
		if s == NoBeacon || s == HasSensor {
			noBeaconCount++
		}
	}
	return noBeaconCount
}

type Range struct {
	Min, Max int
}

type SensorMap struct {
	MaxLen int
	// Covered ranges in order
	// new ranges added that overlap with existing ones will be merged
	Covered []Range
}

func NewSensorMap(squareLen int) SensorMap {
	return SensorMap{MaxLen: squareLen, Covered: []Range{}}
}

func (sm *SensorMap) TrySet(x, miny, maxy int) {
	if x < 0 || x >= sm.MaxLen {
		return
	}

	if miny < 0 {
		miny = 0
	}
	if maxy >= sm.MaxLen {
		maxy = sm.MaxLen - 1
	}

	if miny > maxy {
		return
	}

	// candidate range to insert
	cr := Range{x*sm.MaxLen + miny, x*sm.MaxLen + maxy}

	// reconstruct list of ranges, merging `cr` with existing ones where appropriate
	newCovered := []Range{}
	breakIndex := -1
	for i, r := range sm.Covered {

		if r.Max+1 < cr.Min {
			// looped range is at least 1 space before the candidate
			newCovered = append(newCovered, r)
			continue
		}

		if cr.Max+1 < r.Min {
			// looped range is at least 1 space after the candidate;
			// `cr` is complete, and we can add unappended remainder
			breakIndex = i
			break
		}

		// ranges need to be merged
		cr.Min = util.Min(r.Min, cr.Min)
		cr.Max = util.Max(r.Max, cr.Max)
	}

	newCovered = append(newCovered, cr)
	if breakIndex != -1 {
		newCovered = append(newCovered, sm.Covered[breakIndex:]...)
	}

	sm.Covered = newCovered
}

func PartTwo() any {
	sensors := getPartOneData()
	sensorMap := NewSensorMap(4_000_001) // {max coordinate}+1

	var distressX int
	for t, s := range sensors {
		distressX = t
		beaconlessRange := ManhattanDist(s.Point, s.Beacon)
		for xOffset := 0 - beaconlessRange; xOffset <= beaconlessRange; xOffset++ {
			sensorMap.TrySet(
				s.Point[0]+xOffset,
				s.Point[1]+util.Abs(xOffset)-beaconlessRange,
				s.Point[1]+beaconlessRange-util.Abs(xOffset),
			)
		}
	}

	distressX, distressY := -1, -1
	if len(sensorMap.Covered) > 2 {
		panic("too many ranges")
	}
	coord := sensorMap.Covered[0].Max + 1
	distressX, distressY = coord/sensorMap.MaxLen, coord%sensorMap.MaxLen
	return distressX*4_000_000 + distressY
}

func main() {
	log.Infow(fmt.Sprintf("Answer: %v", PartOne()), "part", 1)
	log.Infow(fmt.Sprintf("Answer: %v", PartTwo()), "part", 2)
}
