package runner

import (
	"fmt"
	"log"
	"sort"

	"../projectfile"

	"../grid"
	"github.com/rakyll/portmidi"
)

// Runner keeps track of its own outputs
type Runner struct {
	outs    []Out
	project projectfile.Project
}

// Out represents the portmidi stream and the associated DeviceID
type Out struct {
	DeviceID portmidi.DeviceID
	Stream   *portmidi.Stream
	Channel  portmidi.Channel
}

// InitializeFromProject func
func InitializeFromProject(project projectfile.Project) Runner {
	return Runner{
		project: project,
	}
}

// CloseOuts close all the open outs
func (runner Runner) CloseOuts() {
	for _, out := range runner.outs {
		out.Stream.Close()
	}
}

// Run func
func (runner Runner) Run() {
	runner.runProject(runner.project)
}

func (runner Runner) runProject(project projectfile.Project) {

	events := []grid.MidiPoint{}

	for _, part := range project.Parts {
		events = append(events, part.Events...)
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].Event.Timestamp < events[j].Event.Timestamp
	})

	fmt.Println(events)

	for _, midiPoint := range events {
		out, isNew := runner.findOrCreateOut(midiPoint.DeviceID, midiPoint.Channel)

		if isNew {
			runner.outs = append(runner.outs, out)
		}

		out.Stream.Write([]portmidi.Event{midiPoint.Event})
	}
}

func (runner Runner) findOrCreateOut(deviceID portmidi.DeviceID, channel portmidi.Channel) (Out, bool) {
	for _, out := range runner.outs {
		if out.DeviceID == deviceID && out.Channel == channel {
			return out, false
		}
	}

	newOut, err := portmidi.NewOutputStream(deviceID, 1024, 1)
	if channel != -1 {
		newOut.SetChannelMask(int(channel))
	}

	if err != nil {
		log.Fatal(err)
	}

	return Out{Stream: newOut, DeviceID: deviceID, Channel: channel}, true
}
