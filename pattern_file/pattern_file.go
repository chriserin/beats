//Package patternfile operates on pattern files
package patternfile

import (
	"bufio"
	"os"
	"strings"
)

//PatternFile contains all data defined in a pattern file
type PatternFile struct {
	DeviceName string
	Text       string
	GridText   string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

//Parse takes a filname, reads, parses and structures all the data in a file
func Parse(fileName string) PatternFile {
	lines := []string{}
	patternLines := []string{}
	postPattern := false

	file, err := os.Open(fileName)
	check(err)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		lines = append(lines, line)

		if postPattern {
			patternLines = append(patternLines, line)
		}

		if strings.Contains(line, "PATTERN") {
			postPattern = true
		}
	}

	return PatternFile{
		DeviceName: "Bus",
		Text:       strings.Join(lines, "\n"),
		GridText:   strings.Join(patternLines, "\n"),
	}
}
