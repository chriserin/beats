package projectfile

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {

	project := Parse(testfileName("simple_project.bp"))

	assert.NotNil(t, project)
	assert.Contains(t, project.Text, "PROJECT")

	assert.Equal(t, 2, len(project.Parts), project.Text)
}

func testfileName(fileName string) string {

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return dir + "/fixtures/" + fileName
}
