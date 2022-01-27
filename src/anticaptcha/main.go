package main

import (
	"anticaptcha/v1"

	"github.com/robomotionio/robomotion-go/runtime"
)

func main() {

	runtime.RegisterNodes(
		&anticaptcha.Image{},
		&anticaptcha.ReCaptcha{},
		&anticaptcha.HCaptcha{},
	)

	runtime.Start()

}
