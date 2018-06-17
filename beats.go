package main

import (
	"fmt"
	"os"
	"time"

	"./projectfile"
	"./runner"
	"github.com/rakyll/portmidi"
)

func main() {

	args := os.Args[1:]

	var fileName string
	if len(args) == 0 {
		fileName = "./projectfile/fixtures/key_pattern_repeat_project.bp"
	} else {
		fileName = args[0]
	}

	portmidi.Initialize()

	project := projectfile.Parse(fileName)
	runner := runner.InitializeFromProject(project)
	runner.Run()
	time.Sleep(time.Duration(project.Length) * time.Millisecond)
	fmt.Println(portmidi.Time())
	defer runner.CloseOuts()

	defer portmidi.Terminate()
}
