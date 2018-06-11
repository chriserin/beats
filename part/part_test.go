package part

import (
	"testing"

	"../pattern_file"

	"github.com/rakyll/portmidi"
	"github.com/stretchr/testify/assert"
)

func TestParseKeyPatternRepeat(t *testing.T) {
	parsedPatterns := make(map[string]patternfile.PatternFile)
	parsedPatterns["A"] = patternfile.PatternFile{}
	part := Parse("2AA", parsedPatterns, portmidi.Timestamp(0))

	assert.Equal(t, 2, part.KeyPatternRepeats)
	assert.Equal(t, 2, len(part.Patterns))
}

func TestParseKeyPatternRepeatDefault(t *testing.T) {
	parsedPatterns := make(map[string]patternfile.PatternFile)
	parsedPatterns["A"] = patternfile.PatternFile{}
	part := Parse("AA", parsedPatterns, portmidi.Timestamp(0))

	assert.Equal(t, 1, part.KeyPatternRepeats)
	assert.Equal(t, 2, len(part.Patterns))
}
