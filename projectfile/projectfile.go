//Package projectfile supports parsing of project files
package projectfile

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"

	"../part"
	"../pattern_file"
	"github.com/rakyll/portmidi"
)

// Project struct
type Project struct {
	Parts        []part.Part
	ProjectLines []string
	Length       portmidi.Timestamp
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Parse a project file
func Parse(fileName string) Project {
	projectDir := filepath.Dir(fileName)
	file, err := os.Open(fileName)
	os.Chdir(projectDir)
	check(err)

	projectLines, patternFiles := splitFileText(file)

	parsedPatternFiles := map[string]patternfile.PatternFile{}

	for key, fileName := range patternFiles {
		parsedPatternFile := patternfile.Parse(fileName)
		parsedPatternFiles[key] = parsedPatternFile
	}

	parts := make([]part.Part, len(projectLines))

	endOfLastPart := portmidi.Timestamp(0)
	for i, projectLine := range projectLines {

		newPart := part.Parse(projectLine, parsedPatternFiles, endOfLastPart)
		endOfLastPart = findLastPartEventTimestamp(newPart) + 1
		parts[i] = newPart
	}

	return Project{
		ProjectLines: projectLines,
		Parts:        parts,
		Length:       endOfLastPart - 1,
	}
}

func splitFileText(file io.Reader) ([]string, map[string]string) {
	projectLines := []string{}
	patternFiles := map[string]string{}

	scanner := bufio.NewScanner(file)

	postPattern := false
	for scanner.Scan() {
		line := scanner.Text()

		if postPattern {
			projectLines = append(projectLines, line)
		} else {
			patternFiles = parseFileDeclaration(patternFiles, line)
		}

		if strings.Contains(line, "PROJECT") {
			postPattern = true
		}
	}

	return projectLines, patternFiles
}

func findLastPartEventTimestamp(part part.Part) portmidi.Timestamp {
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
