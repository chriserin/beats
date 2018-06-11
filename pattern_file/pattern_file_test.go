package patternfile

import (
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {

	pFile := Parse(testfileName("simple_pattern.pf"))

	assert.NotNil(t, pFile)
	assert.Contains(t, pFile.Text, "PATTERN")

	gridTextLinesCount := len(regexp.MustCompile("\n").Split(pFile.GridText, -1))
	assert.Equal(t, 12, gridTextLinesCount, pFile.GridText)
}

func testfileName(fileName string) string {

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return dir + "/fixtures/" + fileName
}
