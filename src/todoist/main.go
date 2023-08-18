package main

import (
	"github.com/robomotionio/robomotion-go/runtime"
	"todoist/v1"
)

func main() {
	runtime.RegisterNodes(
		&todoist.Connect{},
		&todoist.GetProjects{},
		&todoist.AddProject{},
		&todoist.DeleteProject{},
		&todoist.GetTasks{},
		&todoist.AddTask{},
		&todoist.CloseTask{},
		&todoist.DeleteTask{},
		&todoist.Disconnect{},
	)
	runtime.Start()

}
