package notes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseWithDash(t *testing.T) {

	notes := Parse("65-95")

	assert.Equal(t, 65, notes[0])
}

func TestParseWithCommas(t *testing.T) {

	notes := Parse("70,71,72,73,74")

	assert.Equal(t, 70, notes[0])
	assert.Equal(t, 71, notes[1])
}

func TestParseDrum(t *testing.T) {

	notes := Parse("DRUM")

	assert.Equal(t, []int{46, 48, 41, 58, 40, 49, 51, 42, 44, 39}, notes)
}
