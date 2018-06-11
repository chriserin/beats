package main

import (
	"fmt"
	"time"

	"./projectfile"
	"./runner"
	"github.com/rakyll/portmidi"
)

func main() {
	portmidi.Initialize()

	project := projectfile.Parse("./projectfile/fixtures/simple_project.bp")
	runner := runner.InitializeFromProject(project)
	runner.Run()
	time.Sleep(time.Duration(project.Length) * time.Millisecond)
	fmt.Println(portmidi.Time())
	defer runner.CloseOuts()

	defer portmidi.Terminate()
}
