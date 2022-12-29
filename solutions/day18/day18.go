package main

import (
	"embed"
	"fmt"
	"strings"

	"github.com/ShajeshJ/adventofcode_2022/common/logging"
	"github.com/ShajeshJ/adventofcode_2022/common/util"
	"golang.org/x/exp/maps"
)

var log = logging.GetLogger()

//go:embed input.txt
var files embed.FS

type Voxel struct {
	X, Y, Z int
}

func getPartOneData() []Voxel {
	var voxels []Voxel
	for _, line := range util.ReadProblemInput(files, 1) {
		coords := util.Map(
			strings.Split(line, ","),
			func(x string) int { return util.AtoiNoError(x) },
		)
		voxels = append(voxels, Voxel{coords[0], coords[1], coords[2]})
	}
	return voxels
}

func GetAdjacent(v Voxel) []Voxel {
	return []Voxel{
		{v.X - 1, v.Y, v.Z},
		{v.X + 1, v.Y, v.Z},
		{v.X, v.Y - 1, v.Z},
		{v.X, v.Y + 1, v.Z},
		{v.X, v.Y, v.Z - 1},
		{v.X, v.Y, v.Z + 1},
	}
}

func GetNumAdjacentLava(v Voxel, lavaMap map[Voxel]bool) int {
	count := 0
	for _, adj := range GetAdjacent(v) {
		if _, ok := lavaMap[adj]; ok {
			count++
		}
	}
	return count
}

func PartOne() any {
	input := getPartOneData()
	lavaMap := map[Voxel]bool{}

	totalSides := 0

	for _, v := range input {
		// We subtract 2 for each adjacent lava, to account for over-counting done
		// by the adjacent lava previously
		addedSides := 6 - 2*GetNumAdjacentLava(v, lavaMap)
		totalSides += addedSides
		lavaMap[v] = true
	}

	return totalSides
}

type BoundingBox struct {
	MinX, MaxX, MinY, MaxY, MinZ, MaxZ int
}

func (b *BoundingBox) IsOutside(v Voxel) bool {
	return v.X < b.MinX || v.X > b.MaxX ||
		v.Y < b.MinY || v.Y > b.MaxY ||
		v.Z < b.MinZ || v.Z > b.MaxZ
}

func GetBoundingBox(voxels []Voxel) BoundingBox {
	box := BoundingBox{
		voxels[0].X, voxels[0].X,
		voxels[0].Y, voxels[0].Y,
		voxels[0].Z, voxels[0].Z,
	}

	for _, v := range voxels[1:] {
		if v.X < box.MinX {
			box.MinX = v.X
		}
		if v.X > box.MaxX {
			box.MaxX = v.X
		}
		if v.Y < box.MinY {
			box.MinY = v.Y
		}
		if v.Y > box.MaxY {
			box.MaxY = v.Y
		}
		if v.Z < box.MinZ {
			box.MinZ = v.Z
		}
		if v.Z > box.MaxZ {
			box.MaxZ = v.Z
		}
	}

	// Make bounding box 1 larger so that the lava
	// contents are fully contained by an air box
	box.MinX--
	box.MinY--
	box.MinZ--
	box.MaxX++
	box.MaxY++
	box.MaxZ++

	return box
}

func GetOpenAirVoxels(lavaMap map[Voxel]bool, box BoundingBox) map[Voxel]bool {
	// Bounding box edges are all air by construction
	start := Voxel{box.MinX, box.MinY, box.MinZ}
	openAirVoxels := map[Voxel]bool{start: true}
	toProcess := map[Voxel]bool{start: true}

	for len(toProcess) > 0 {
		nextOpenAirV := maps.Keys(toProcess)[0]
		delete(toProcess, nextOpenAirV)

		for _, adjV := range GetAdjacent(nextOpenAirV) {
			if _, ok := openAirVoxels[adjV]; ok {
				continue
			}
			if _, ok := lavaMap[adjV]; ok {
				continue
			}
			if box.IsOutside(adjV) {
				continue
			}
			openAirVoxels[adjV] = true
			toProcess[adjV] = true
		}
	}

	return openAirVoxels
}

func PartTwo() any {
	input := getPartOneData()
	box := GetBoundingBox(input)

	lavaMap := map[Voxel]bool{}
	for _, v := range input {
		lavaMap[v] = true
	}
	openAirMap := GetOpenAirVoxels(lavaMap, box)

	// Start by assuming all lava voxel sides will cool
	totalSides := 6 * len(lavaMap)

	for x := box.MinX; x <= box.MaxX; x++ {
		for y := box.MinY; y <= box.MaxY; y++ {
			for z := box.MinZ; z <= box.MaxZ; z++ {
				v := Voxel{x, y, z}
				numAdjacentLava := GetNumAdjacentLava(v, lavaMap)
				if _, ok := lavaMap[v]; ok {
					// is a lava voxel; reduce overcounting by adjacent
					// lava voxels' sides, same as part 1
					totalSides -= numAdjacentLava
					continue
				}
				if _, ok := openAirMap[v]; !ok {
					// is a trapped air voxel; reduce overcounting by adjacent
					// lava voxels' sides, as if the air was a lava voxel
					totalSides -= numAdjacentLava
				}
			}
		}
	}

	return totalSides
}

func main() {
	log.Infow(fmt.Sprintf("Answer: %v", PartOne()), "part", 1)
	log.Infow(fmt.Sprintf("Answer: %v", PartTwo()), "part", 2)
}
