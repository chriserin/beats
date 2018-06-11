//Package patternfile operates on pattern files
package patternfile

import (
	"bufio"
	"os"
	"strings"

	"../grid"
)

//PatternFile contains all data defined in a pattern file
type PatternFile struct {
	DeviceName string
	Text       string
	GridText   string
	MidiPoints []grid.MidiPoint
}

type option struct {
	key   string
	value string
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
	options := make(map[string]string)

	file, err := os.Open(fileName)
	check(err)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		lines = append(lines, line)

		if postPattern {
			patternLines = append(patternLines, line)
		} else {
			appendOption(options, parseForConfiguration(line))
		}

		if strings.Contains(line, "PATTERN") {
			postPattern = true
		}
	}

	gridText := strings.Join(patternLines, "\n")
	midiPoints := grid.TransformGridToMidi(gridText, options)

	return PatternFile{
		DeviceName: options["DeviceName"],
		Text:       strings.Join(lines, "\n"),
		GridText:   gridText,
		MidiPoints: midiPoints,
	}
}

func appendOption(options map[string]string, parsedOption *option) {
	if parsedOption != nil {
		options[parsedOption.key] = parsedOption.value
	}
}

func parseForConfiguration(line string) *option {
	if strings.Contains(line, "=") {
		splitValues := strings.Split(line, "=")

		key := splitValues[0]
		value := splitValues[1]

		return &option{key, value}
	}

	return nil
}
