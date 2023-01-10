package main

import (
	"embed"
	"fmt"
	"math"
	"strings"

	"github.com/ShajeshJ/adventofcode_2022/common/logging"
	"github.com/ShajeshJ/adventofcode_2022/common/util"
	"golang.org/x/exp/slices"
)

var log = logging.GetLogger()

//go:embed input.txt
var files embed.FS

type Valve struct {
	ID      string
	Rate    float64
	LeadsTo []string
}

func getPartOneData() map[string]Valve {
	lines := util.ReadProblemInput(files)
	valves := map[string]Valve{}

	// [0]  [1] [2]  [3]   [4]     [5]   [6] [7]  [8]  [9]
	// Valve XX has flow rate=X; tunnels lead to valves XX, XX, ...
	for _, line := range lines {
		tokens := strings.Split(line, " ")
		id := tokens[1]
		rate := float64(util.AtoiNoError(strings.Split(strings.Trim(tokens[4], ";"), "=")[1]))
		leadTo := []string{}
		for _, t := range tokens[9:] {
			leadTo = append(leadTo, strings.Trim(t, ","))
		}
		valves[id] = Valve{id, rate, leadTo}
	}

	return valves
}

type GraphConnectivity struct {
	Graph  map[string]map[string]float64
	Valves map[string]Valve
}

func InitShortestDistGraph(valves map[string]Valve) GraphConnectivity {
	g := GraphConnectivity{map[string]map[string]float64{}, valves}

	for id, valve := range valves {
		g.Graph[id] = map[string]float64{id: 0}
		for _, nextValveID := range valve.LeadsTo {
			g.Graph[id][nextValveID] = 1
		}
	}

	for dstID := range g.Valves {
		for srcID := range g.Valves {
			if _, ok := g.Graph[dstID][srcID]; !ok {
				g.Graph[dstID][srcID] = math.Inf(1)
			}
		}
	}

	return g
}

func GetAllShortestDist(valves map[string]Valve) GraphConnectivity {
	g := InitShortestDistGraph(valves)

	for inter := range g.Valves {
		for src := range g.Graph {
			for dest := range g.Graph[src] {
				if g.Graph[src][inter]+g.Graph[inter][dest] < g.Graph[src][dest] {
					g.Graph[src][dest] = g.Graph[src][inter] + g.Graph[inter][dest]
				}
			}
		}
	}

	return g
}

func GetReleasedAmount(duration float64, opened []string, g GraphConnectivity) float64 {
	var released float64 = 0
	for _, id := range opened {
		released += g.Valves[id].Rate * duration
	}
	return released
}

func GetHighestPressure(src, dest string, remaining float64, opened []string, g GraphConnectivity) float64 {
	var released float64 = 0

	if remaining <= g.Graph[src][dest] {
		// Spend remaining time just running
		return GetReleasedAmount(remaining, opened, g)
	}

	// Run to dest valve, and open it
	released += GetReleasedAmount(g.Graph[src][dest]+1, opened, g)
	remaining -= g.Graph[src][dest] + 1

	if remaining == 0 {
		// Spent remaining time just opening the valve
		return released
	}

	// Add our destination as an opened valve
	nextOpen := make([]string, len(opened))
	copy(nextOpen, opened)
	nextOpen = append(nextOpen, dest)

	// Try remaining unopened valves, and use the highest released amount
	var nextRelease float64 = 0
	for nextDest := range g.Valves {
		if !slices.Contains(nextOpen, nextDest) {
			destRelease := GetHighestPressure(dest, nextDest, remaining, nextOpen, g)
			if destRelease > nextRelease {
				nextRelease = destRelease
			}
		}
	}

	// No more valves to open; just sit and wait for more pressure
	if nextRelease == 0 {
		nextRelease = GetReleasedAmount(remaining, nextOpen, g)
	}

	return released + nextRelease
}

func PartOne() any {
	valves := getPartOneData()
	g := GetAllShortestDist(valves)

	// We don't care to travel to valves with 0 release rate
	for k, v := range g.Valves {
		if v.Rate == 0 {
			delete(g.Valves, k)
		}
	}

	var highestRelease float64 = 0
	for id := range g.Valves {
		test := GetHighestPressure("AA", id, 30, []string{}, g)
		if test > highestRelease {
			highestRelease = test
		}
	}

	return highestRelease
}

