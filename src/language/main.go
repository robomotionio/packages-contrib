package main

import (
	language "main.go/v1"

	"github.com/robomotionio/robomotion-go/runtime"
)

func main() {
	runtime.RegisterNodes(
		&language.Connect{},
		&language.AnalyzeEntities{},
		&language.AnalyzeSentiment{},
		&language.Disconnect{},
	)
	runtime.Start()

}
