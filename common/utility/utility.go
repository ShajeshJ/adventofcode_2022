package utility

import (
	"embed"
	"fmt"
	"strings"
)

func ReadProblemInput(fileReader embed.FS, part int) []string {
	bytes, _ := fileReader.ReadFile(fmt.Sprintf("part%v.txt", part))
	return strings.Split(string(bytes), "\n")
}
