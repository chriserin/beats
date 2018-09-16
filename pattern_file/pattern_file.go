//Package patternfile operates on pattern files
package patternfile

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"../grid"
	"../notes"
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
func Parse(fileName string, projectOptions grid.Options) PatternFile {
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

	patternNotes := parseNotes(options)
	channel := parseChannel(options)
	partOptions := grid.Options{
		Tempo:      projectOptions.Tempo,
		DeviceName: options["DeviceName"],
		Notes:      patternNotes,
		Channel:    channel,
		Start:      portmidi.Time() + 500,
	}

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

func parseNotes(options map[string]string) []int {
	if notesOption, ok := options["Notes"]; ok {
		return notes.Parse(notesOption)
	}

	return []int{60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83}
}

func parseChannel(options map[string]string) portmidi.Channel {
	if channelOption, ok := options["Channel"]; ok {
		if channelNum, err := strconv.Atoi(channelOption); err == nil {
			if channelNum >= 0 && channelNum < 15 {
				return portmidi.Channel(channelNum)
			}
		}
	}

	return portmidi.Channel(-1)
}
