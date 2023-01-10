package main

import (
	"embed"
	"fmt"
	"strings"

	"github.com/ShajeshJ/adventofcode_2022/common/logging"
	"github.com/ShajeshJ/adventofcode_2022/common/util"
)

var log = logging.GetLogger()

//go:embed input.txt
var files embed.FS

type Dir struct {
	Name   string
	Parent *Dir
	Dirs   map[string]*Dir
	Size   int
}

type Computer struct {
	Root               *Dir
	Cwd                *Dir
	FolderExitCallback func(c *Dir)
}

func NewComputer(folderExitCallback func(c *Dir)) Computer {
	return Computer{
		Root:               &Dir{"/", nil, map[string]*Dir{}, 0},
		Cwd:                nil,
		FolderExitCallback: folderExitCallback,
	}
}

func (c *Computer) Mkdir(name string) {
	c.Cwd.Dirs[name] = &Dir{name, c.Cwd, map[string]*Dir{}, 0}
}

func (c *Computer) FinalizeCwd() {
	for _, d := range c.Cwd.Dirs {
		c.Cwd.Size += d.Size
	}
	c.FolderExitCallback(c.Cwd)
}

func (c *Computer) RunCmd(tokens ...string) {
	if tokens[0] == "ls" {
		return
	}

	// else command is "cd"
	if tokens[1] == "/" {
		c.Cwd = c.Root
		return
	}
	if tokens[1] != ".." {
		c.Cwd = c.Cwd.Dirs[tokens[1]]
		return
	}

	// exiting the folder; do exit-folder processing
	c.FinalizeCwd()
	c.Cwd = c.Cwd.Parent
}

func (c *Computer) ReadlsOutput(tokens ...string) {
	if tokens[0] == "dir" {
		c.Mkdir(tokens[1])
	} else {
		// else is a file
		c.Cwd.Size += util.AtoiNoError(tokens[0])
	}
}

func BuildComputer(folderExitCallback func(cwd *Dir)) Computer {
	c := NewComputer(folderExitCallback)

	for _, line := range util.ReadProblemInput(files) {
		tokens := strings.Split(line, " ")

		if tokens[0] == "$" {
			c.RunCmd(tokens[1:]...)
			continue
		}

		// Otherwise we're reading ls output
		c.ReadlsOutput(tokens...)
	}

	// Processing remaining folders back up to the root
	for c.Cwd.Parent != nil {
		c.RunCmd("cd", "..")
	}
	c.FinalizeCwd() // Also process root itself

	return c
}

func PartOne() any {
	const maxDirSize = 100_000
	totalSizeOfUnder100k := 0
	BuildComputer(func(cwd *Dir) {
		if cwd.Size <= maxDirSize {
			totalSizeOfUnder100k += cwd.Size
		}
		return
	})
	return totalSizeOfUnder100k
}

func PartTwo() any {
	allDirs := []*Dir{}

	root := BuildComputer(func(cwd *Dir) {
		allDirs = append(allDirs, cwd)
	}).Root

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
