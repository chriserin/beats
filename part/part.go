package part

import (
	"../grid"
	"../pattern_file"
	"github.com/rakyll/portmidi"
)

//Part struct
type Part struct {
	Patterns []patternfile.PatternFile
	Events   []grid.MidiPoint
}

// Parse lines for part
func Parse(projectLine string, parsedPatternFiles map[string]patternfile.PatternFile, endOfLastPart portmidi.Timestamp) Part {
	newPart := Part{Patterns: []patternfile.PatternFile{}}

	for _, key := range projectLine {
		pattern := parsedPatternFiles[string(key)]
		newPart.Patterns = append(newPart.Patterns, pattern)
		newPart.Events = append(newPart.Events, grid.ShiftEvents(pattern.MidiPoints, endOfLastPart)...)
	}

	return newPart
}
