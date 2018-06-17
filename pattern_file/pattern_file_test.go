package patternfile

import (
	"os"
	"regexp"
	"testing"

	"../grid"
	"github.com/rakyll/portmidi"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {

	pFile := Parse(testfileName("simple_pattern.pf"), grid.Options{Tempo: grid.Tempo(120), DeviceName: "beats"})

	assert.NotNil(t, pFile)
	assert.Contains(t, pFile.Text, "PATTERN")

	gridTextLinesCount := len(regexp.MustCompile("\n").Split(pFile.GridText, -1))
	assert.Equal(t, 12, gridTextLinesCount, pFile.GridText)
	assert.Equal(t, int64(60), pFile.MidiPoints[0].Event.Data1)
}

func TestParseNotes(t *testing.T) {
	pFile := Parse(testfileName("defined_notes_pattern.pf"), grid.Options{Tempo: grid.Tempo(120), DeviceName: "beats"})

	assert.NotNil(t, pFile)
	assert.Contains(t, pFile.Text, "PATTERN")

	gridTextLinesCount := len(regexp.MustCompile("\n").Split(pFile.GridText, -1))
	assert.Equal(t, 12, gridTextLinesCount, pFile.GridText)
	assert.Equal(t, int64(60), pFile.MidiPoints[0].Event.Data1)
}

func TestParseChannel(t *testing.T) {
	pFile := Parse(testfileName("channeled_pattern.pf"), grid.Options{})

	assert.Equal(t, portmidi.Channel(4), pFile.MidiPoints[0].Channel)
}

func testfileName(fileName string) string {

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return dir + "/fixtures/" + fileName
}
