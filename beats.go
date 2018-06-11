package main

import (
	"fmt"
	"log"

	"./grid"
	"./pattern_file"
	"github.com/rakyll/portmidi"
)

func main() {

	portmidi.Initialize()
	defer portmidi.Terminate()

	pattern := patternfile.Parse("pattern_file/fixtures/simple_pattern.pf")
	midiPoints := grid.TransformGridToMidi(pattern.GridText)
	deviceID := portmidi.DefaultOutputDeviceID()

	out, err := portmidi.NewOutputStream(deviceID, 1024, 0)
	if err != nil {
		log.Fatal(err)
	}

	info := portmidi.Info(deviceID)
	count := portmidi.CountDevices()

	fmt.Println(count)

	fmt.Println("INFO")
	fmt.Println(info)

	fmt.Println("Writing note")
	out.WriteShort(0x90, 60, 100)
	out.Write(midiPoints)
	out.Close()

}
