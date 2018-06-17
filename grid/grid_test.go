package grid

import (
	"strings"
	"testing"

	"github.com/rakyll/portmidi"
	"github.com/stretchr/testify/assert"
)

func TestTransformTextGrd(t *testing.T) {

	textGrid := `
	X
	 X
	  X
	`

	textGrid = strings.Replace(textGrid, "\t", "", -1)
	textGrid = strings.Trim(textGrid, "\n")

	result := TransformTextGrid(textGrid)

	point1 := TextGridPoint{0, 0, 'X'}
	point2 := TextGridPoint{1, 1, 'X'}
	point3 := TextGridPoint{2, 2, 'X'}

	assert.Equal(t, point1, result[0])
	assert.Equal(t, point2, result[1])
	assert.Equal(t, point3, result[2])
}

func TestTransformTextGridPoints(t *testing.T) {

	points := []TextGridPoint{TextGridPoint{1, 0, 'X'}, TextGridPoint{2, 1, 'X'}, TextGridPoint{3, 2, 'X'}}

	results := TransformTextGridPoints(points, "down")

	point1 := RawPitchPoint{1, 0, 'X'}
	point2 := RawPitchPoint{2, 1, 'X'}
	point3 := RawPitchPoint{3, 2, 'X'}

	assert.Equal(t, 3, len(results))

	assert.Equal(t, point1, results[0])
	assert.Equal(t, point2, results[1])
	assert.Equal(t, point3, results[2])
}

func TestTransformTextGridPointsLeft2Right(t *testing.T) {

	points := []TextGridPoint{TextGridPoint{0, 1, 'X'}, TextGridPoint{1, 2, 'X'}, TextGridPoint{2, 3, 'X'}}

	results := TransformTextGridPoints(points, "right")

	point1 := RawPitchPoint{1, 0, 'X'}
	point2 := RawPitchPoint{2, 1, 'X'}
	point3 := RawPitchPoint{3, 2, 'X'}

	assert.Equal(t, 3, len(results))

	assert.Equal(t, point1, results[0])
	assert.Equal(t, point2, results[1])
	assert.Equal(t, point3, results[2])
}

func TestTransformRawPitchPoints(t *testing.T) {
	points := []RawPitchPoint{{0, 0, 'X'}, {1, 1, 'X'}, {2, 2, 'X'}}

	results := TransformRawPitchPoints(points)

	point1 := BeatPoint{0, 0, 1, 'X'}
	point2 := BeatPoint{1, 1, 1, 'X'}
	point3 := BeatPoint{2, 2, 1, 'X'}

	assert.Equal(t, 3, len(results))

	assert.Equal(t, point1, results[0])
	assert.Equal(t, point2, results[1])
	assert.Equal(t, point3, results[2])
}

func TestTransformRawPitchPointsWithStops(t *testing.T) {
	points := []RawPitchPoint{{0, 0, 'X'}, {0, 2, '!'}}

	results := TransformRawPitchPoints(points)

	point1 := BeatPoint{0, 0, 3, 'X'}

	assert.Equal(t, 1, len(results))

	assert.Equal(t, point1, results[0])
}

func TestTransformRawPitchPointsWithNextNote(t *testing.T) {
	points := []RawPitchPoint{{0, 0, 'X'}, {0, 2, 'X'}, {0, 13, '!'}}

	results := TransformRawPitchPoints(points)

	point1 := BeatPoint{0, 0, 1, 'X'}

	assert.Equal(t, 2, len(results))

	assert.Equal(t, point1, results[0])
}

func TestTransformBeatPoints(t *testing.T) {
	points := []BeatPoint{{0, 0, 1, 'X'}, {1, 1, 1, 'X'}, {2, 2, 1, 'X'}}

	results := TransformBeatPoints(points, []int{41, 42, 43})

	point1 := PitchPoint{41, 0, 1, 'X'}
	point2 := PitchPoint{42, 1, 1, 'X'}
	point3 := PitchPoint{43, 2, 1, 'X'}

	assert.Equal(t, 3, len(results))

	assert.Equal(t, point1, results[0])
	assert.Equal(t, point2, results[1])
	assert.Equal(t, point3, results[2])
}

func TestTransformPitchPoints(t *testing.T) {
	points := []PitchPoint{{42, 0, 1, 'X'}, {43, 1, 1, 'X'}, {44, 2, 1, 'X'}}
	tempo := 120

	results := TransformPitchPoints(points, Tempo(tempo))

	point1 := TimedPoint{42, 0, 499, 'X'}
	point2 := TimedPoint{43, 500, 499, 'X'}
	point3 := TimedPoint{44, 1000, 499, 'X'}

	assert.Equal(t, point1, results[0])
	assert.Equal(t, point2, results[1])
	assert.Equal(t, point3, results[2])
}

func TestTransformTimedPoints(t *testing.T) {
	points := []TimedPoint{{42, 0, 499, 'X'}, {43, 500, 499, 'X'}, {44, 1000, 499, 'X'}}

	results := TransformTimedPoints(points)

	point1 := VelocityPoint{42, 0, 499, 120, 'X'}
	point2 := VelocityPoint{43, 500, 499, 120, 'X'}
	point3 := VelocityPoint{44, 1000, 499, 120, 'X'}

	assert.Equal(t, point1, results[0])
	assert.Equal(t, point2, results[1])
	assert.Equal(t, point3, results[2])
}

func TestTransformVelocityPoints(t *testing.T) {
	points := []VelocityPoint{{42, 0, 499, 100, 'X'}, {43, 500, 499, 100, 'X'}, {44, 1000, 499, 100, 'X'}}

	results := TransformVelocityPoints(points)

	event1 := MidiPoint{Event: portmidi.Event{Timestamp: portmidi.Time(), Status: 0x90, Data1: 42, Data2: 100}}
	event1End := MidiPoint{Event: portmidi.Event{Timestamp: portmidi.Time() + 499, Status: 0x80, Data1: 42, Data2: 100}}
	event2 := MidiPoint{Event: portmidi.Event{Timestamp: portmidi.Time() + 500, Status: 0x90, Data1: 43, Data2: 100}}
	event2End := MidiPoint{Event: portmidi.Event{Timestamp: portmidi.Time() + 500 + 499, Status: 0x80, Data1: 43, Data2: 100}}
	event3 := MidiPoint{Event: portmidi.Event{Timestamp: portmidi.Time() + 1000, Status: 0x90, Data1: 44, Data2: 100}}
	event3End := MidiPoint{Event: portmidi.Event{Timestamp: portmidi.Time() + 1000 + 499, Status: 0x80, Data1: 44, Data2: 100}}

	assert.Equal(t, event1, results[0])
	assert.Equal(t, event1End, results[1])
	assert.Equal(t, event2, results[2])
	assert.Equal(t, event2End, results[3])
	assert.Equal(t, event3, results[4])
	assert.Equal(t, event3End, results[5])
}
