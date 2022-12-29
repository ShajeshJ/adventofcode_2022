package util

import (
	"embed"
	"strings"
)

// ReadProblemInput reads and returns data from the file named "input.txt" from the `fileCollection`
func ReadProblemInput(fileCollection embed.FS, part int) []string {
	bytes, _ := fileCollection.ReadFile("input.txt")
	return strings.Split(string(bytes), "\n")
}
