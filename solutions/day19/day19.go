package main

import (
	"embed"
	"fmt"
	"regexp"
	"sync"

	"github.com/ShajeshJ/adventofcode_2022/common/logging"
	"github.com/ShajeshJ/adventofcode_2022/common/util"
	"golang.org/x/exp/slices"
)

var log = logging.GetLogger()

//go:embed part1.txt
var files embed.FS

const (
	None = iota
	Ore
	Clay
	Obsidian
	Geode
)

type ResourceType int

type Items struct {
	Ore      int
	Clay     int
	Obsidian int
	Geode    int
}

func (i *Items) Add(t ResourceType, amt int) {
	switch t {
	case Ore:
		i.Ore += amt
	case Clay:
		i.Clay += amt
	case Obsidian:
		i.Obsidian += amt
	case Geode:
		i.Geode += amt
	}
}

type Robot = Items

func (i *Items) Contains(other Items) bool {
	return i.Ore >= other.Ore &&
		i.Clay >= other.Clay &&
		i.Obsidian >= other.Obsidian
	// don't bother checking geode
}

// Add adds the items in other to i
func Add(i, other Items) Items {
	return Items{
		Ore:      i.Ore + other.Ore,
		Clay:     i.Clay + other.Clay,
		Obsidian: i.Obsidian + other.Obsidian,
		Geode:    i.Geode + other.Geode,
	}
}

// Subtract subtracts the items in other from i
func Subtract(i, other Items) Items {
	return Items{
		Ore:      i.Ore - other.Ore,
		Clay:     i.Clay - other.Clay,
		Obsidian: i.Obsidian - other.Obsidian,
		Geode:    i.Geode - other.Geode,
	}
}

type Blueprint struct {
	ID    int
	Costs map[ResourceType]Items
}

// SumN is a function which finds the sum from 1 to n
func SumN(n int) int {
	return (n * (n + 1)) / 2
}

func getPartOneData() []Blueprint {
	// Regex Indicies:
	// Blueprint {0}: Each ore robot costs {1} ore.
	// Each clay robot costs {2} ore.
	// Each obsidian robot costs {3} ore and {4} clay.
	// Each geode robot costs {5} ore and {6} obsidian.
	regex := regexp.MustCompile(`\d+`)
	blueprints := []Blueprint{}
	for _, line := range util.ReadProblemInput(files, 1) {
		vals := regex.FindAllStringSubmatch(line, -1)
		bp := Blueprint{util.AtoiNoError(vals[0][0]), make(map[ResourceType]Items)}
		bp.Costs[Ore] = Items{Ore: util.AtoiNoError(vals[1][0])}
		bp.Costs[Clay] = Items{Ore: util.AtoiNoError(vals[2][0])}
		bp.Costs[Obsidian] = Items{Ore: util.AtoiNoError(vals[3][0]), Clay: util.AtoiNoError(vals[4][0])}
		bp.Costs[Geode] = Items{Ore: util.AtoiNoError(vals[5][0]), Obsidian: util.AtoiNoError(vals[6][0])}
		blueprints = append(blueprints, bp)
	}
	return blueprints
}

type DecisionFactors struct {
	bp          Blueprint
	minute      int
	items       Items
	bots        Robot
	nextBotType ResourceType
	curMax      Items
	exclude     []ResourceType
}

func Simulate(d DecisionFactors) Items {
	d.exclude = []ResourceType{}
	if d.nextBotType == None {
		// If we had the ability to build a bot in this minute,
		// but chose to do nothing, then it doesn't make sense to
		// build those bots during the next minute; so add those to the exclude list
		for bot, cost := range d.bp.Costs {
			if d.items.Contains(cost) {
				d.exclude = append(d.exclude, bot)
			}
		}
	} else {
		d.items = Subtract(d.items, d.bp.Costs[d.nextBotType])
	}

	d.items = Add(d.items, d.bots)
	if d.nextBotType != None {
		d.bots.Add(d.nextBotType, 1)
	}
	d.minute--

	return Decide(d)
}

func Decide(d DecisionFactors) Items {
	if d.minute == 0 {
		return d.items
	}
	if d.items.Geode+d.bots.Geode*d.minute+SumN(d.minute-1) <= d.curMax.Geode {
		// curMax is unbeatable even if we make geode bots for the remaining time; prune this branch
		return d.items
	}

	// Intentionally reverse order to prioritize most useful bots
	for _, bot := range []ResourceType{Geode, Obsidian, Clay, Ore, None} {
		if slices.Contains(d.exclude, bot) || !d.items.Contains(d.bp.Costs[bot]) {
			continue
		}
		d.nextBotType = bot
		test := Simulate(d)
		if test.Geode > d.curMax.Geode {
			d.curMax = test
		}
	}

	return d.curMax
}

func GetMaxGeodes(bp Blueprint, minute int) int {
	return Decide(DecisionFactors{
		bp:          bp,
		minute:      minute,
		items:       Items{Ore: 0, Clay: 0, Obsidian: 0, Geode: 0},
		bots:        Robot{Ore: 1, Clay: 0, Obsidian: 0, Geode: 0},
		nextBotType: None,
		curMax:      Items{Ore: 0, Clay: 0, Obsidian: 0, Geode: 0},
		exclude:     []ResourceType{},
	}).Geode
}

func PartOne() any {
	total := 0
	var wg sync.WaitGroup

	for _, bp := range getPartOneData() {
		bp := bp
		wg.Add(1)
		go func() {
			defer wg.Done()
			total += GetMaxGeodes(bp, 24) * bp.ID
		}()
	}
	wg.Wait()
	return total
}

func PartTwo() any {
	const max = 3
	numGeodes := make([]int, max)
	var wg sync.WaitGroup

	for _, bp := range getPartOneData()[:max] {
		bp := bp
		wg.Add(1)
		go func() {
			defer wg.Done()
			numGeodes[bp.ID-1] = GetMaxGeodes(bp, 32)
		}()
	}
	wg.Wait()
	return numGeodes[0] * numGeodes[1] * numGeodes[2]
}

func main() {
	log.Infow(fmt.Sprintf("Answer: %v", PartOne()), "part", 1)
	log.Infow(fmt.Sprintf("Answer: %v", PartTwo()), "part", 2)
}
