package runner

import (
	"fmt"
	"log"

	"../projectfile"

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

	for _, part := range project.Parts {
		for i, midiPoint := range part.Events {
			fmt.Printf("Scheduling point %d\n", i)
			fmt.Println(midiPoint.Event)
			out, isNew := runner.findOrCreateOut(midiPoint.DeviceID)
			if isNew {
				runner.outs = append(runner.outs, out)
			}
			out.Stream.Write([]portmidi.Event{midiPoint.Event})
		}
	}
}

func (runner Runner) findOrCreateOut(deviceID portmidi.DeviceID) (Out, bool) {
	fmt.Printf("Finding Device for %d\n", int(deviceID))
	for _, out := range runner.outs {
		if out.DeviceID == deviceID {
			return out, false
		}
	}

	newOut, err := portmidi.NewOutputStream(deviceID, 1024, 1)
	if err != nil {
		log.Fatal(err)
	}

	return Out{Stream: newOut, DeviceID: deviceID}, true
}
