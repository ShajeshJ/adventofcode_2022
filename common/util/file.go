package util

import (
	"embed"
	"fmt"
	"strings"
)

// ReadProblemInput reads and returns data from the file named "part{part}.txt" from the `fileCollection`
func ReadProblemInput(fileCollection embed.FS, part int) []string {
	bytes, _ := fileCollection.ReadFile(fmt.Sprintf("part%v.txt", part))
	return strings.Split(string(bytes), "\n")
}
