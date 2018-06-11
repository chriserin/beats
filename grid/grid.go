package grid

import (
	"fmt"
	"strings"

	"github.com/rakyll/portmidi"
)

// TextGridPoint contains the coordinates and symbol
type TextGridPoint struct {
	X      int
	Y      int
	Symbol rune
}

//TransformGridToMidi transforms a grid into midi notes
func TransformGridToMidi(gridText string) []portmidi.Event {
	textGridPoints := TransformTextGrid(gridText)
	rawPitchPoints := TransformTextGridPoints(textGridPoints, "down")
	beatPoints := TransformRawPitchPoints(rawPitchPoints)
	pitchPoints := TransformBeatPoints(beatPoints, []int{62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73})
	timedPoints := TransformPitchPoints(pitchPoints, 120)
	fmt.Println(timedPoints)
	midiPoints := TransformTimedPoints(timedPoints)
	return midiPoints
}

// TransformTextGrid the text representation of Grid into list of GridPoints
func TransformTextGrid(gridText string) []TextGridPoint {
	results := []TextGridPoint{}

	lines := strings.Split(gridText, "\n")
	for x, line := range lines {
		for y, r := range line {
			if r > 32 {
				results = append(results, TextGridPoint{x, y, r})
			}
		}
	}

	return results
}

//RawPitchPoint struct of data before pitches converted to scale
type RawPitchPoint struct {
	RawPitch int
	Beat     int
	Symbol   rune
}

//TransformTextGridPoints converts XY grid points to Pitch Beat grid points for either left2right or up2down
func TransformTextGridPoints(points []TextGridPoint, direction string) []RawPitchPoint {
	results := make([]RawPitchPoint, len(points))

	for i, point := range points {
		if direction == "down" {
			results[i] = RawPitchPoint{point.X, point.Y, point.Symbol}
		} else {
			results[i] = RawPitchPoint{point.Y, point.X, point.Symbol}
		}
	}

	return results
}

//BeatPoint point struct with Length added
type BeatPoint struct {
	RawPitch int
	Beat     int
	Length   int
	Symbol   rune
}

//TransformRawPitchPoints determines the length of each beat point
//If the next symbol on the same pitch is a `!` then the length goes to the end of that beat
//If the next symbol on the same pitch is a character, then the length is 1
//If there is no next symbol then the length is 1
func TransformRawPitchPoints(points []RawPitchPoint) []BeatPoint {
	results := []BeatPoint{}

	//iterate to first symbol
	//find next symbol of same pitch
	for i, point := range points {
		if point.Symbol != '!' {
			beatPoint := BeatPoint{
				RawPitch: point.RawPitch,
				Beat:     point.Beat,
				Symbol:   point.Symbol,
			}

			nextPoint, found := findNextNoteWithPitch(points[i+1:], point.RawPitch)

			if !found {
				beatPoint.Length = 1
			} else if nextPoint.Symbol == '!' {
				beatPoint.Length = nextPoint.Beat - beatPoint.Beat + 1
			} else {
				beatPoint.Length = 1
			}

			results = append(results, beatPoint)
		}
	}

	return results
}

func findNextNoteWithPitch(points []RawPitchPoint, rawPitch int) (RawPitchPoint, bool) {
	for _, point := range points {
		if point.RawPitch == rawPitch {
			return point, true
		}
	}

	return RawPitchPoint{}, false
}

//PitchPoint is a note described as a specific pitch on a certain beat
type PitchPoint struct {
	Pitch  int
	Beat   int
	Length int
	Symbol rune
}

//TransformBeatPoints transfroms the grid points into pitch points
func TransformBeatPoints(points []BeatPoint, pitches []int) []PitchPoint {
	results := make([]PitchPoint, len(points))

	for i, point := range points {
		results[i] = PitchPoint{pitches[point.RawPitch], point.Beat, point.Length, point.Symbol}
	}

	return results
}

//Milliseconds represent the exact time the event will occur
type Milliseconds int

//TimedPoint is the event information with time relative to 0
type TimedPoint struct {
	Pitch  int
	Start  Milliseconds
	Length Milliseconds
	Symbol rune
}

//Tempo used to convert beat to milliseconds
type Tempo int

//TransformPitchPoints uses tempo to convert beats to seconds
func TransformPitchPoints(points []PitchPoint, tempo Tempo) []TimedPoint {
	results := make([]TimedPoint, len(points))

	millisecondsPerBeat := 1000 / (float64(tempo) / 60)

	for i, point := range points {
		start := int(point.Beat) * int(millisecondsPerBeat)
		length := int(point.Length)*int(millisecondsPerBeat) - 1
		results[i] = TimedPoint{point.Pitch, Milliseconds(start), Milliseconds(length), point.Symbol}
	}

	return results
}

//TransformTimedPoints transforms timed points into portmidi Events
func TransformTimedPoints(points []TimedPoint) []portmidi.Event {
	results := make([]portmidi.Event, len(points)*2)

	fmt.Println("portmidi time")
	fmt.Println(portmidi.Time())

	for i, point := range points {
		results[i*2] = portmidi.Event{Timestamp: portmidi.Time() + portmidi.Timestamp(point.Start), Status: 0x90, Data1: int64(point.Pitch), Data2: 100}
		results[i*2+1] = portmidi.Event{Timestamp: portmidi.Time() + portmidi.Timestamp(point.Start+point.Length), Status: 0x80, Data1: int64(point.Pitch), Data2: 100}
	}

	return results
}
