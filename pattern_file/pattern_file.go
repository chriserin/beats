//Package patternfile operates on pattern files
package patternfile

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"../grid"
	"github.com/rakyll/portmidi"
)

//PatternFile contains all data defined in a pattern file
type PatternFile struct {
	DeviceName string
	Text       string
	GridText   string
	MidiPoints []grid.MidiPoint
	Length     portmidi.Timestamp
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
func Parse(fileName string, projectOptions grid.PartOptions) PatternFile {
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

	partOptions := grid.PartOptions{Tempo: projectOptions.Tempo, DeviceName: options["DeviceName"]}

	gridText := strings.Join(patternLines, "\n")
	midiPoints := grid.TransformGridToMidi(gridText, partOptions)

	return PatternFile{
		DeviceName: options["DeviceName"],
		Text:       strings.Join(lines, "\n"),
		GridText:   gridText,
		MidiPoints: midiPoints,
		Length:     length(options, midiPoints),
	}
}

func length(options map[string]string, midiPoints []grid.MidiPoint) portmidi.Timestamp {
	if lengthOption, ok := options["Length"]; ok {
		if length, err := strconv.Atoi(lengthOption); err == nil {
			millisecondsPerBeat := 1000 / (float64(120) / 60)
			return portmidi.Timestamp(float64(length) * millisecondsPerBeat)
		}
	}

	return midiPoints[len(midiPoints)-1].Event.Timestamp + 1
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
