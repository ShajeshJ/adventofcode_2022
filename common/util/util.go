package util

import (
	"embed"
	"fmt"
	"strings"
)

func ReadProblemInput(fileCollection embed.FS, part int) []string {
	bytes, _ := fileCollection.ReadFile(fmt.Sprintf("part%v.txt", part))
	return strings.Split(string(bytes), "\n")
}
