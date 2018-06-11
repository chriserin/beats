//Package projectfile supports parsing of project files
package projectfile

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"../pattern_file"
)

// Project struct
type Project struct {
	Text         string
	Parts        []Part
	ProjectLines []string
}

// Part struct
type Part struct {
	Patterns []patternfile.PatternFile
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

	file, err := os.Open(fileName)
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

	for i, projectLine := range projectLines {
		newPart := Part{Patterns: []patternfile.PatternFile{}}

		for _, key := range projectLine {
			newPart.Patterns = append(newPart.Patterns, parsedPatternFiles[string(key)])
		}

		parts[i] = newPart
	}

	fmt.Println("Returning Project")
	fmt.Println(lines)

	return Project{
		Text:         strings.Join(lines, "\n"),
		ProjectLines: projectLines,
		Parts:        parts,
	}
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
