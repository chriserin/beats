package part

import (
	"strconv"

	"../grid"
	"../pattern_file"
	"github.com/rakyll/portmidi"
)

//Part struct
type Part struct {
	KeyPatternRepeats int
	Patterns          []patternfile.PatternFile
	Events            []grid.MidiPoint
}

// Parse lines for part
func Parse(projectLine string, parsedPatternFiles map[string]patternfile.PatternFile, endOfLastPart portmidi.Timestamp) Part {
	newPart := Part{Patterns: []patternfile.PatternFile{}, KeyPatternRepeats: 1}

	for _, r := range projectLine {
		key := string(r)
		if repeats, err := strconv.Atoi(key); err == nil {
			newPart.KeyPatternRepeats = repeats
			continue
		}

		pattern := parsedPatternFiles[key]
		newPart.Patterns = append(newPart.Patterns, pattern)
	}

	//Shift events based on start time of part
	for patternIndex, pattern := range newPart.Patterns {
		newPart.Events = append(newPart.Events, grid.ShiftEvents(pattern.MidiPoints, endOfLastPart)...)

		for i := 1; patternIndex == 0 && i < newPart.KeyPatternRepeats; i++ {
			shiftAmount := endOfLastPart + (pattern.Length * (portmidi.Timestamp(i)))
			newPart.Events = append(newPart.Events, grid.ShiftEvents(pattern.MidiPoints, shiftAmount)...)
		}
	}

	return newPart
}
