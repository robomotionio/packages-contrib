package main

import (
	"datetime/nodes"

	"github.com/robomotionio/robomotion-go/runtime"
)

func main() {

	runtime.RegisterNodes(
		&nodes.Add{},
		&nodes.Format{},
		&nodes.Now{},
		&nodes.Span{},
		&nodes.Leap{},
		&nodes.Split{},
	)

	runtime.Start()

}
