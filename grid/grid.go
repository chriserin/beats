package grid

import (
	"strings"

	"../devices"
	"../velocity"
	"github.com/rakyll/portmidi"
)

// TextGridPoint contains the coordinates and symbol
type TextGridPoint struct {
	X      int
	Y      int
	Symbol rune
}

//Options passed to grid transformer
type Options struct {
	Tempo      Tempo
	DeviceName string
	Notes      []int
	Channel    portmidi.Channel
	Start      portmidi.Timestamp
}

//TransformGridToMidi transforms a grid into midi notes
func TransformGridToMidi(gridText string, options Options) []MidiPoint {
	textGridPoints := TransformTextGrid(gridText)
	rawPitchPoints := TransformTextGridPoints(textGridPoints, "down")
	beatPoints := TransformRawPitchPoints(rawPitchPoints)
	pitchPoints := TransformBeatPoints(beatPoints, options.Notes)
	timedPoints := TransformPitchPoints(pitchPoints, options.Tempo)
	velocityPoints := TransformTimedPoints(timedPoints)
	midiPoints := TransformVelocityPoints(velocityPoints, options.Start)
	midiPoints = SetDeviceID(midiPoints, options.DeviceName, options.Channel)

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

//VelocityPoint is the event informattion with the Symbol converted to a velocity
type VelocityPoint struct {
	Pitch    int
	Start    Milliseconds
	Length   Milliseconds
	Velocity int
	Symbol   rune
}

//TransformTimedPoints transforms velocity points
func TransformTimedPoints(points []TimedPoint) []VelocityPoint {
	results := make([]VelocityPoint, len(points))

	for i, point := range points {
		point := VelocityPoint{
			Pitch:    point.Pitch,
			Start:    point.Start,
			Length:   point.Length,
			Velocity: velocity.ConvertSymbol(point.Symbol),
			Symbol:   point.Symbol,
		}

		results[i] = point
	}

	return results
}

//MidiPoint containts a midi event and the device ID (and later the channel)
type MidiPoint struct {
	Event    portmidi.Event
	DeviceID portmidi.DeviceID
	Channel  portmidi.Channel
}

//TransformVelocityPoints transforms velocity points into portmidi Events
func TransformVelocityPoints(points []VelocityPoint, startTime portmidi.Timestamp) []MidiPoint {
	results := make([]MidiPoint, len(points)*2)

	for i, point := range points {
		startEvent := portmidi.Event{Timestamp: startTime + portmidi.Timestamp(point.Start), Status: 0x90, Data1: int64(point.Pitch), Data2: int64(point.Velocity)}
		endEvent := portmidi.Event{Timestamp: startTime + portmidi.Timestamp(point.Start+point.Length), Status: 0x80, Data1: int64(point.Pitch), Data2: int64(point.Velocity)}
		results[i*2] = MidiPoint{Event: startEvent}
		results[i*2+1] = MidiPoint{Event: endEvent}
	}

	return results
}

//SetDeviceID sets the device ID on all midi points
func SetDeviceID(points []MidiPoint, name string, channel portmidi.Channel) []MidiPoint {
	results := make([]MidiPoint, len(points))

	deviceID := devices.FindDeviceID(name)

	for i, point := range points {
		results[i] = MidiPoint{Event: point.Event, DeviceID: deviceID, Channel: channel}
	}

	return results
}

//ShiftEvents moves events by the shiftAmount
func ShiftEvents(points []MidiPoint, shiftAmount portmidi.Timestamp) []MidiPoint {
	results := make([]MidiPoint, len(points))

	for i, point := range points {
		event := point.Event
		newEvent := portmidi.Event{Timestamp: event.Timestamp + shiftAmount, Status: event.Status, Data1: event.Data1, Data2: event.Data2}
		point.Event = newEvent
		results[i] = point
	}

	return results
}
