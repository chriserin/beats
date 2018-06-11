//Package projectfile supports parsing of project files
package projectfile

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"../grid"
	"../pattern_file"
	"github.com/rakyll/portmidi"
)

// Project struct
type Project struct {
	Text         string
	Parts        []Part
	ProjectLines []string
	Length       portmidi.Timestamp
}

// Part struct
type Part struct {
	Patterns []patternfile.PatternFile
	Events   []grid.MidiPoint
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Parse a project file
func Parse(fileName string) Project {
	lines := []string{}
	projectLines := []string{}
	postPattern := false

	patternFiles := map[string]string{}

	projectDir := filepath.Dir(fileName)
	file, err := os.Open(fileName)
	os.Chdir(projectDir)
	check(err)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		lines = append(lines, line)

		if postPattern {
			projectLines = append(projectLines, line)
		} else {
			patternFiles = parseFileDeclaration(patternFiles, line)
		}

		if strings.Contains(line, "PROJECT") {
			postPattern = true
		}
	}

	parsedPatternFiles := map[string]patternfile.PatternFile{}

	for key, fileName := range patternFiles {
		parsedPatternFile := patternfile.Parse(fileName)
		parsedPatternFiles[key] = parsedPatternFile
	}

	parts := make([]Part, len(projectLines))

	endOfLastPart := portmidi.Timestamp(0)
	for i, projectLine := range projectLines {
		newPart := Part{Patterns: []patternfile.PatternFile{}}

		for _, key := range projectLine {
			pattern := parsedPatternFiles[string(key)]
			newPart.Patterns = append(newPart.Patterns, pattern)
			newPart.Events = append(newPart.Events, grid.ShiftEvents(pattern.MidiPoints, endOfLastPart)...)
		}

		endOfLastPart = findLastPartEventTimestamp(newPart) + 1

		parts[i] = newPart
	}

	return Project{
		Text:         strings.Join(lines, "\n"),
		ProjectLines: projectLines,
		Parts:        parts,
		Length:       findLastProjectEventTimestamp(parts),
	}
}

func findLastProjectEventTimestamp(parts []Part) portmidi.Timestamp {
	var lastTimestamp = portmidi.Timestamp(0)
	for _, part := range parts {
		partEndTimestamp := findLastPartEventTimestamp(part)
		if partEndTimestamp > lastTimestamp {
			lastTimestamp = partEndTimestamp
		}
	}
	return lastTimestamp
}

func findLastPartEventTimestamp(part Part) portmidi.Timestamp {
	return part.Events[len(part.Events)-1].Event.Timestamp
}

func parseFileDeclaration(patternFiles map[string]string, line string) map[string]string {
	if strings.Contains(line, "=") {
		splitValues := strings.Split(line, "=")

		identifier := splitValues[0]
		filePath := splitValues[1]

		patternFiles[identifier] = filePath
	}

	return patternFiles
}
