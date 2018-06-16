package projectfile

import (
	"os"
	"testing"

	"github.com/rakyll/portmidi"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {

	portmidi.Initialize()
	project := Parse(testfileName("simple_project.bp"))

	assert.NotNil(t, project)

	assert.Equal(t, 2, len(project.Parts))
	assert.Equal(t, 12719, int(project.Length))
}

func TestParseWithKeyPatternRepeats(t *testing.T) {

	portmidi.Initialize()
	project := Parse(testfileName("key_pattern_repeat_project.bp"))

	assert.NotNil(t, project)

	assert.Equal(t, 1, len(project.Parts))
	assert.Equal(t, 11999, int(project.Length))
}

func TestParseForTempo(t *testing.T) {

	portmidi.Initialize()
	project := Parse(testfileName("simple_project.bp"))

	assert.NotNil(t, project)

	assert.Equal(t, 113, int(project.Tempo))
	assert.Equal(t, 12719, int(project.Length))
}

func testfileName(fileName string) string {

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return dir + "/fixtures/" + fileName
}