type Runner struct {
	Src  string
	Dest string
	Dur  float64
}

func GetRemainingValveIDs(opened []string, g GraphConnectivity) []string {
	remaining := []string{}
	for id := range g.Valves {
		if !slices.Contains(opened, id) {
			remaining = append(remaining, id)
		}
	}
	return remaining
}

func SimPressureWithHelp(h, e Runner, remaining float64, opened []string, g GraphConnectivity) float64 {
	if remaining == 0 || len(opened) == len(g.Valves) {
		// If no time remaining, return 0
		// If everything's opened, just wait for remaining time
		return GetReleasedAmount(remaining, opened, g)
	}

	released := GetReleasedAmount(1, opened, g)
	remaining--
	h.Dur--
	e.Dur--

	if (h.Dur > 0 || h.Dest == "") && (e.Dur > 0 || e.Dest == "") {
		// At least one runner still going, and nobody opened a new valve
		return released + SimPressureWithHelp(h, e, remaining, opened, g)
	}

	// At least 1 valve opening

	nextOpen := make([]string, len(opened))
	copy(nextOpen, opened)

	// Set new runners for ones that finished
	hdone, edone := h.Dur <= 0, e.Dur <= 0
	if hdone {
		nextOpen = append(nextOpen, h.Dest)
		h = Runner{Src: h.Dest}
	}
	if edone {
		nextOpen = append(nextOpen, e.Dest)
		e = Runner{Src: e.Dest}
	}

	var bestRelease float64 = 0
	simNextRelease := func() { // Helper for toil
		test := SimPressureWithHelp(h, e, remaining, nextOpen, g)
		if test > bestRelease {
			bestRelease = test
		}
	}
	resetDur := func(r *Runner) {
		r.Dur = g.Graph[r.Src][r.Dest] + 1
	}

	var remainingIDs []string
	if hdone && !edone {
		remainingIDs = GetRemainingValveIDs(append(nextOpen, e.Dest), g)
	} else if edone && !hdone {
		remainingIDs = GetRemainingValveIDs(append(nextOpen, h.Dest), g)
	} else {
		remainingIDs = GetRemainingValveIDs(nextOpen, g)
	}

	if len(remainingIDs) == 0 {
		// Runners don't matter; just sim arbitrarily
		simNextRelease()
		return released + bestRelease
	}

	// Exactly one runner choosing a new dest
	if hdone != edone {
		for _, id := range remainingIDs {
			if hdone {
				h.Dest = id
				resetDur(&h)
			} else {
				e.Dest = id
				resetDur(&e)
			}
			simNextRelease()
		}
		return released + bestRelease
	}

	// Both can choose a new dest

	if len(remainingIDs) == 1 {
		// If only 1 left unopened, don't bother simulating elephant
		h.Dest = remainingIDs[0]
		resetDur(&h)
		simNextRelease()
		return released + bestRelease
	}

	for i := 0; i < len(remainingIDs); i++ {
		h.Dest = remainingIDs[i]
		resetDur(&h)

		for j := i + 1; j < len(remainingIDs); j++ {
			e.Dest = remainingIDs[j]
			resetDur(&e)
			simNextRelease()
		}
	}

	return released + bestRelease
}

func PartTwo() any {
	valves := getPartOneData()
	g := GetAllShortestDist(valves)

	// We don't care to travel to valves with 0 release rate
	idlist := []string{}
	for k, v := range g.Valves {
		if v.Rate == 0 {
			delete(g.Valves, k)
		} else {
			idlist = append(idlist, k)
		}
	}

	var highestRelease float64 = 0
	h, e := Runner{Src: "AA"}, Runner{Src: "AA"}

	for i := 0; i < len(idlist); i++ {
		h.Dest = idlist[i]
		h.Dur = g.Graph[h.Src][h.Dest] + 1

		for j := i + 1; j < len(idlist); j++ {
			e.Dest = idlist[j]
			e.Dur = g.Graph[e.Src][e.Dest] + 1

			test := SimPressureWithHelp(h, e, 26, []string{}, g)
			if test > highestRelease {
				highestRelease = test
			}
		}
	}

	return highestRelease
}

func main() {
	log.Infow(fmt.Sprintf("Answer: %v", PartOne()), "part", 1)
	log.Infow(fmt.Sprintf("Answer: %v", PartTwo()), "part", 2)
}
