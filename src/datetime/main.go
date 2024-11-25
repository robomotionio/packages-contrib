package main

import (
	"datetime/v1"

	"github.com/robomotionio/robomotion-go/runtime"
)

func main() {

	runtime.RegisterNodes(
		&datetime.Add{},
		&datetime.Format{},
		&datetime.Now{},
		&datetime.Span{},
		&datetime.Leap{},
		&datetime.Split{},
	)

	runtime.Start()

}
