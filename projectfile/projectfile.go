//Package projectfile supports parsing of project files
package projectfile

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"../grid"
	"../part"
	"../pattern_file"
	"github.com/rakyll/portmidi"
)

// Project struct
type Project struct {
	Parts        []part.Part
	ProjectLines []string
	Length       portmidi.Timestamp
	Tempo        grid.Tempo
}

// Options for project to pass around to patterns
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Parse a project file
func Parse(fileName string) Project {
	dir, _ := os.Getwd()
	defer os.Chdir(dir)
	projectDir := filepath.Dir(fileName)
	file, err := os.Open(fileName)
	os.Chdir(projectDir)
	check(err)

	projectLines, patternFiles, options := splitFileText(file)

	parsedTempo := parseTempo(options)

	projectOptions := grid.Options{Tempo: parsedTempo}

	parsedPatternFiles := map[string]patternfile.PatternFile{}

	for key, fileName := range patternFiles {
		parsedPatternFile := patternfile.Parse(fileName, projectOptions)
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
		Tempo:        parsedTempo,
	}

}

func parseTempo(options map[string]string) grid.Tempo {
	var parsedTempo int
	var err error

	if tempo, ok := options["TEMPO"]; ok {
		parsedTempo, err = strconv.Atoi(tempo)
		if err != nil {
			panic("Tempo must be an integer")
		}
	} else {
		return grid.Tempo(120)
	}

	return grid.Tempo(parsedTempo)
}

func splitFileText(file io.Reader) ([]string, map[string]string, map[string]string) {
	projectLines := []string{}
	patternFiles := map[string]string{}
	options := map[string]string{}

	scanner := bufio.NewScanner(file)

	postPattern := false
	for scanner.Scan() {
		line := scanner.Text()

		if postPattern {
			projectLines = append(projectLines, line)
		} else {
			patternFiles, options = parseProjectOption(line, patternFiles, options)
		}

		if strings.Contains(line, "PROJECT") {
			postPattern = true
		}
	}

	return projectLines, patternFiles, options
}

func findLastPartEventTimestamp(part part.Part) portmidi.Timestamp {
	return part.Events[len(part.Events)-1].Event.Timestamp
}

func parseProjectOption(line string, patternFiles map[string]string, options map[string]string) (map[string]string, map[string]string) {
	if strings.Contains(line, "=") {
		splitValues := strings.Split(line, "=")

		identifier := splitValues[0]
		value := splitValues[1]

		if len(identifier) == 1 {
			patternFiles[identifier] = value
		} else {
			options[identifier] = value
		}
	}

	return patternFiles, options
}
