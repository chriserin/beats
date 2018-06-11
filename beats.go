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
	time.Sleep(12 * time.Second)
	fmt.Println(portmidi.Time())
	defer runner.CloseOuts()

	defer portmidi.Terminate()
}
