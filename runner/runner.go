package runner

import (
	"fmt"
	"log"

	"../devices"
	"../grid"
	"../projectfile"

	"github.com/rakyll/portmidi"
)

// Runner keeps track of its own outputs
type Runner struct {
	outs    []out
	project projectfile.Project
}

type out struct {
	deviceID portmidi.DeviceID
	stream   *portmidi.Stream
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
		out.stream.Close()
	}
}

// Run func
func (runner Runner) Run() {
	runner.runProject(runner.project)
}

func (runner Runner) runProject(project projectfile.Project) {
	fmt.Println("Running Project")
	fmt.Println(project)

	for _, part := range project.Parts {
		for _, pattern := range part.Patterns {
			deviceID := devices.FindDeviceID(pattern.DeviceName)
			stream := runner.findOrCreateOut(deviceID)
			runner.outs = append(runner.outs, out{stream: stream, deviceID: deviceID})
		}
	}

	for _, part := range project.Parts {
		for _, pattern := range part.Patterns {
			deviceID := devices.FindDeviceID(pattern.DeviceName)
			midiPoints := grid.TransformGridToMidi(pattern.GridText)
			runner.schedule(midiPoints, deviceID)
		}
	}
}

func (runner Runner) schedule(midiPoints []portmidi.Event, deviceID portmidi.DeviceID) {
	out := runner.findOrCreateOut(deviceID)
	fmt.Println(midiPoints)
	out.Write(midiPoints)
}

func (runner Runner) findOrCreateOut(deviceID portmidi.DeviceID) *portmidi.Stream {
	for _, out := range runner.outs {
		if out.deviceID == deviceID {
			return out.stream
		}
	}

	newOut, err := portmidi.NewOutputStream(deviceID, 1024, 1)
	if err != nil {
		log.Fatal(err)
	}

	return newOut
}
