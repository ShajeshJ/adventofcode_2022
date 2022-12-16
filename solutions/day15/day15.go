package main

import (
	"embed"
	"fmt"
	"regexp"
	"sync"

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
	for _, line := range util.ReadProblemInput(files, 1) {
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

		// sub the Y distance, and focus only on X
		beaconlessRange := s.NoBeaconRange - util.Abs(s.Point[1]-targetY)

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

func CanHaveDistressBeacon(x, y int, sensors []Sensor) bool {
	for _, s := range sensors {
		if ManhattanDist(s.Point, []int{x, y}) <= s.NoBeaconRange {
			return false
		}
	}
	return true
}

// answer was 13673971349056 at x=3418400-3418800
func PartTwo() any {
	sensors := getPartOneData()
	maxCoords := 4_000_000

	distressX, distressY := -1, -1
	var wg sync.WaitGroup
	for x := 0; x <= maxCoords; x++ {
		x := x
		wg.Add(1)
		go func() {
			defer wg.Done()
			for y := 0; y <= maxCoords; y++ {
				if CanHaveDistressBeacon(x, y, sensors) {
					distressX, distressY = x, y
					log.Infof("Early exit value: %v", distressX*4_000_000+distressY)
					break
				}
			}
		}()

		if x%400 == 0 {
			log.Infof("waiting for batch x=%v", x)
			wg.Wait()
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
