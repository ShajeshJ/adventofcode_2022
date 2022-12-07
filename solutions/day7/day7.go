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

// type File struct {
// 	name string
// 	size int
// }

type Dir struct {
	Parent *Dir
	Name   string
	// files []File
	Size int
	Dirs map[string]*Dir
}

var (
	cdCmd      = regexp.MustCompile(`\$ cd (.+)`)
	dirOutput  = regexp.MustCompile(`dir (.+)`)
	fileOutput = regexp.MustCompile(`(\d+) (.+)`)
)

func traverseDir(finishedDirCallback func(cwd *Dir)) *Dir {
	terminal := util.ReadProblemInput(files, 1)

	cwd := &Dir{Name: "/", Dirs: map[string]*Dir{}}

	finishDir := func() {
		// Finished cwd; add sub-dir sizes to cwd before moving up
		for _, d := range cwd.Dirs {
			cwd.Size += d.Size
		}
		finishedDirCallback(cwd)
	}

	for i := 0; i < len(terminal); i++ {
		if res := cdCmd.FindStringSubmatch(terminal[i]); len(res) > 1 {
			if res[1] == ".." {
				finishDir()
				cwd = cwd.Parent
			} else if nextDir, ok := cwd.Dirs[res[1]]; ok {
				cwd = nextDir
			}
			continue
		}

		// Otherwise we process an ls command output
		for {
			i++
			if i >= len(terminal) {
				break // EOF
			} else if res := fileOutput.FindStringSubmatch(terminal[i]); len(res) > 1 {
				cwd.Size += util.AtoiNoError(res[1])
			} else if res := dirOutput.FindStringSubmatch(terminal[i]); len(res) > 1 {
				cwd.Dirs[res[1]] = &Dir{Name: res[1], Parent: cwd, Dirs: map[string]*Dir{}}
			} else {
				break // ls output complete
			}
		}
		i-- // Go back 1 so command isn't skipped in next iter
	}

	// Calculate leftovers all the way up to root dir
	for cwd.Parent != nil {
		finishDir()
		cwd = cwd.Parent
	}

	finishDir() // Calculate root dir
	return cwd
}

func PartOne() any {
	const maxDirSize = 100_000
	totalSizeOfUnder100k := 0
	traverseDir(func(cwd *Dir) {
		if cwd.Size <= maxDirSize {
			totalSizeOfUnder100k += cwd.Size
		}
		return
	})
	return totalSizeOfUnder100k
}

func PartTwo() any {
	allDirs := []*Dir{}

	root := traverseDir(func(cwd *Dir) {
		allDirs = append(allDirs, cwd)
	})

	const (
		totalDisk    = 70_000_000
		requiredDisk = 30_000_000
	)

	needToClear := requiredDisk - (totalDisk - root.Size)
	candidateDir := root // Worst case, we delete the whole root dir

	for _, d := range allDirs {
		if d.Size >= needToClear && d.Size < candidateDir.Size {
			candidateDir = d
		}
	}

	return candidateDir.Size
}

func main() {
	log.Infow(fmt.Sprintf("Answer: %v", PartOne()), "part", 1)
	log.Infow(fmt.Sprintf("Answer: %v", PartTwo()), "part", 2)
}
